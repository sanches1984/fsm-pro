package fsm_pro

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_RegisterState_Success(t *testing.T) {
	c := &builder{
		states: make(map[State]StateProcessor),
	}

	p1 := NewTestStateProcessor()
	err := c.RegisterState(State1, p1)
	require.NoError(t, err)

	p2 := NewTestStateProcessor()
	err = c.RegisterState(State2, p2)
	require.NoError(t, err)

	require.Equal(t, map[State]StateProcessor{
		State1: p1,
		State2: p2,
	}, c.states)
}

func TestClient_RegisterState_Error(t *testing.T) {
	c := &builder{
		states: make(map[State]StateProcessor),
	}

	p1 := NewTestStateProcessor()
	err := c.RegisterState(State1, p1)
	require.NoError(t, err)

	p2 := NewTestStateProcessor()
	err = c.RegisterState(State1, p2)
	require.EqualError(t, err, ErrStateAlreadyRegistered.Error())
}

func TestClient_RegisterEvent_Success(t *testing.T) {
	c := &builder{
		events: make(map[Event]EventProcessor),
	}

	p1 := NewTestEventProcessor()
	err := c.RegisterEvent(Event12, p1)
	require.NoError(t, err)

	p2 := NewTestEventProcessor()
	err = c.RegisterEvent(Event23, p2)
	require.NoError(t, err)

	require.Equal(t, map[Event]EventProcessor{
		Event12: p1,
		Event23: p2,
	}, c.events)
}

func TestClient_RegisterEvent_Error(t *testing.T) {
	c := &builder{
		events: make(map[Event]EventProcessor),
	}

	p1 := NewTestEventProcessor()
	err := c.RegisterEvent(Event12, p1)
	require.NoError(t, err)

	p2 := NewTestEventProcessor()
	err = c.RegisterEvent(Event12, p2)
	require.EqualError(t, err, ErrEventAlreadyRegistered.Error())
}

func TestClient_RegisterTransition_Success(t *testing.T) {
	c := &builder{
		transitions: make([]Transition, 0, 1),
		states: map[State]StateProcessor{
			State1: nil,
			State2: nil,
			State3: nil,
		},
		events: map[Event]EventProcessor{
			Event12: nil,
			Event23: nil,
		},
	}

	err := c.RegisterTransition(State1, State2, Event12)
	require.NoError(t, err)

	err = c.RegisterTransition(State2, State3, Event23)
	require.NoError(t, err)

	require.Equal(t, []Transition{
		{
			StartState: State1,
			Event:      Event12,
			EndState:   State2,
		},
		{
			StartState: State2,
			Event:      Event23,
			EndState:   State3,
		},
	}, c.transitions)
}

func TestClient_RegisterTransition_Error(t *testing.T) {
	c := &builder{
		transitions: make([]Transition, 0, 1),
		states: map[State]StateProcessor{
			State1: nil,
			State2: nil,
			State3: nil,
		},
		events: map[Event]EventProcessor{
			Event12: nil,
			Event23: nil,
		},
	}

	err := c.RegisterTransition(State1, State2, Event12)
	require.NoError(t, err)
	require.Equal(t, []Transition{
		{
			StartState: State1,
			Event:      Event12,
			EndState:   State2,
		},
	}, c.transitions)

	err = c.RegisterTransition(State1, State2, Event23)
	require.EqualError(t, err, ErrTransitionAlreadyRegistered.Error())

	err = c.RegisterTransition(State1, State3, Event12)
	require.EqualError(t, err, ErrTransitionEventAlreadyUsed.Error())

	err = c.RegisterTransition(State4, State2, Event23)
	require.EqualError(t, err, ErrStateNotFound.Error())

	err = c.RegisterTransition(State1, State4, Event23)
	require.EqualError(t, err, ErrStateNotFound.Error())

	err = c.RegisterTransition(State2, State3, "event3")
	require.EqualError(t, err, ErrEventNotFound.Error())
}

func TestClient_Init_Success(t *testing.T) {
	c := &builder{
		states: map[State]StateProcessor{
			State1: nil,
			State2: nil,
			State3: nil,
		},
		events: map[Event]EventProcessor{
			Event12: nil,
			Event23: nil,
		},
		transitions: []Transition{
			{
				StartState: State1,
				Event:      Event12,
				EndState:   State2,
			},
			{
				StartState: State2,
				Event:      Event23,
				EndState:   State3,
			},
		},
	}

	_, err := c.Init(State1, State3)
	require.NoError(t, err)
}

func TestClient_Init_Error(t *testing.T) {
	c := &builder{
		states: map[State]StateProcessor{
			State1: nil,
			State2: nil,
			State3: nil,
		},
		events: map[Event]EventProcessor{
			Event12: nil,
		},
		transitions: make([]Transition, 0, 1),
	}

	_, err := c.Init(State1, State4)
	require.EqualError(t, err, ErrStateNotFound.Error())

	_, err = c.Init(State4, State3)
	require.EqualError(t, err, ErrStateNotFound.Error())

	_, err = c.Init(State1, State3)
	require.EqualError(t, err, ErrProcessNotSet.Error())

	c.transitions = append(c.transitions, Transition{
		StartState: State1,
		Event:      Event12,
		EndState:   State2,
	})

	_, err = c.Init(State1, State3)
	require.EqualError(t, err, ErrNotConnectedStates.Error())

	c.events[Event23] = nil

	_, err = c.Init(State1, State3)
	require.EqualError(t, err, ErrNotConnectedEvents.Error())
}

//-----------------------------------------------------------------
// Test data

type testStateProcessor struct {
}

func NewTestStateProcessor() StateProcessor {
	return &testStateProcessor{}
}

func (p *testStateProcessor) Process(ctx context.Context, args ...interface{}) (Event, error) {
	return "", nil
}

type testEventProcessor struct {
}

func NewTestEventProcessor() EventProcessor {
	return &testEventProcessor{}
}

func (p *testEventProcessor) Process(ctx context.Context, args ...interface{}) error {
	return nil
}
