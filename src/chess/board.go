package chess

import "fmt"

var coord map[bitmap][2]uint8

func init() {
	coord = make(map[bitmap][2]uint8)

	var s bitmap = 1
	var x uint8
	var y uint8

	for y = 0; y < 8; y++ {
		for x = 7; x <= 7; x-- {
			coord[s] = [2]uint8{x, y}
			s <<= 1
		}
	}
}

type bitmap uint64

func (b bitmap) Coordinates() [2]uint8 {
	return coord[b]
}

// Color is either White or Black
type Color uint8

// Color representations
const (
	White = iota
	Black = iota
)

// Board is a representation of a Chess board.
type Board struct {
	// Pieces[0]  = White King
	// Pieces[1]  = White Queens
	// Pieces[2]  = White Knights
	// Pieces[3]  = White Bishops
	// Pieces[4]  = White Rooks
	// Pieces[5]  = White Pawns
	// Pieces[6]  = Black King
	// Pieces[7]  = Black Queens
	// Pieces[8]  = Black Knights
	// Pieces[9]  = Black Bishops
	// Pieces[10] = Black Rooks
	// Pieces[11] = Black Pawns
	Pieces [12]bitmap

	// 1 in locs where "invisible" pawns can be captured via en-passent
	EnPassent bitmap

	CanWhiteCastleKingside  bool
	CanWhiteCastleQueenside bool

	CanBlackCastleKingside  bool
	CanBlackCastleQueenside bool

	Turn Color

	AllPieces      *bitmap
	WhiteAttackMap *bitmap
	BlackAttackMap *bitmap
}

// NewBoard creates a new board, and returns it.
func NewBoard() *Board {
	return &Board{
		Pieces:                  [12]bitmap{8, 16, 66, 36, 129, 65280, 576460752303423488, 1152921504606846976, 4755801206503243776, 2594073385365405696, 9295429630892703744, 71776119061217280},
		EnPassent:               bitmap(0),
		CanWhiteCastleKingside:  true,
		CanWhiteCastleQueenside: true,
		CanBlackCastleKingside:  true,
		CanBlackCastleQueenside: true,
		Turn:                    White,
	}
}

// Copy creates a copy of this board
func (b *Board) Copy() *Board {
	nb := *b

	// Don't copy pointers
	if b.AllPieces != nil {
		tmp := *b.AllPieces
		nb.AllPieces = &tmp
	}

	if b.WhiteAttackMap != nil {
		tmp := *b.WhiteAttackMap
		nb.WhiteAttackMap = &tmp
	}

	if b.BlackAttackMap != nil {
		tmp := *b.BlackAttackMap
		nb.BlackAttackMap = &tmp
	}

	return &nb
}

// Move performs a move on this board. Returns an error if the move is invalid.
func (b *Board) Move(m *Move) error {
	if err := b.CheckMove(m); err != nil {
		return err
	}

	b.UnsafeMove(m)

	return nil
}

// String returns a human-readable representation of this board.
func (b *Board) String() string {
	var board string
	var cursor bitmap = 1 << 63
	var p Piece
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			p = b.detectPiece(cursor)
			if b.EnPassent&cursor != 0 { // ignore "invisible" pawns on EnPassent squares
				p = EmptyPiece
			}
			board += p.String()
			cursor >>= 1
		}

		board += "\n"
	}

	return board
}

// Returns the piece on a given square on this board. Returns nil if no piece exists on the square.
// Returns an "invisible" pawn of the appropriate color if square is an EnPassent square.
func (b *Board) detectPiece(square bitmap) Piece {
	if b.EnPassent&square != 0 {
		if b.Turn == White {
			return BlackPawn
		}

		return WhitePawn
	}

	for _, pt := range AllPieceTypes {
		if b.Pieces[pt]&square != 0 {
			return pt
		}
	}

	return EmptyPiece
}

