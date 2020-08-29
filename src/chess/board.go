package chess

type bitmap uint64

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

	Moves []*Move

	CanWhiteCastleKingside  bool
	CanWhiteCastleQueenside bool

	CanBlackCastleKingside  bool
	CanBlackCastleQueenside bool

	Turn uint8 // 0 --> White, 1 --> Black
}

// NewBoard creates a new board, and returns it.
func NewBoard() *Board {
	return &Board{
		Pieces:                  [12]bitmap{8, 16, 66, 36, 129, 65280, 576460752303423488, 1152921504606846976, 4755801206503243776, 2594073385365405696, 9295429630892703744, 71776119061217280},
		EnPassent:               bitmap(0),
		Moves:                   make([]*Move, 0),
		CanWhiteCastleKingside:  true,
		CanWhiteCastleQueenside: true,
		CanBlackCastleKingside:  true,
		CanBlackCastleQueenside: true,
	}
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
			board += p.String()
			cursor >>= 1
		}

		board += "\n"
	}

	return board
}

// Returns the piece on a given square on this board. Returns nil if no piece exists on the square.
func (b *Board) detectPiece(square bitmap) Piece {
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
	// EnPassent *bitmap
}

// UnsafeMove performs a move on this board without any validity checking.
func (b *Board) UnsafeMove(m *Move) {
	b.unsafeMoveWithCache(m, &moveCache{})
}

func (b *Board) unsafeMoveWithCache(m *Move, c *moveCache) {
	fromSquare := bitmap(m.From)
	toSquare := bitmap(m.To)

	if c.FromPiece == nil {
		x := b.detectPiece(fromSquare)
		c.FromPiece = &x
	}

	if c.ToPiece == nil {
		x := b.detectPiece(toSquare)
		c.ToPiece = &x
	}

	// if c.EnPassent == nil {
	// 	var x bitmap
	// 	if b.EnPassent == toSquare {
	// 		x = toSquare
	// 	}

	// 	c.EnPassent = &x // c.EnPassent != 0 iff EnPassent ocurred
	// }

	// because validity check was skipped, assumes:
	//  0. toSquare != fromSquare
	//	1. fromPiece is not empty
	// 	2. if toPiece is not empty, that toPiece is opposite color of fromPiece
	fromPiece := *c.FromPiece
	toPiece := *c.ToPiece

	// moves fromPiece fromSquare -> toSquare
	// TODO handle promotion
	b.Pieces[fromPiece] ^= (fromSquare | toSquare)

	// removes toPiece from existing square
	// TODO handle en-passent
	if toPiece != EmptyPiece {
		b.Pieces[toPiece] ^= toSquare
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

func (b *Board) checkMoveWithCache(m *Move, c *moveCache) error {
	// TODO
	return nil
}
