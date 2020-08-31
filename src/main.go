package main

import (
	"Chess2020/src/chess"
	"fmt"
	"time"
)

var totalElapsed int64 = 0
var total int64 = 0

func makeMove(b *chess.Board, f, t chess.Coordinate) {
	makeMoveUtil(b, f, t, chess.EmptyPiece)
}

func makeMoveUtil(b *chess.Board, f, t chess.Coordinate, promotion chess.Piece) {
	var m *chess.Move
	var err error

	if m, err = chess.NewMoveCoordPromotion(f, t, promotion); err != nil {
		panic(err)
	}

	start := time.Now().UnixNano()
	if err = b.Move(m); err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Millisecond)
	end := time.Now().UnixNano()

	fmt.Printf("difference: %v\n", end-start-1000000)
	totalElapsed += (end - start - 1000000)
	total += 1

	fmt.Println(b)
}

func main() {
	b := chess.NewBoard()
	fmt.Println(b)

	makeMove(b, "e2", "e4")
	makeMove(b, "e7", "e5")
	makeMove(b, "g1", "f3")
	makeMove(b, "b8", "c6")

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
	makeMoveUtil(b, "g7", "h8", chess.WhiteKnight)

	makeMove(b, "d8", "e7")
	makeMove(b, "f1", "e2")
	makeMove(b, "d7", "d6")
	makeMove(b, "e1", "g1")
	makeMove(b, "c8", "d7")
	makeMove(b, "b1", "c3")

	// makeMove(b, "e8", "d8")
	// makeMove(b, "f1", "e1")
	// makeMove(b, "d8", "e8")
	// makeMove(b, "h2", "h3")

	makeMove(b, "e8", "c8")

	fmt.Println(float64(totalElapsed) / float64(total) / 1000 / 1000 / 1000)
}