type moveCache struct {
	FromPiece *Piece
	ToPiece   *Piece
	EnPassent *bitmap
}

func (c *moveCache) resolveFromPiece(b *Board, fromSquare bitmap) {
	if c.FromPiece == nil {
		x := b.detectPiece(fromSquare)
		c.FromPiece = &x
	}
}

func (c *moveCache) resolveToPiece(b *Board, toSquare bitmap) {
	if c.ToPiece == nil {
		x := b.detectPiece(toSquare)
		c.ToPiece = &x
	}
}

func (c *moveCache) resolveEnPassent(b *Board, toSquare bitmap) {
	if c.EnPassent == nil {
		var x bitmap
		if b.EnPassent == toSquare {
			x = toSquare
		}

		c.EnPassent = &x // c.EnPassent != 0 iff current move is EnPassent
	}
}

// UnsafeMove performs a move on this board without any validity checking.
func (b *Board) UnsafeMove(m *Move) {
	b.unsafeMoveWithCache(m, &moveCache{})
}

func (b *Board) unsafeMoveWithCache(m *Move, c *moveCache) {
	fromSquare := bitmap(m.From)
	toSquare := bitmap(m.To)

	c.resolveFromPiece(b, fromSquare)
	c.resolveToPiece(b, toSquare)
	c.resolveEnPassent(b, toSquare)

	// because validity check was skipped, assumes:
	//  0. toSquare != fromSquare
	//	1. fromPiece is not empty
	// 	2. if toPiece is not empty, that toPiece is opposite color of fromPiece
	//  3. if move is a promotion, that:
	//      a. fromPiece is a pawn
	//      b. toPiece is one of {Q,N,B,R} of the same color as fromPiece
	//      c. toSquare is on the 1st rank if fromPiece is black, and on the 8th rank if fromPiece is white
	//  4. if a pawn moved up (or down) two squares, it started from its initial position
	//  5. if a king is moved 2 squares to left or right:
	//		a. the king has never been moved before
	//		b. there is a rook that has never been moved in the corner of the direction of movement
	//		c. the king is not currently in check
	// 		d. the squares between the king and the rook are empty
	//		e. the two squares the king is moving through are not controlled by a piece of the opposite color
	fromPiece := *c.FromPiece
	toPiece := *c.ToPiece

	// moves fromPiece fromSquare -> toSquare
	b.Pieces[fromPiece] ^= (fromSquare | toSquare)
	if m.Promotion != EmptyPiece {
		b.Pieces[fromPiece] ^= toSquare   // remove piece from destination square
		b.Pieces[m.Promotion] ^= toSquare // add piece to destination square
	}

	// removes toPiece from existing square
	if toPiece != EmptyPiece {
		var capturedSquare bitmap

		capturedSquare = toSquare

		if *c.EnPassent != 0 {
			if b.Turn == White {
				capturedSquare >>= 8 // captured pawn is one square below destination square
			} else {
				capturedSquare <<= 8 // captured pawn is one square above destination square
			}
		}

		b.Pieces[toPiece] ^= capturedSquare
	}

	// check if EnPassent is possible
	switch {
	case fromPiece == WhitePawn && (fromSquare<<16 == toSquare): // white pawn moved up 2 squares
		b.EnPassent = fromSquare << 8
	case fromPiece == BlackPawn && (fromSquare>>16 == toSquare): // black pawn moved down 2 squares
		b.EnPassent = fromSquare >> 8
	default:
		b.EnPassent = 0
	}

	switch {
	case fromPiece == WhiteKing && fromSquare>>2 == toSquare: // white castling king-side
		b.Pieces[WhiteRook] ^= 5
	case fromPiece == WhiteKing && fromSquare<<2 == toSquare: // white castling queen-side
		b.Pieces[WhiteRook] ^= 144
	case fromPiece == BlackKing && fromSquare>>2 == toSquare: // black castling king-side
		b.Pieces[BlackRook] ^= 360287970189639680
	case fromPiece == BlackKing && fromSquare<<2 == toSquare: // black castling queen-side
		b.Pieces[BlackRook] ^= 10376293541461622784
	}

	// moving the king strips castling rights
	if fromPiece == WhiteKing {
		b.CanWhiteCastleKingside = false
		b.CanWhiteCastleQueenside = false
	}

	if fromPiece == BlackKing {
		b.CanBlackCastleKingside = false
		b.CanBlackCastleQueenside = false
	}

	// moving the rook strips castling rights on that side
	if fromPiece == WhiteRook && fromSquare == 128 {
		b.CanWhiteCastleQueenside = false
	}

	if fromPiece == WhiteRook && fromSquare == 1 {
		b.CanWhiteCastleKingside = false
	}

	if fromPiece == BlackRook && fromSquare == 9223372036854775808 {
		b.CanBlackCastleQueenside = false
	}

	if fromPiece == BlackRook && fromSquare == 72057594037927936 {
		b.CanBlackCastleKingside = false
	}

	// pieces changed, reset cache
	b.AllPieces = nil
	b.WhiteAttackMap = nil
	b.BlackAttackMap = nil

	// toggle turn
	if b.Turn == White {
		b.Turn = Black
	} else {
		b.Turn = White
	}
}

