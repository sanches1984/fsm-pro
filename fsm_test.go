package fsm_pro

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFSMSuccess(t *testing.T) {
	ctx := context.Background()
	fsmClient := getFSMClient()

	cases := []struct {
		object        *Object
		startState    State
		expectedState State
	}{
		{
			object:        &Object{MoveToState3: false},
			startState:    State1,
			expectedState: State2,
		},
		{
			object:        &Object{MoveToState3: true},
			startState:    State1,
			expectedState: State3,
		},
		{
			object:        &Object{MoveToState4: false},
			startState:    State2,
			expectedState: State3,
		},
		{
			object:        &Object{MoveToState4: true},
			startState:    State2,
			expectedState: State4,
		},
		{
			object:        &Object{},
			startState:    State3,
			expectedState: State4,
		},
	}

	for _, c := range cases {
		fsm, err := fsmClient.Init(c.startState, State4)
		require.NoError(t, err)

		err = fsm.MakeTransition(ctx, c.object)
		require.NoError(t, err)
		require.Equal(t, c.expectedState, fsm.CurrentState())
	}
}

func TestFSMError(t *testing.T) {
	ctx := context.Background()
	fsmClient := getFSMClient()

	cases := []struct {
		object     interface{}
		startState State
		err        error
	}{
		{
			object:     "some text",
			startState: State1,
			err:        errors.New("bad object"),
		},
		{
			object:     &Object{ThrowState3Error: true},
			startState: State3,
			err:        errors.New("state error"),
		},
		{
			object:     &Object{ThrowEvent34Error: true},
			startState: State3,
			err:        errors.New("event error"),
		},
		{
			object:     &Object{},
			startState: State4,
			err:        errors.New("process is finished"),
		},
	}

	for _, c := range cases {
		fsm, err := fsmClient.Init(c.startState, State4)
		require.NoError(t, err)

		err = fsm.MakeTransition(ctx, c.object)
		require.EqualError(t, err, c.err.Error())
	}
}

//-----------------------------------------------------------------
// Test data

const (
	State1 State = "state1"
	State2 State = "state2"
	State3 State = "state3"
	State4 State = "state4"

	Event12 Event = "st1-st2"
	Event13 Event = "st1-st3"
	Event23 Event = "st2-st3"
	Event24 Event = "st2-st4"
	Event34 Event = "st3-st4"
)

type Object struct {
	MoveToState3      bool
	MoveToState4      bool
	ThrowState3Error  bool
	ThrowEvent34Error bool
}

func getFSMClient() FSMBuilder {
	client := New()
	client.RegisterState(State1, NewState1Processor())
	client.RegisterState(State2, NewState2Processor())
	client.RegisterState(State3, NewState3Processor())
	client.RegisterState(State4, nil)
	client.RegisterEvent(Event12, NewTestEventProcessor())
	client.RegisterEvent(Event13, NewTestEventProcessor())
	client.RegisterEvent(Event23, NewTestEventProcessor())
	client.RegisterEvent(Event24, NewTestEventProcessor())
	client.RegisterEvent(Event34, NewEvent34Processor())
	client.RegisterTransition(State1, State2, Event12)
	client.RegisterTransition(State1, State3, Event13)
	client.RegisterTransition(State2, State3, Event23)
	client.RegisterTransition(State2, State4, Event24)
	client.RegisterTransition(State3, State4, Event34)
	return client
}

type state1Processor struct {
}

func NewState1Processor() StateProcessor {
	return &state1Processor{}
}

func (p *state1Processor) Process(ctx context.Context, args ...interface{}) (Event, error) {
	if len(args) == 0 {
		return "", errors.New("no object")
	}
	data, ok := args[0].(*Object)
	if !ok {
		return "", errors.New("bad object")
	}

	if data.MoveToState3 {
		return Event13, nil
	}
	return Event12, nil
}

type state2Processor struct {
}

func NewState2Processor() StateProcessor {
	return &state2Processor{}
}

func (p *state2Processor) Process(ctx context.Context, args ...interface{}) (Event, error) {
	if len(args) == 0 {
		return "", errors.New("no object")
	}
	data, ok := args[0].(*Object)
	if !ok {
		return "", errors.New("bad object")
	}

	if data.MoveToState4 {
		return Event24, nil
	}
	return Event23, nil
}

type state3Processor struct {
}

func NewState3Processor() StateProcessor {
	return &state3Processor{}
}

func (p *state3Processor) Process(ctx context.Context, args ...interface{}) (Event, error) {
	if len(args) == 0 {
		return "", errors.New("no object")
	}
	data, ok := args[0].(*Object)
	if !ok {
		return "", errors.New("bad object")
	}

	if data.ThrowState3Error {
		return "", errors.New("state error")
	}
	return Event34, nil
}

type event34Processor struct {
}

func NewEvent34Processor() EventProcessor {
	return &event34Processor{}
}

func (p *event34Processor) Process(ctx context.Context, args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("no object")
	}
	data, ok := args[0].(*Object)
	if !ok {
		return errors.New("bad object")
	}

	if data.ThrowEvent34Error {
		return errors.New("event error")
	}
	return nil
}
