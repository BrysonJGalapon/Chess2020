package chess

import (
	"testing"
)

func newMove(c1, c2 Coordinate) *Move {
	m, _ := NewMoveCoord(c1, c2)
	return m
}

func newMovePromotion(c1, c2 Coordinate, p Piece) *Move {
	m, _ := NewMoveCoordPromotion(c1, c2, p)
	return m
}

func TestCheckmate(t *testing.T) {
	var b *Board

	b = NewBoard()
	_ = b.Move(newMove("f2", "f3"))
	_ = b.Move(newMove("e7", "e6"))
	_ = b.Move(newMove("g2", "g4"))

	if b.IsCheckmate() {
		t.Fatalf("Position is not checkmate, but board returns true:\n%v", b.String())
	}

	_ = b.Move(newMove("d8", "h4"))
	if !b.IsCheckmate() {
		t.Fatalf("Position is checkmate, but board returns false:\n%v", b.String())
	}

	/////////////////////////

	b = NewBoard()

	_ = b.Move(newMove("e2", "e4"))
	_ = b.Move(newMove("e7", "e6"))
	_ = b.Move(newMove("a2", "a3"))
	_ = b.Move(newMove("d8", "h4"))
	_ = b.Move(newMove("g1", "f3"))

	_ = b.Move(newMove("h4", "e4"))
	if b.IsCheckmate() {
		t.Fatalf("Position is not checkmate (bishop or queen can block), but board returns true:\n%v", b.String())
	}

	/////////////////////////
	b = NewBoard()

	_ = b.Move(newMove("e2", "e4"))
	_ = b.Move(newMove("e7", "e6"))
	_ = b.Move(newMove("a2", "a3"))
	_ = b.Move(newMove("d8", "h4"))
	_ = b.Move(newMove("g1", "f3"))

	_ = b.Move(newMove("h4", "f2"))
	if b.IsCheckmate() {
		t.Fatalf("Position is not checkmate (king can capture), but board returns true:\n%v", b.String())
	}

	/////////////////////////
	b = NewBoard()

	_ = b.Move(newMove("d2", "d4"))
	_ = b.Move(newMove("e7", "e6"))
	_ = b.Move(newMove("b1", "d2"))
	_ = b.Move(newMove("f8", "b4"))
	_ = b.Move(newMove("e2", "e4"))
	_ = b.Move(newMove("d8", "h4"))
	_ = b.Move(newMove("g1", "f3"))

	_ = b.Move(newMove("h4", "e4"))
	if b.IsCheckmate() {
		t.Fatalf("Position is not checkmate (bishop must block), but board returns true:\n%v", b.String())
	}

	/////////////////////////
	b = NewBoard()

	_ = b.Move(newMove("d2", "d3"))
	_ = b.Move(newMove("e7", "e6"))
	_ = b.Move(newMove("b2", "b4"))
	_ = b.Move(newMove("c7", "c6"))
	_ = b.Move(newMove("g2", "g3"))
	_ = b.Move(newMove("g7", "g6"))
	_ = b.Move(newMove("b4", "b5"))
	_ = b.Move(newMove("f8", "g7"))
	_ = b.Move(newMove("c1", "a3"))
	_ = b.Move(newMove("d8", "a5"))
	_ = b.Move(newMove("a3", "b4"))
	_ = b.Move(newMove("h7", "h5"))
	_ = b.Move(newMove("f1", "h3"))
	_ = b.Move(newMove("h8", "h7"))
	_ = b.Move(newMove("h3", "e6"))
	_ = b.Move(newMove("h7", "h8"))
	_ = b.Move(newMove("e6", "b3"))
	_ = b.Move(newMove("h8", "h7"))
	_ = b.Move(newMove("b1", "a3"))
	_ = b.Move(newMove("h7", "h8"))
	_ = b.Move(newMove("c2", "c4"))
	_ = b.Move(newMove("h8", "h7"))
	_ = b.Move(newMove("d1", "b1"))
	_ = b.Move(newMove("h7", "h8"))
	_ = b.Move(newMove("g1", "h3"))
	_ = b.Move(newMove("h8", "h7"))
	_ = b.Move(newMove("h1", "f1"))
	_ = b.Move(newMove("h7", "h8"))
	_ = b.Move(newMove("b3", "d1"))

	_ = b.Move(newMove("g7", "c3"))
	if b.IsCheckmate() {
		t.Fatalf("Position is not checkmate (bishop must capture), but board returns true:\n%v", b.String())
	}

	/////////////////////////
	b = NewBoard()

	_ = b.Move(newMove("e2", "e4"))
	_ = b.Move(newMove("e7", "e5"))
	_ = b.Move(newMove("g1", "f3"))
	_ = b.Move(newMove("d7", "d6"))
	_ = b.Move(newMove("d2", "d4"))
	_ = b.Move(newMove("c8", "g4"))
	_ = b.Move(newMove("d4", "e5"))
	_ = b.Move(newMove("g4", "f3"))
	_ = b.Move(newMove("d1", "f3"))
	_ = b.Move(newMove("d6", "e5"))
	_ = b.Move(newMove("f1", "c4"))
	_ = b.Move(newMove("g8", "f6"))
	_ = b.Move(newMove("f3", "b3"))
	_ = b.Move(newMove("d8", "e7"))
	_ = b.Move(newMove("b1", "c3"))
	_ = b.Move(newMove("c7", "c6"))
	_ = b.Move(newMove("c1", "g5"))
	_ = b.Move(newMove("b7", "b5"))
	_ = b.Move(newMove("c3", "b5"))
	_ = b.Move(newMove("c6", "b5"))
	_ = b.Move(newMove("c4", "b5"))
	_ = b.Move(newMove("b8", "d7"))
	_ = b.Move(newMove("e1", "c1"))
	_ = b.Move(newMove("a8", "d8"))
	_ = b.Move(newMove("d1", "d7"))
	_ = b.Move(newMove("d8", "d7"))
	_ = b.Move(newMove("h1", "d1"))
	_ = b.Move(newMove("e7", "e6"))
	_ = b.Move(newMove("b5", "d7"))
	_ = b.Move(newMove("f6", "d7"))
	_ = b.Move(newMove("b3", "b8"))
	_ = b.Move(newMove("d7", "b8"))

	_ = b.Move(newMove("d1", "d8"))
	if !b.IsCheckmate() {
		t.Fatalf("Position is checkmate, but board returns false:\n%v", b.String())
	}
}

