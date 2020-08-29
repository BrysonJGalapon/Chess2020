package main

import (
	"Chess2020/src/chess"
	"fmt"
)

func makeMove(b *chess.Board, f, t chess.Coordinate) {
	makeMoveUtil(b, f, t, chess.EmptyPiece)
}

func makeMoveUtil(b *chess.Board, f, t chess.Coordinate, promotion chess.Piece) {
	var m *chess.Move
	var err error

	if m, err = chess.NewMoveCoordPromotion(f, t, promotion); err != nil {
		panic(err)
	}

	b.Move(m)
	fmt.Println(b)
}

func main() {
	b := chess.NewBoard()
	fmt.Println(b)

	makeMove(b, "e2", "e4")
	makeMove(b, "e7", "e5")
	makeMove(b, "g1", "f3")
	makeMove(b, "g8", "c6")

	makeMove(b, "d2", "d4")
	makeMove(b, "e5", "d4")
	makeMove(b, "c2", "c4")
	makeMove(b, "d4", "c3")

	makeMove(b, "e4", "e5")
	makeMove(b, "f7", "f5")
	makeMove(b, "e5", "f6")

	makeMove(b, "c3", "b2")
	makeMove(b, "f6", "g7")

	makeMoveUtil(b, "b2", "a1", chess.BlackQueen)
	makeMoveUtil(b, "g7", "g8", chess.WhiteKnight)
}
