package fsm_pro

type FSMBuilder interface {
	RegisterState(state State, condition StateProcessor) error
	RegisterEvent(event Event, condition EventProcessor) error
	RegisterTransition(startState, endState State, event Event) error
	Init(initState, finishState State) (FSM, error)
}

type builder struct {
	states      map[State]StateProcessor
	events      map[Event]EventProcessor
	transitions []Transition
}

func New() *builder {
	return &builder{
		states:      make(map[State]StateProcessor),
		events:      make(map[Event]EventProcessor),
		transitions: make([]Transition, 0, 1),
	}
}

func (c *builder) RegisterState(state State, condition StateProcessor) error {
	if _, ok := c.states[state]; ok {
		return ErrStateAlreadyRegistered
	}

	c.states[state] = condition
	return nil
}

func (c *builder) RegisterEvent(event Event, condition EventProcessor) error {
	if _, ok := c.events[event]; ok {
		return ErrEventAlreadyRegistered
	}

	c.events[event] = condition
	return nil
}

func (c *builder) RegisterTransition(startState, endState State, event Event) error {
	if _, ok := c.states[startState]; !ok {
		return ErrStateNotFound
	}
	if _, ok := c.states[endState]; !ok {
		return ErrStateNotFound
	}
	if _, ok := c.events[event]; !ok {
		return ErrEventNotFound
	}

	for _, t := range c.transitions {
		if t.StartState == startState && t.Event == event {
			return ErrTransitionEventAlreadyUsed
		} else if t.StartState == startState && t.EndState == endState {
			return ErrTransitionAlreadyRegistered
		}
	}

	c.transitions = append(c.transitions, Transition{
		StartState: startState,
		Event:      event,
		EndState:   endState,
	})
	return nil
}

func (c *builder) Init(initState, finishState State) (FSM, error) {
	if _, ok := c.states[initState]; !ok {
		return nil, ErrStateNotFound
	}
	if _, ok := c.states[finishState]; !ok {
		return nil, ErrStateNotFound
	}
	if len(c.states) == 0 || len(c.events) == 0 || len(c.transitions) == 0 {
		return nil, ErrProcessNotSet
	}
	if len(c.events) < len(c.states)-1 {
		return nil, ErrNotConnectedStates
	}
	if len(c.transitions) < len(c.events) {
		return nil, ErrNotConnectedEvents
	}

	return &fsm{
		currentState: initState,
		finishState:  finishState,
		states:       c.states,
		events:       c.events,
		transitions:  c.transitions,
	}, nil
}
