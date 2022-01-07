package fsm_pro

import "context"

type State string
type Event string

type StateProcessor interface {
	Process(ctx context.Context, args ...interface{}) (Event, error)
}

type EventProcessor interface {
	Process(ctx context.Context, args ...interface{}) error
}

type Transition struct {
	StartState State
	Event      Event
	EndState   State
}
