package chess

import "time"

const (
	INFINITY = 999999999
)

type TimeControl interface {
	InitialTime() time.Duration
	Increment() time.Duration
}

type ThreeMinute struct{}

func (tm ThreeMinute) InitialTime() time.Duration {
	return 3 * time.Minute
}

func (tm ThreeMinute) Increment() time.Duration {
	return 0 * time.Second
}

type InfiniteTime struct{}

func (it InfiniteTime) InitialTime() time.Duration {
	return INFINITY * time.Second
}

func (it InfiniteTime) Increment() time.Duration {
	return 0 * time.Second
}
