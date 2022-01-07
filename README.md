# FSM
Finite-state machine with processors.

## Builder
### Register state processors

```go
type state1Processor struct {
    // clients
}

func NewState1Processor(...) StateProcessor {
    return &state1Processor { 
    	// init clients
    }
}

func (p *state1Processor) Process(ctx context.Context, args ...interface{}) (Event, error) {
	// get objects from args
	
	// make some operations
	
	// choose event to use
	return event, nil
}
```

### Register event processors

```go
type event1Processor struct {
    // clients
}

func NewEvent1Processor(...) EventProcessor {
    return &event1Processor { 
    	// init clients
    }
}

func (p *event1Processor) Process(ctx context.Context, args ...interface{}) error {
	// get objects from args
	
	// make some operations
	
	// return error for cancel transition
	return nil
}
```

### Build FSM

```go
builder := New()
builder.RegisterState("State1", NewState1Processor(...))
builder.RegisterState("State2", NewState2Processor(...))
builder.RegisterState("State3", NewState3Processor(...))
builder.RegisterEvent("Event1", NewEvent1Processor(...))
builder.RegisterEvent("Event2", NewEvent2Processor(...))
builder.RegisterTransition("State1", "State2", "Event1")
builder.RegisterTransition("State2", "State3", "Event2")

fsm, err := builder.Init("State1", "State3")
```

## Using
### Make transition

```go
err := fsm.MakeTransition(ctx, myObject1, myObject2, ...)
```
