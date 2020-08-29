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
		Turn:                    0,
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
		if b.Turn == 0 {
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

	if c.EnPassent == nil {
		var x bitmap
		if b.EnPassent == toSquare {
			x = toSquare
		}

		c.EnPassent = &x // c.EnPassent != 0 iff current move is EnPassent
	}

	// because validity check was skipped, assumes:
	//  0. toSquare != fromSquare
	//	1. fromPiece is not empty
	// 	2. if toPiece is not empty, that toPiece is opposite color of fromPiece
	//  3. if move is a promotion, that:
	//      a. fromPiece is a pawn
	//      b. toPiece is one of {Q,N,B,R} of the same color as fromPiece
	//      c. toSquare is on the 1st rank if fromPiece is black, and on the 8th rank if fromPiece is white
	//  4. if a pawn moved up (or down) two squares, it started from its initial position
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
			if b.Turn == 0 {
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

	// TODO handle castling

	// toggle turn
	b.Turn ^= 1
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
