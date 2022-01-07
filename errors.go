package fsm_pro

import "errors"

var ErrStateAlreadyRegistered = errors.New("state already registered")
var ErrStateNotFound = errors.New("state not found")
var ErrStateProcessorNotSet = errors.New("state processor not set")

var ErrEventAlreadyRegistered = errors.New("event already registered")
var ErrEventNotFound = errors.New("event not found")
var ErrEventProcessorNotSet = errors.New("event processor not set")

var ErrTransitionAlreadyRegistered = errors.New("transition already registered")
var ErrTransitionEventAlreadyUsed = errors.New("transition event already used")
var ErrTransitionNotFound = errors.New("transition not found")
var ErrProcessNotSet = errors.New("process not set")
var ErrNotConnectedStates = errors.New("process has not-connected states")
var ErrNotConnectedEvents = errors.New("process has not-connected events")
var ErrProcessFinished = errors.New("process is finished")
