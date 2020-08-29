package main

import (
	"Chess2020/src/chess"
	"fmt"
)

func main() {
	b := chess.NewBoard()
	fmt.Println(b)

	var m *chess.Move
	var err error

	if m, err = chess.NewMoveCoord("e2", "e4"); err != nil {
		panic(err)
	}

	b.Move(m)
	fmt.Println(b)

	if m, err = chess.NewMoveCoord("e7", "e5"); err != nil {
		panic(err)
	}

	b.Move(m)
	fmt.Println(b)

	if m, err = chess.NewMoveCoord("g1", "f3"); err != nil {
		panic(err)
	}

	b.Move(m)
	fmt.Println(b)
}
