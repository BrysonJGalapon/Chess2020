package chess

import "fmt"

var (
	// AllPieceTypes a list of all supported piece types
	AllPieceTypes = [12]Piece{
		WhiteKing,
		WhiteQueen,
		WhiteKnight,
		WhiteBishop,
		WhiteRook,
		WhitePawn,
		BlackKing,
		BlackQueen,
		BlackKnight,
		BlackBishop,
		BlackRook,
		BlackPawn,
	}

	// All supported black pieces
	BlackPieceTypes = [6]Piece{
		BlackKing,
		BlackPawn,
		BlackQueen,
		BlackBishop,
		BlackKnight,
		BlackRook,
	}

	WhitePieceTypes = [6]Piece{
		WhiteKing,
		WhitePawn,
		WhiteQueen,
		WhiteBishop,
		WhiteKnight,
		WhiteRook,
	}
)

// Piece Enumerations
const (
	WhiteKing   = iota
	WhiteQueen  = iota
	WhiteKnight = iota
	WhiteBishop = iota
	WhiteRook   = iota
	WhitePawn   = iota

	BlackKing   = iota
	BlackQueen  = iota
	BlackKnight = iota
	BlackBishop = iota
	BlackRook   = iota
	BlackPawn   = iota

	EmptyPiece = iota
)

// Piece represents a Chess piece
type Piece uint8

func (p Piece) Color() Color {
	switch p {
	case WhiteKing:
		return White
	case WhiteQueen:
		return White
	case WhiteKnight:
		return White
	case WhiteBishop:
		return White
	case WhiteRook:
		return White
	case WhitePawn:
		return White
	case BlackKing:
		return Black
	case BlackQueen:
		return Black
	case BlackKnight:
		return Black
	case BlackBishop:
		return Black
	case BlackRook:
		return Black
	case BlackPawn:
		return Black
	default:
		panic(fmt.Sprintf("Unhandled piece value: %d", p))
	}
}

func (p Piece) String() string {
	switch p {
	case WhiteKing:
		return "K"
	case WhiteQueen:
		return "Q"
	case WhiteKnight:
		return "N"
	case WhiteBishop:
		return "B"
	case WhiteRook:
		return "R"
	case WhitePawn:
		return "P"
	case BlackKing:
		return "k"
	case BlackQueen:
		return "q"
	case BlackKnight:
		return "n"
	case BlackBishop:
		return "b"
	case BlackRook:
		return "r"
	case BlackPawn:
		return "p"
	case EmptyPiece:
		return "-"
	default:
		panic(fmt.Sprintf("Unhandled piece value: %d", p))
	}
}
