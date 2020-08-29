package chess

import "fmt"

var alphabet = "abcdefgh"
var digits = "12345678"

var letters map[byte]struct{} // contains individual letters of the alphabet
var numbers map[byte]struct{} // contains individual numbers of the set of possible digits

func init() {
	letters = make(map[byte]struct{})
	for i := 0; i < len(alphabet); i++ {
		letters[alphabet[i]] = struct{}{}
	}

	numbers = make(map[byte]struct{})
	for i := 0; i < len(digits); i++ {
		numbers[digits[i]] = struct{}{}
	}
}

// Square is a representation of a Chess square
type Square uint64

// toCoord translates this Square into a Coordinate representation
func (s Square) toCoord() (Coordinate, error) {
	// TODO
	return "", nil
}

// Coordinate is a human-readable representation of a Chess square
type Coordinate string

// toSquare translates this Coordinate into a Square representation
func (c Coordinate) toSquare() (Square, error) {
	if len(c) != 2 {
		return 0, fmt.Errorf("coordinate must be a length-2 string: %v", c)
	}

	if _, ok := letters[c[0]]; !ok {
		return 0, fmt.Errorf("first letter of coordinate must be one of {a,b,c,d,e,f,g,h}, is: %v", string(c[0]))
	}

	if _, ok := numbers[c[1]]; !ok {
		return 0, fmt.Errorf("second letter of coordinate must be one of {1,2,3,4,5,6,7,8}, is: %v", string(c[1]))
	}

	var s Square = 1

	x := c[0] - "a"[0]
	y := c[1] - "1"[0]

	return s << (8*y + (7 - x)), nil
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

// NewMoveCoord returns a new Move from coordinate f to coordinate t. Returns an
// error if Coordinate could not be parsed correctly.
func NewMoveCoord(f, t Coordinate) (*Move, error) {
	var from Square
	var to Square
	var err error

	if from, err = f.toSquare(); err != nil {
		return nil, fmt.Errorf("could not parse from coordinate: %v", err)
	}

	if to, err = t.toSquare(); err != nil {
		return nil, fmt.Errorf("could not parse to coordinate: %v", err)
	}

	return NewMove(from, to), nil
}
