package main

import (
	"Chess2020/src/chess"
	"Chess2020/src/players/interactive"
)

func main() {
	g := chess.NewGame(
		interactive.Player(),
		interactive.Player(),
		chess.InfiniteTime{},
	)

	g.Start()
}
