package main

import (
	"Chess2020/src/chess"
	"fmt"
)

func main() {
	b := chess.NewBoard()
	fmt.Println(b)

	b.Move(chess.NewMove(32768, 2147483648))
	fmt.Println(b)
}
