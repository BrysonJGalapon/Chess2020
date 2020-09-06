package main

import (
	"Chess2020/src/chess"
	"Chess2020/src/players/random"
)

func main() {
	white := random.Player()
	black := random.Player()

	timeControl := chess.InfiniteTime{}

	g := chess.NewGame(white, black, timeControl)
	g.Start()
}
