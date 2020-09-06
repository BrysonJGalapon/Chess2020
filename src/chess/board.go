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

type CheckingPiece struct {
	Piece  Piece
	Square bitmap
	Check  bitmap
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

func (c Color) String() string {
	switch c {
	case White:
		return "White"
	case Black:
		return "Black"
	default:
		panic("Unhandled color type")
	}
}

func (c Color) Opposite() Color {
	switch c {
	case White:
		return Black
	case Black:
		return White
	default:
		panic("Unhandled color type")
	}
}

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

	WhitePieces    *bitmap
	BlackPieces    *bitmap
	WhiteAttackMap *bitmap
	BlackAttackMap *bitmap

	PinnedPieces *bitmap

	Counts [12]uint8
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
		Counts:                  [12]uint8{1, 1, 2, 2, 2, 8, 1, 1, 2, 2, 2, 8},
	}
}

// Copy creates a copy of this board
func (b *Board) Copy() *Board {
	nb := *b

	// Don't copy pointers
	if b.WhitePieces != nil {
		tmp := *b.WhitePieces
		nb.WhitePieces = &tmp
	}

	if b.BlackPieces != nil {
		tmp := *b.BlackPieces
		nb.WhitePieces = &tmp
	}

	if b.WhiteAttackMap != nil {
		tmp := *b.WhiteAttackMap
		nb.WhiteAttackMap = &tmp
	}

	if b.BlackAttackMap != nil {
		tmp := *b.BlackAttackMap
		nb.BlackAttackMap = &tmp
	}

	if b.PinnedPieces != nil {
		tmp := *b.PinnedPieces
		nb.PinnedPieces = &tmp
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

		b.Counts[fromPiece]--
		b.Counts[m.Promotion]++
	}

	// removes toPiece from existing square
	if toPiece != EmptyPiece {
		b.Counts[toPiece]--

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
	b.WhitePieces = nil
	b.BlackPieces = nil
	b.WhiteAttackMap = nil
	b.BlackAttackMap = nil
	b.PinnedPieces = nil

	// toggle turn
	b.Turn = b.Turn.Opposite()
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

func (b *Board) GetPiece(c Color, pt string) Piece {
	if c == White {
		switch pt {
		case "King":
			return WhiteKing
		case "Queen":
			return WhiteQueen
		case "Knight":
			return WhiteKnight
		case "Bishop":
			return WhiteBishop
		case "Rook":
			return WhiteRook
		case "Pawn":
			return WhitePawn
		default:
			panic("Unhandled piece case")
		}
	}

	switch pt {
	case "King":
		return BlackKing
	case "Queen":
		return BlackQueen
	case "Knight":
		return BlackKnight
	case "Bishop":
		return BlackBishop
	case "Rook":
		return BlackRook
	case "Pawn":
		return BlackPawn
	default:
		panic("Unhandled piece case")
	}
}

func (b *Board) CheckingPieces() []CheckingPiece {
	cps := make([]CheckingPiece, 0)

	var pt [6]Piece

	if b.Turn == White {
		pt = BlackPieceTypes
	} else {
		pt = WhitePieceTypes
	}

	king := b.Pieces[b.GetPiece(b.Turn, "King")]
	var pinned bitmap
	var attacked bitmap = 0
	for _, p := range pt {
		var cursor bitmap = 1

		for i := 0; i < 64; i++ {
			// ignore empty squares
			if cursor&b.Pieces[p] == 0 {
				cursor <<= 1
				continue
			}

			guide := AttackMap[p][cursor] ^ cursor

			handleScan := func(cursor, guide bitmap, dirs []Direction) {
				for _, d := range dirs {
					att, pin := b.pseudoScan(cursor, guide, d)
					attacked |= att
					pinned |= pin

					if king&att != 0 && pin == 0 {
						cps = append(cps, CheckingPiece{Piece: p, Square: cursor, Check: att})
					}
				}
			}

			switch p {
			case WhiteKnight, BlackKnight, WhitePawn, BlackPawn, WhiteKing, BlackKing:
				att := AttackMap[p][cursor]
				if king&att != 0 {
					cps = append(cps, CheckingPiece{Piece: p, Square: cursor, Check: 0})
				}
				attacked |= att
			case WhiteRook, BlackRook:
				handleScan(cursor, guide, []Direction{N, S, E, W})
			case WhiteBishop, BlackBishop:
				handleScan(cursor, guide, []Direction{NE, NW, SE, SW})
			case WhiteQueen, BlackQueen:
				handleScan(cursor, guide, []Direction{N, S, E, W, NE, NW, SE, SW})
			}

			cursor <<= 1
		}
	}

	b.PinnedPieces = &pinned

	return cps
}

func (b *Board) pieces(c Color) bitmap {
	if c == White {
		return b.whitePieces()
	}

	return b.blackPieces()
}

// IsCheckmate returns true iff the position is checkmate
func (b *Board) IsCheckmate() bool {
	king := b.GetPiece(b.Turn, "King")
	kingSquare := b.Pieces[king]
	attacked := b.dynamicAttackMap(b.Turn.Opposite())

	if attacked&kingSquare == 0 {
		return false // not in check, therefore not checkmate
	}

	cps := b.CheckingPieces()

	obstacles := b.pieces(b.Turn) &^ b.Pieces[king]
	possible := AttackMap[king][kingSquare]

	switch len(cps) {
	case 1:
		cp := cps[0]

		// check blocks, if checking piece is not a knight
		if cp.Piece != WhiteKnight && cp.Piece != BlackKnight {
			if cp.Check&(b.dynamicAttackMapQueen(b.Turn)|b.dynamicAttackMapRook(b.Turn)|b.dynamicAttackMapBishop(b.Turn)|b.dynamicMovementMapPawn(b.Turn)) != 0 {
				return false
			}
		}

		// check captures, for pieces that are unpinned
		if b.dynamicAttackMapNoPin(b.Turn)&cp.Square != 0 {
			return false
		}

		// check king moves
		return (possible&^obstacles)&^attacked == 0 // true iff there are no possible king moves that go to a blocked or attacked square
	case 2: // double-check, must respond with a king move
		return (possible&^obstacles)&^attacked == 0 // true iff there are no possible king moves that go to a blocked or attacked square
	default:
		panic("there can never be more than 2 checking pieces")
	}
}

func (b *Board) dynamicAttackMapQueen(c Color) bitmap {
	switch c {
	case White:
		return b.attackedSquaresNoPin(WhiteQueen)
	case Black:
		return b.attackedSquaresNoPin(BlackQueen)
	default:
		panic("Unhandled color type")
	}
}

func (b *Board) dynamicAttackMapBishop(c Color) bitmap {
	switch c {
	case White:
		return b.attackedSquaresNoPin(WhiteBishop)
	case Black:
		return b.attackedSquaresNoPin(BlackBishop)
	default:
		panic("Unhandled color type")
	}
}

func (b *Board) dynamicAttackMapRook(c Color) bitmap {
	switch c {
	case White:
		return b.attackedSquaresNoPin(WhiteRook)
	case Black:
		return b.attackedSquaresNoPin(BlackRook)
	default:
		panic("Unhandled color type")
	}
}

func (b *Board) dynamicMovementMapPawn(c Color) bitmap {
	var p Piece
	var d Direction

	if c == White {
		p = WhitePawn
		d = N
	} else {
		p = BlackPawn
		d = S
	}

	pinned := b.pinnedPieces()

	var cursor bitmap = 1
	var m bitmap

	for i := 0; i < 64; i++ {
		// ignore empty squares
		if cursor&b.Pieces[p] == 0 {
			cursor <<= 1
			continue
		}

		// ignore pinned pawns
		if cursor&pinned != 0 {
			cursor <<= 1
			continue
		}

		m |= (b.scan(cursor, MoveMap[p][cursor]^cursor, d)) &^ b.allPieces()

		cursor <<= 1
	}

	return m
}

// // IsCheckmate returns true iff the position is checkmate
// func (b *Board) IsCheckmate() bool {
// 	var king Piece
// 	var attacked bitmap
// 	var obstacles bitmap

// 	switch b.Turn {
// 	case White:
// 		king = WhiteKing
// 		attacked = b.dynamicAttackMap(Black)
// 		obstacles = b.whitePieces()
// 	case Black:
// 		king = BlackKing
// 		attacked = b.dynamicAttackMap(White)
// 		obstacles = b.blackPieces()
// 	default:
// 		panic("Unhandled turn case")
// 	}

// 	square := b.Pieces[king]
// 	possible := AttackMap[king][square]

// 	// in check, and all possible movement squares are empty
// 	return square&attacked != 0 && (possible&^obstacles)&^attacked == 0
// }

// IsStalemate returns true iff the position is stalemate
func (b *Board) IsStalemate() bool {
	// TODO
	return false
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

func (b *Board) dynamicAttackMapNoPin(c Color) bitmap {
	switch c {
	case White:
		if b.WhiteAttackMap != nil {
			return *b.WhiteAttackMap
		}

		var dam bitmap = 0
		for _, p := range WhitePieceTypes {
			dam |= b.attackedSquaresNoPin(p)
		}

		b.WhiteAttackMap = &dam
		return dam
	case Black:
		if b.BlackAttackMap != nil {
			return *b.BlackAttackMap
		}

		var dam bitmap = 0
		for _, p := range BlackPieceTypes {
			dam |= b.attackedSquaresNoPin(p)
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

func (b *Board) pseudoScan(start, guide bitmap, d Direction) (bitmap, bitmap) {
	same := b.pieces(b.Turn) &^ start
	opp := b.pieces(b.Turn.Opposite()) &^ start

	var s bitmap = start
	var passThrough bool = true
	var pinned bitmap

	var cursor bitmap = start
	for (cursor & guide) != 0 {
		if cursor&same != 0 {
			break
		}

		if cursor&opp != 0 {
			if !passThrough {
				break
			}
			pinned = cursor
			passThrough = false
		}

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

	if cursor != b.Pieces[b.GetPiece(b.Turn.Opposite(), "King")] {
		pinned = 0
	}

	return s ^ (cursor & guide), pinned
}

func (b *Board) scanUtil(start, guide bitmap, d Direction, passThroughKing bool) bitmap {
	obstacles := b.allPieces() ^ start
	var s bitmap = start

	if passThroughKing && b.detectPiece(start).Color() != b.Turn {
		obstacles = obstacles &^ (b.Pieces[b.GetPiece(b.Turn, "King")])
	}

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

	return (s ^ (cursor & guide))
}

func (b *Board) scan(start, guide bitmap, d Direction) bitmap {
	return b.scanUtil(start, guide, d, true)
}

func (b *Board) pinnedPieces() bitmap {
	if b.PinnedPieces != nil {
		return *b.PinnedPieces
	}

	panic("Expected *b.PinnedPieces is non-nil")
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

// attackedSquares returns the set of currently attacked squares of given piece
func (b *Board) attackedSquaresNoPin(p Piece) bitmap {
	pinned := b.pinnedPieces()

	var attacked bitmap = 0
	var cursor bitmap = 1
	for i := 0; i < 64; i++ {
		// ignore empty squares
		if cursor&b.Pieces[p] == 0 {
			cursor <<= 1
			continue
		}

		// ignore pinned pieces
		if cursor&pinned != 0 {
			cursor <<= 1
			continue
		}

		guide := AttackMap[p][cursor] ^ cursor

		switch p {
		case WhiteKing, BlackKing:
			attacked |= AttackMap[p][cursor] &^ b.dynamicAttackMap(b.Turn.Opposite())
		case WhiteKnight, BlackKnight, WhitePawn, BlackPawn:
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

	return attacked &^ b.pieces(b.Turn)
}

func arrEqual(a1, a2 [6]uint8) bool {
	for i := 0; i < 6; i++ {
		if a1[i] != a2[i] {
			return false
		}
	}

	return true
}

func (b *Board) InsufficientMaterial() bool {
	var whiteCounts [6]uint8
	var blackCounts [6]uint8

	for i := 0; i < 6; i++ {
		whiteCounts[i] = b.Counts[i]
		blackCounts[i] = b.Counts[i+6]
	}

	if arrEqual(whiteCounts, [6]uint8{1, 0, 0, 0, 0, 0}) && arrEqual(blackCounts, [6]uint8{1, 0, 0, 0, 0, 0}) {
		// king v king
		return true
	}

	if arrEqual(whiteCounts, [6]uint8{1, 0, 0, 1, 0, 0}) && arrEqual(blackCounts, [6]uint8{1, 0, 0, 0, 0, 0}) ||
		arrEqual(whiteCounts, [6]uint8{1, 0, 0, 0, 0, 0}) && arrEqual(blackCounts, [6]uint8{1, 0, 0, 1, 0, 0}) {
		// king + bishop v king
		return true
	}

	if arrEqual(whiteCounts, [6]uint8{1, 0, 1, 0, 0, 0}) && arrEqual(blackCounts, [6]uint8{1, 0, 0, 0, 0, 0}) ||
		arrEqual(whiteCounts, [6]uint8{1, 0, 0, 0, 0, 0}) && arrEqual(blackCounts, [6]uint8{1, 0, 1, 0, 0, 0}) {
		// king + knight v king
		return true
	}

	if arrEqual(whiteCounts, [6]uint8{1, 0, 0, 1, 0, 0}) && arrEqual(blackCounts, [6]uint8{1, 0, 0, 1, 0, 0}) && ((b.Pieces[WhiteBishop]&EvenSquares != 0 && b.Pieces[BlackBishop]&EvenSquares != 0) || (b.Pieces[WhiteBishop]&OddSquares != 0 && b.Pieces[BlackBishop]&OddSquares != 0)) {
		// king + bishop v king + bishop of same color
		return true
	}

	return false
}

func (b *Board) checkMoveWithCache(m *Move, c *moveCache) error {
	// if b.IsCheckmate() {
	// 	return fmt.Errorf("game is over")
	// }

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

	if (fromPiece == WhitePawn || fromPiece == BlackPawn) && toPiece != EmptyPiece && toSquare&MoveMap[fromPiece][fromSquare] != 0 {
		return fmt.Errorf("can't move a pawn onto a piece")
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
		if fromPiece == WhitePawn && toSquare >= 72057594037927936 {
			return fmt.Errorf("white pawn on 8th rank has to promote to some piece")
		}

		if fromPiece == BlackPawn && toSquare <= 128 {
			return fmt.Errorf("black pawn on 1st rank has to promote to some piece")
		}
	case WhitePawn, BlackPawn, WhiteKing, BlackKing:
		return fmt.Errorf("can't promote to a pawn or king")
	default:
		if fromPiece != WhitePawn && fromPiece != BlackPawn {
			return fmt.Errorf("only pawns can promote")
		}

		if m.Promotion.Color() != fromPiece.Color() {
			return fmt.Errorf("can only promote to a piece of the same color as the pawn that moved")
		}

		if fromPiece == WhitePawn && toSquare < 72057594037927936 {
			return fmt.Errorf("white pawn can only promote on 8th rank")
		}

		if fromPiece == BlackPawn && toSquare > 128 {
			return fmt.Errorf("black pawn can only promote on 1st rank")
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

func (b *Board) allPieces() bitmap {
	return b.whitePieces() | b.blackPieces()
}

func (b *Board) whitePieces() bitmap {
	if b.WhitePieces != nil {
		return *b.WhitePieces
	}

	var m bitmap = 0
	for _, p := range WhitePieceTypes {
		m |= b.Pieces[p]
	}

	b.WhitePieces = &m
	return m
}

func (b *Board) blackPieces() bitmap {
	if b.BlackPieces != nil {
		return *b.BlackPieces
	}

	var m bitmap = 0
	for _, p := range BlackPieceTypes {
		m |= b.Pieces[p]
	}

	b.BlackPieces = &m
	return m
}
