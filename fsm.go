package fsm_pro

import "context"

type FSM interface {
	MakeTransition(ctx context.Context, args ...interface{}) error
	CurrentState() State
	IsFinished() bool
}

type fsm struct {
	currentState State
	finishState  State
	states       map[State]StateProcessor
	events       map[Event]EventProcessor
	transitions  []Transition
}

func (f *fsm) MakeTransition(ctx context.Context, args ...interface{}) error {
	if f.IsFinished() {
		return ErrProcessFinished
	}

	event, err := f.processState(ctx, f.currentState, args...)
	if err != nil {
		return err
	}

	nextState, err := f.getNextState(event)
	if err != nil {
		return err
	}

	if err := f.processEvent(ctx, event, args...); err != nil {
		return err
	}

	f.currentState = nextState
	return nil
}

func (f *fsm) CurrentState() State {
	return f.currentState
}

func (f *fsm) IsFinished() bool {
	return f.currentState == f.finishState
}

func (f *fsm) getNextState(event Event) (State, error) {
	for _, t := range f.transitions {
		if t.StartState == f.currentState && t.Event == event {
			return t.EndState, nil
		}
	}

	return f.currentState, ErrTransitionNotFound
}

func (f *fsm) processState(ctx context.Context, state State, args ...interface{}) (Event, error) {
	ps, ok := f.states[state]
	if !ok {
		return "", ErrStateNotFound
	} else if ps == nil {
		return "", ErrStateProcessorNotSet
	}

	return ps.Process(ctx, args...)
}

func (f *fsm) processEvent(ctx context.Context, event Event, args ...interface{}) error {
	pe, ok := f.events[event]
	if !ok {
		return ErrEventNotFound
	} else if pe == nil {
		return ErrEventProcessorNotSet
	}

	return pe.Process(ctx, args...)
}