func TestMove(t *testing.T) {
	b := NewBoard()

	var err error
	i := 0

	// 1. e4
	i++
	if err = b.Move(newMove("e2", "e4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// making white move on black's turn
	i++
	if err = b.Move(newMove("a2", "a4")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// e5
	i++
	if err = b.Move(newMove("e7", "e5")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// trying to move piece that isn't there
	i++
	if err = b.Move(newMove("e2", "e3")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// 2. d4
	i++
	if err = b.Move(newMove("d2", "d4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// exd4
	i++
	if err = b.Move(newMove("e5", "d4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 3. c4
	i++
	if err = b.Move(newMove("c2", "c4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// dxc3
	i++
	if err = b.Move(newMove("d4", "c3")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// invalid knight move
	i++
	if err = b.Move(newMove("b1", "c4")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// 4. Nxc3
	i++
	if err = b.Move(newMove("b1", "c3")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// d6
	i++
	if err = b.Move(newMove("d7", "d6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't capture non-existent pawn
	i++
	if err = b.Move(newMove("e4", "d5")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// 5. e5
	i++
	if err = b.Move(newMove("e4", "e5")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// f5
	i++
	if err = b.Move(newMove("f7", "f5")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 6.Nf3
	i++
	if err = b.Move(newMove("g1", "f3")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't move queen through pawns (1)
	i++
	if err = b.Move(newMove("d8", "d5")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// can't move queen through pawns (2)
	i++
	if err = b.Move(newMove("d8", "a5")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// Qh4
	i++
	if err = b.Move(newMove("d8", "h4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't en-passent after more than 1 move
	i++
	if err = b.Move(newMove("e5", "f6")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// 7. Bb5+
	i++
	if err = b.Move(newMove("f1", "b5")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// doesn't resolve the check
	i++
	if err = b.Move(newMove("h4", "f2")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// c6
	i++
	if err = b.Move(newMove("c7", "c6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// invalid queen move
	i++
	if err = b.Move(newMove("d8", "e3")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// 8. Qxd6
	i++
	if err = b.Move(newMove("d1", "d6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't move pinned piece
	i++
	if err = b.Move(newMove("c6", "c5")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// Bxd6
	i++
	if err = b.Move(newMove("f8", "d6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 9. Bg5
	i++
	if err = b.Move(newMove("c1", "g5")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't castle kingside, due to knight blocking
	i++
	if err = b.Move(newMove("e8", "g8")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// can't land on your own piece
	i++
	if err = b.Move(newMove("b8", "c6")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	//  Na6
	i++
	if err = b.Move(newMove("b8", "a6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 10. O-O-O
	i++
	if err = b.Move(newMove("e1", "c1")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// Be6
	i++
	if err = b.Move(newMove("c8", "e6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 11. Rd4
	i++
	if err = b.Move(newMove("d1", "d4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't castle through check
	i++
	if err = b.Move(newMove("e8", "c8")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// Nf6
	i++
	if err = b.Move(newMove("g8", "f6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 12. Bxc6+
	i++
	if err = b.Move(newMove("b5", "c6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't castle out of check (1)
	i++
	if err = b.Move(newMove("e8", "c8")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// can't castle out of check (2)
	i++
	if err = b.Move(newMove("e8", "c8")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// bxc6
	i++
	if err = b.Move(newMove("b7", "c6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 13. exf6
	i++
	if err = b.Move(newMove("e5", "f6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// Rg8
	i++
	if err = b.Move(newMove("h8", "g8")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 14. fxg7
	i++
	if err = b.Move(newMove("f6", "g7")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// Rh8
	i++
	if err = b.Move(newMove("g8", "h8")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 15. Rc4
	i++
	if err = b.Move(newMove("d4", "c4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// Qh6
	i++
	if err = b.Move(newMove("h4", "h6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't promote to opposite color piece
	i++
	if err = b.Move(newMovePromotion("g7", "h8", BlackKnight)); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// can't promote to pawn
	i++
	if err = b.Move(newMovePromotion("g7", "h8", WhitePawn)); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// can't promote to king
	i++
	if err = b.Move(newMovePromotion("g7", "h8", WhiteKing)); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// 16. gxh8=Q+
	i++
	if err = b.Move(newMovePromotion("g7", "h8", WhiteQueen)); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// Bf8
	i++
	if err = b.Move(newMove("d6", "f8")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't move pinned bishop
	i++
	if err = b.Move(newMove("g5", "e7")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// 17. Bxh6
	i++
	if err = b.Move(newMove("g5", "h6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// f4
	i++
	if err = b.Move(newMove("f5", "f4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 18. Rxc6
	i++
	if err = b.Move(newMove("c4", "c6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't castle into check
	i++
	if err = b.Move(newMove("e8", "c8")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// Nc7
	i++
	if err = b.Move(newMove("a6", "c7")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 19. Re1
	i++
	if err = b.Move(newMove("h1", "e1")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// Kd7
	i++
	if err = b.Move(newMove("e8", "d7")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 20. h3
	i++
	if err = b.Move(newMove("h2", "h3")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// Ke8
	i++
	if err = b.Move(newMove("d7", "e8")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 21. g4
	i++
	if err = b.Move(newMove("g2", "g4")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't castle, as king has moved
	i++
	if err = b.Move(newMove("e8", "c8")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// Kf7
	i++
	if err = b.Move(newMove("e8", "f7")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 22. Qxf8+
	i++
	if err = b.Move(newMove("h8", "f8")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// can't capture protected piece with king
	i++
	if err = b.Move(newMove("f7", "f8")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an error, but got: %v", err)
	}

	// Kg6
	i++
	if err = b.Move(newMove("f7", "g6")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// 23. Qg7#
	i++
	if err = b.Move(newMove("f8", "g7")); err != nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected no errors, but got: %v", err)
	}

	// game over, no legal moves
	i++
	if err = b.Move(newMove("a8", "g8")); err == nil {
		t.Errorf("Test %v failed", i)
		t.Fatalf("Expected an errors, but got: %v", err)
	}
}
