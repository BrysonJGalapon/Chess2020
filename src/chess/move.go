package chess

// Square is a representation of a Chess square
type Square uint64

// toCoord translates this Square into a Coordinate representation
func (s Square) toCoord() Coordinate {
	// TODO
	return ""
}

// Coordinate is a human-readable representation of a Chess square
type Coordinate string

// toSquare translates this Coordinate into a Square representation
func (c Coordinate) toSquare() Square {
	// TODO
	return 0
}

// Move is a representation of a Chess move
type Move struct {
	From Square
	To   Square
}

// NewMove returns a new Move from square f to square t
func NewMove(f, t Square) *Move {
	return &Move{
		From: f,
		To:   t,
	}
}

// NewMoveCoord returns a new Move from coordinate f to coordinate t
func NewMoveCoord(f, t Coordinate) *Move {
	return NewMove(f.toSquare(), t.toSquare())
}