// UndoLastMove undos the last move played on this board.
func (b *Board) UndoLastMove() {
	// TODO
}

// CheckMove checks if a given move is valid on this board. Returns an error if the move is invalid.
func (b *Board) CheckMove(m *Move) error {
	return b.checkMoveWithCache(m, &moveCache{})
}

// InCheck returns true iff the king of the given color is currently in check
func (b *Board) InCheck(c Color) bool {
	if c == White {
		return b.dynamicAttackMap(Black)&b.Pieces[WhiteKing] != 0
	}

	return b.dynamicAttackMap(White)&b.Pieces[BlackKing] != 0
}

func (b *Board) dynamicAttackMap(c Color) bitmap {
	switch c {
	case White:
		if b.WhiteAttackMap != nil {
			return *b.WhiteAttackMap
		}

		var dam bitmap = 0
		for _, p := range WhitePieceTypes {
			dam |= b.attackedSquares(p)
		}

		b.WhiteAttackMap = &dam
		return dam
	case Black:
		if b.BlackAttackMap != nil {
			return *b.BlackAttackMap
		}

		var dam bitmap = 0
		for _, p := range BlackPieceTypes {
			dam |= b.attackedSquares(p)
		}

		b.BlackAttackMap = &dam
		return dam
	default:
		panic("Unhandled Color case")
	}
}

func (b bitmap) String() string {
	var board string
	var cursor bitmap = 1 << 63
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if cursor&b != 0 {
				board += "x"
			} else {
				board += "-"
			}
			cursor >>= 1
		}

		board += "\n"
	}

	return board
}

func (b *Board) scan(start, guide bitmap, d Direction) bitmap {
	obstacles := b.allPieces() ^ start
	var s bitmap = start

	var cursor bitmap = start
	for (cursor&guide)&^obstacles != 0 {
		s ^= cursor

		switch d {
		case N:
			cursor <<= 8
		case S:
			cursor >>= 8
		case E:
			cursor >>= 1
		case W:
			cursor <<= 1
		case NE:
			cursor <<= 7
		case NW:
			cursor <<= 9
		case SE:
			cursor >>= 9
		case SW:
			cursor >>= 7
		default:
			panic("Unhandled direction")
		}
	}

	return s ^ (cursor & guide)
}

