package random

import (
	"Chess2020/src/chess"
	"Chess2020/src/players/utils"
	"fmt"
	"math/rand"
	"time"
)

type RandomPlayer struct {
	Color      chess.Color
	GameClient chess.GameClient
	Prompt     chan chess.Prompt
	Move       chan *chess.Move

	Board *chess.Board
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomCoordinate() chess.Coordinate {
	i := rand.Intn(len(utils.Alphabet))
	j := rand.Intn(len(utils.Digits))

	return chess.Coordinate(utils.Alphabet[i:i+1] + utils.Digits[j:j+1])
}

func randomPiece() chess.Piece {
	// all promotions, + 1 potentially no promotion
	i := rand.Intn(len(utils.PromotionPieces) + 1)
	if i == len(utils.PromotionPieces) {
		return chess.EmptyPiece
	}

	switch utils.PromotionPieces[i : i+1] {
	case "R":
		return chess.WhiteRook
	case "B":
		return chess.WhiteBishop
	case "Q":
		return chess.WhiteQueen
	case "N":
		return chess.WhiteKnight
	case "r":
		return chess.BlackRook
	case "b":
		return chess.BlackBishop
	case "q":
		return chess.BlackQueen
	case "n":
		return chess.BlackKnight
	default:
		panic("unhandled promotion piece case")
	}
}

func (rp *RandomPlayer) Init(c chess.Color, gc chess.GameClient, prompt chan chess.Prompt, move chan *chess.Move) {
	rp.Color = c
	rp.Prompt = prompt
	rp.Move = move
	rp.GameClient = gc

	rp.Board = gc.GetBoard()
}

func (rp *RandomPlayer) Run() {
	fmt.Printf("Random Player [%v] started\n", rp.Color)

	for {
		// Wait for our turn
		p := <-rp.Prompt
		if p.OppMove != nil {
			rp.Board.UnsafeMove(p.OppMove)
		}

		var m *chess.Move
		var err error
		for {
			// Generate a random move
			c1 := randomCoordinate()
			c2 := randomCoordinate()
			p := randomPiece()

			if m, err = chess.NewMoveCoordPromotion(c1, c2, p); err != nil {
				panic(fmt.Errorf("Unable to generate a move: %v", err))
			}

			// Make move on the board, until random move is legal
			if err := rp.Board.Move(m); err != nil {
				// fmt.Printf("Invalid move: %v\n", err)
				continue
			}

			break
		}

		// Send move
		rp.Move <- m
	}
}

func Player() *RandomPlayer {
	return &RandomPlayer{}
}
