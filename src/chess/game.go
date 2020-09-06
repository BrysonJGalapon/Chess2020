package chess

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Player interface {
	Init(c Color, g GameClient, prompt chan Prompt, move chan *Move)
	Run()
}

type Prompt struct {
	OppMove *Move
}

type Game struct {
	board *Board

	timeControl TimeControl

	blackTimeLeft time.Duration
	whiteTimeLeft time.Duration

	whitePlayer Player
	blackPlayer Player

	timestamp time.Time

	promptWhite chan Prompt
	promptBlack chan Prompt

	moveWhite chan *Move
	moveBlack chan *Move
}

type GameClient interface {
	GetBoard() *Board
	GetTimeLeft(c Color) time.Duration
	GetTimeControl() TimeControl
}

func NewGame(white, black Player, tc TimeControl) *Game {
	return &Game{
		board: NewBoard(),

		timeControl: tc,

		blackTimeLeft: tc.InitialTime(),
		whiteTimeLeft: tc.InitialTime(),

		whitePlayer: white,
		blackPlayer: black,

		timestamp: time.Unix(0, 0), // last time the current player to move was prompted

		promptWhite: make(chan Prompt, 1),
		promptBlack: make(chan Prompt, 1),

		moveWhite: make(chan *Move, 1),
		moveBlack: make(chan *Move, 1),
	}
}

func (g *Game) GetTimeControl() TimeControl {
	return g.timeControl
}

func (g *Game) GetBoard() *Board {
	return g.board.Copy()
}

func (g *Game) GetTimeLeft(c Color) time.Duration {
	switch {
	case c == White && g.board.Turn == Black:
		return g.whiteTimeLeft
	case c == Black && g.board.Turn == White:
		return g.blackTimeLeft
	case c == White && g.board.Turn == White:
		return g.whiteTimeLeft - time.Since(g.timestamp)
	case c == Black && g.board.Turn == Black:
		return g.blackTimeLeft - time.Since(g.timestamp)
	default:
		panic("Unhandled turn case")
	}
}

func (g *Game) handleMove(c Color) (*Move, error) {
	switch c {
	case White:
		ctx, cancel := context.WithTimeout(context.Background(), g.whiteTimeLeft)
		defer cancel()

		select {
		case m := <-g.moveWhite:
			tmp := *m

			err := g.board.Move(&tmp)
			if err != nil {
				return nil, fmt.Errorf("white made an invalid move: %v", err)
			}

			g.whiteTimeLeft -= time.Since(g.timestamp)
			g.whiteTimeLeft += g.timeControl.Increment()

			return &tmp, nil
		case <-ctx.Done():
			return nil, fmt.Errorf("white ran out of time")
		}

	case Black:
		ctx, cancel := context.WithTimeout(context.Background(), g.blackTimeLeft)
		defer cancel()
		select {
		case m := <-g.moveBlack:
			tmp := *m

			err := g.board.Move(&tmp)
			if err != nil {
				return nil, fmt.Errorf("black made an invalid move: %v", err)
			}

			g.blackTimeLeft -= time.Since(g.timestamp)
			g.blackTimeLeft += g.timeControl.Increment()

			return &tmp, nil
		case <-ctx.Done():
			return nil, fmt.Errorf("black ran out of time")
		}
	default:
		panic("Unhandled color type")
	}
}

func (g *Game) Start() {
	wp := g.whitePlayer
	bp := g.blackPlayer

	wp.Init(White, g, g.promptWhite, g.moveWhite)
	bp.Init(Black, g, g.promptBlack, g.moveBlack)

	go wp.Run()
	go bp.Run()

	g.promptWhite <- Prompt{}
	g.timestamp = time.Now()

	var m *Move
	var err error

game:
	for {
		// white's turn
		m, err = g.handleMove(White)

		switch {
		case err != nil:
			fmt.Println()
			fmt.Println(g.board.String())
			log.Printf("White lost: %v", err)
			break game
		case g.board.IsCheckmate():
			fmt.Println()
			fmt.Println(g.board.String())
			log.Printf("White won via checkmate")
			break game
		case g.board.IsStalemate():
			fmt.Println()
			fmt.Println(g.board.String())
			log.Printf("White drew via stalemate")
			break game
		case g.board.InsufficientMaterial():
			fmt.Println()
			fmt.Println(g.board.String())
			log.Printf("White drew via insufficient mating material")
			break game
		}

		g.promptBlack <- Prompt{m}
		g.timestamp = time.Now()

		// black's turn
		m, err = g.handleMove(Black)

		switch {
		case err != nil:
			fmt.Println()
			fmt.Println(g.board.String())
			log.Printf("Black lost: %v", err)
			break game
		case g.board.IsCheckmate():
			fmt.Println()
			fmt.Println(g.board.String())
			log.Printf("Black won via checkmate")
			break game
		case g.board.IsStalemate():
			fmt.Println()
			fmt.Println(g.board.String())
			log.Printf("Black drew via stalemate")
			break game
		case g.board.InsufficientMaterial():
			fmt.Println()
			fmt.Println(g.board.String())
			log.Printf("Black drew via insufficient mating material")
			break game
		}

		g.promptWhite <- Prompt{m}
		g.timestamp = time.Now()

		// TODO cap the number of moves in a game to 200
		// TODO handle 3-fold repetition
		// TODO handle 50-move rule
	}
}
