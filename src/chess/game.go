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

func (g *Game) handleMove(c Color) error {
	switch c {
	case White:
		ctx, cancel := context.WithTimeout(context.Background(), g.whiteTimeLeft)
		defer cancel()

		select {
		case m := <-g.moveWhite:
			tmp := *m

			err := g.board.Move(&tmp)
			if err != nil {
				return fmt.Errorf("white made an invalid move: %v", err)
			}

			g.whiteTimeLeft -= time.Since(g.timestamp)
			g.whiteTimeLeft += g.timeControl.Increment()

			g.promptBlack <- Prompt{&tmp}
			g.timestamp = time.Now()
			return nil
		case <-ctx.Done():
			return fmt.Errorf("white ran out of time")
		}

	case Black:
		ctx, cancel := context.WithTimeout(context.Background(), g.blackTimeLeft)
		defer cancel()

		select {
		case m := <-g.moveBlack:
			tmp := *m

			err := g.board.Move(&tmp)
			if err != nil {
				return fmt.Errorf("black made an invalid move: %v", err)
			}

			g.blackTimeLeft -= time.Since(g.timestamp)
			g.blackTimeLeft += g.timeControl.Increment()

			g.promptWhite <- Prompt{&tmp}
			g.timestamp = time.Now()
			return nil
		case <-ctx.Done():
			return fmt.Errorf("black ran out of time")
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

	var err error

game:
	for {
		// white's turn
		err = g.handleMove(White)

		switch {
		case err != nil:
			log.Printf("White lost: %v", err)
			break game
		case g.board.IsCheckmate():
			log.Printf("White won via checkmate")
			break game
		case g.board.IsStalemate():
			log.Printf("White drew via stalemate")
		}

		// black's turn
		err = g.handleMove(Black)

		switch {
		case err != nil:
			log.Printf("Black lost: %v", err)
			break game
		case g.board.IsCheckmate():
			log.Printf("Black won via checkmate")
			break game
		case g.board.IsStalemate():
			log.Printf("Black drew via stalemate")
			break game
		}

		// TODO cap the number of moves in a game to 200
		// TODO handle 3-fold repetition
		// TODO handle 50-move rule
	}
}
