package interactive

import (
	"Chess2020/src/chess"
	"bufio"
	"fmt"
	"os"
	"strings"
)

var alphabet = "abcdefgh"
var digits = "12345678"
var pieces = "NQBRnqbr"

type InteractivePlayer struct {
	Color      chess.Color
	GameClient chess.GameClient
	Prompt     chan chess.Prompt
	Move       chan *chess.Move

	Board *chess.Board
}

func (ip *InteractivePlayer) Init(c chess.Color, gc chess.GameClient, prompt chan chess.Prompt, move chan *chess.Move) {
	ip.Color = c
	ip.Prompt = prompt
	ip.Move = move
	ip.GameClient = gc

	ip.Board = gc.GetBoard()
}

func (ip *InteractivePlayer) readInput(r *bufio.Reader) string {
	fmt.Print("\n>$ ")
	text, _ := r.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	return text
}

func validateCoord(c string) error {
	if len(c) != 2 {
		return fmt.Errorf("expected two characters")
	}

	fc := c[0:1]
	sc := c[1:2]

	if !strings.Contains(alphabet, fc) {
		return fmt.Errorf("first character must be in: %v", alphabet)
	}

	if !strings.Contains(digits, sc) {
		return fmt.Errorf("second character must be in %v", digits)
	}

	// all tests passed
	return nil
}

func validatePromotionPiece(p string) (chess.Piece, error) {
	if len(p) != 1 {
		return chess.EmptyPiece, fmt.Errorf("piece must be a single character in %v", pieces)
	}

	switch p {
	case "N":
		return chess.WhiteKnight, nil
	case "B":
		return chess.WhiteBishop, nil
	case "Q":
		return chess.WhiteQueen, nil
	case "R":
		return chess.WhiteRook, nil
	case "n":
		return chess.BlackKnight, nil
	case "b":
		return chess.BlackBishop, nil
	case "q":
		return chess.BlackQueen, nil
	case "r":
		return chess.BlackRook, nil
	default:
		return chess.EmptyPiece, fmt.Errorf("piece must be a single character in %v", pieces)
	}
}

func (ip *InteractivePlayer) parseInput(inp string) (*chess.Move, error) {
	f := strings.Fields(inp)

	if len(f) != 2 && len(f) != 3 {
		return nil, fmt.Errorf("expected exactly 2 or 3 fields")
	}

	var c1, c2 string
	var p chess.Piece = chess.EmptyPiece

	c1, c2 = f[0], f[1]

	if err := validateCoord(c1); err != nil {
		return nil, fmt.Errorf("bad first coordinate: %v", err)
	}

	if err := validateCoord(c2); err != nil {
		return nil, fmt.Errorf("bad second coordinate: %v", err)
	}

	if len(f) == 3 {
		var err error
		if p, err = validatePromotionPiece(f[2]); err != nil {
			return nil, fmt.Errorf("bad promotion piece: %v", err)
		}
	}

	return chess.NewMoveCoordPromotion(chess.Coordinate(c1), chess.Coordinate(c2), p)
}

func (ip *InteractivePlayer) Run() {
	fmt.Printf("Interactive Player [%v] started\n", ip.Color)

	reader := bufio.NewReader(os.Stdin)
	for {
		// Wait for our turn
		p := <-ip.Prompt
		if p.OppMove != nil {
			fmt.Printf("Interactive Player [%v] Opponent played: %v\n", ip.Color, p.OppMove.String())
			ip.Board.UnsafeMove(p.OppMove)
		}

		var m *chess.Move
		var err error
		for {
			fmt.Printf("Interactive Player [%v] to move...\n", ip.Color)
			fmt.Println()
			fmt.Println(ip.Board)

			// Get move from CLI arg
			inp := ip.readInput(reader)
			if m, err = ip.parseInput(inp); err != nil {
				fmt.Printf("Could not parse move. Try again: %v\n", err)
				continue
			}

			// Make move on the board
			if err := ip.Board.Move(m); err != nil {
				fmt.Printf("Invalid move: %v\n", err)
				continue
			}

			break
		}

		// Send move
		ip.Move <- m
	}
}

func Player() *InteractivePlayer {
	return &InteractivePlayer{}
}