// attackedSquares returns the set of currently attacked squares of given piece
func (b *Board) attackedSquares(p Piece) bitmap {
	var attacked bitmap = 0
	var cursor bitmap = 1
	for i := 0; i < 64; i++ {
		// ignore empty squares
		if cursor&b.Pieces[p] == 0 {
			cursor <<= 1
			continue
		}

		guide := AttackMap[p][cursor] ^ cursor

		switch p {
		case WhiteKnight, BlackKnight, WhitePawn, BlackPawn, WhiteKing, BlackKing:
			attacked |= AttackMap[p][cursor]
		case WhiteRook, BlackRook:
			attacked |= (b.scan(cursor, guide, N) | b.scan(cursor, guide, S) | b.scan(cursor, guide, E) | b.scan(cursor, guide, W))
		case WhiteBishop, BlackBishop:
			attacked |= (b.scan(cursor, guide, NE) | b.scan(cursor, guide, NW) | b.scan(cursor, guide, SE) | b.scan(cursor, guide, SW))
		case WhiteQueen, BlackQueen:
			attacked |= (b.scan(cursor, guide, N) | b.scan(cursor, guide, S) | b.scan(cursor, guide, E) | b.scan(cursor, guide, W) | b.scan(cursor, guide, NE) | b.scan(cursor, guide, NW) | b.scan(cursor, guide, SE) | b.scan(cursor, guide, SW))
		}

		cursor <<= 1
	}

	return attacked
}

func (b *Board) checkMoveWithCache(m *Move, c *moveCache) error {
	fromSquare := bitmap(m.From)
	toSquare := bitmap(m.To)

	c.resolveFromPiece(b, fromSquare)
	c.resolveToPiece(b, toSquare)
	c.resolveEnPassent(b, toSquare)

	fromPiece := *c.FromPiece
	toPiece := *c.ToPiece

	if fromPiece == EmptyPiece {
		return fmt.Errorf("there must be a piece on the source square")
	}

	if toPiece != EmptyPiece && toPiece.Color() == fromPiece.Color() {
		return fmt.Errorf("can't capture piece of same color")
	}

	if fromPiece.Color() != b.Turn {
		return fmt.Errorf("cannot move piece of different color of turn")
	}

	if (MoveMap[fromPiece][fromSquare]|AttackMap[fromPiece][fromSquare])&toSquare == 0 {
		return fmt.Errorf("invalid piece movement of: %v", fromPiece.String())
	}

	if (fromPiece == WhitePawn || fromPiece == BlackPawn) && AttackMap[fromPiece][fromSquare]&toSquare != 0 && toPiece == EmptyPiece {
		return fmt.Errorf("pawn captures can't occur on empty squares")
	}

	if fromSquare == toSquare {
		return fmt.Errorf("destination square can't be same as source square")
	}

	if (fromPiece != WhiteKnight && fromPiece != BlackKnight) && ray(fromSquare, toSquare)&b.allPieces() != 0 {
		return fmt.Errorf("non-knight pieces are not allowed to jump over other pieces")
	}

	copy := b.Copy()
	copy.UnsafeMove(m)
	if copy.InCheck(b.Turn) {
		return fmt.Errorf("can't make a move that leaves king in check")
	}

	switch {
	case fromPiece == WhiteKing && fromSquare>>2 == toSquare: // white castling king-side
		if !b.CanWhiteCastleKingside {
			return fmt.Errorf("white is not allowed to castle kingside")
		}

		copy := b.Copy()
		copy.UnsafeMove(NewMove(m.From, m.From>>1, EmptyPiece))
		if copy.InCheck(b.Turn) {
			return fmt.Errorf("can't castle through check")
		}

		if b.InCheck(b.Turn) {
			return fmt.Errorf("can't castle out of check")
		}
	case fromPiece == WhiteKing && fromSquare<<2 == toSquare: // white castling queen-side
		if !b.CanWhiteCastleQueenside {
			return fmt.Errorf("white is not allowed to castle queenside")
		}

		copy := b.Copy()
		copy.UnsafeMove(NewMove(m.From, m.From<<1, EmptyPiece))
		if copy.InCheck(b.Turn) {
			return fmt.Errorf("can't castle through check")
		}

		if b.InCheck(b.Turn) {
			return fmt.Errorf("can't castle out of check")
		}
	case fromPiece == BlackKing && fromSquare>>2 == toSquare: // black castling king-side
		if !b.CanBlackCastleKingside {
			return fmt.Errorf("black is not allowed to castle kingside")
		}

		copy := b.Copy()
		copy.UnsafeMove(NewMove(m.From, m.From>>1, EmptyPiece))
		if copy.InCheck(b.Turn) {
			return fmt.Errorf("can't castle through check")
		}

		if b.InCheck(b.Turn) {
			return fmt.Errorf("can't castle out of check")
		}
	case fromPiece == BlackKing && fromSquare<<2 == toSquare: // black castling queen-side
		if !b.CanBlackCastleQueenside {
			return fmt.Errorf("black is not allowed to castle queenside")
		}

		copy := b.Copy()
		copy.UnsafeMove(NewMove(m.From, m.From<<1, EmptyPiece))
		if copy.InCheck(b.Turn) {
			return fmt.Errorf("can't castle through check")
		}

		if b.InCheck(b.Turn) {
			return fmt.Errorf("can't castle out of check")
		}
	}

	switch m.Promotion {
	case EmptyPiece:
		break // do nothing
	case WhitePawn, BlackPawn, WhiteKing, BlackKing:
		return fmt.Errorf("can't promote to a pawn or king")
	default:
		if fromPiece != WhitePawn && fromPiece != BlackPawn {
			return fmt.Errorf("only pawns can promote")
		}

		if m.Promotion.Color() != fromPiece.Color() {
			return fmt.Errorf("can only promote to a piece of the same color as the pawn that moved")
		}
	}

	// All checks passed
	return nil
}

func direction(fromSquare, toSquare bitmap) Direction {
	c1 := fromSquare.Coordinates()
	c2 := toSquare.Coordinates()

	x1, y1 := c1[0], c1[1]
	x2, y2 := c2[0], c2[1]

	switch {
	case y1 == y2 && x1 < x2:
		return E
	case y1 == y2 && x1 > x2:
		return W
	case x1 == x2 && y1 < y2:
		return N
	case x1 == x2 && y1 > y2:
		return S
	case x1 < x2 && y1 < y2:
		return NE
	case x1 < x2 && y1 > y2:
		return SE
	case x1 > x2 && y1 < y2:
		return NW
	case x1 > x2 && y1 > y2:
		return SW
	default:
		panic(fmt.Sprintf("Unhandled direction case c1=(%v, %v), c2=(%v,%v)", x1, x2, y1, y2))
	}
}

// ray returns a bitmap s.t. all squares between fromSquare and toSquare are all 1, else 0.
// Assumes fromSquare != toSquare and that they share a diagonal, column, or row.
func ray(fromSquare, toSquare bitmap) bitmap {
	if fromSquare > toSquare {
		// force fromSquare < toSquare, so don't need to deal with negative shifts
		return ray(toSquare, fromSquare)
	}

	var shift int
	switch direction(fromSquare, toSquare) {
	case W:
		shift = 1
	case N:
		shift = 8
	case NE:
		shift = 7
	case NW:
		shift = 9
	default:
		panic(fmt.Sprintf("Unhandled direction case: %v -> %v", fromSquare, toSquare))
	}

	// shift in the direction of ray
	var r bitmap = fromSquare
	var cursor bitmap = fromSquare
	for cursor != toSquare {
		r ^= cursor
		cursor <<= shift
	}

	return r
}

func (b *Board) Test(c Color) bitmap {
	return b.dynamicAttackMap(c)
}

func (b *Board) allPieces() bitmap {
	if b.AllPieces != nil {
		return *b.AllPieces
	}

	var m bitmap = 0
	for _, p := range AllPieceTypes {
		m |= b.Pieces[p]
	}

	b.AllPieces = &m
	return m
}
