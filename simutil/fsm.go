package simutil

type FiniteStateMachine interface {
	State() interface{}
	PreviousState() interface{}
	SetState(condition interface{}, params interface{}, state interface{})
}

type BaseFiniteStateMachine struct {
	state, previousState interface{}
	settingStates        bool
}

func NewBaseFiniteStateMachine(state interface{}) *BaseFiniteStateMachine {
	var fsm = &BaseFiniteStateMachine{
		state: state,
	}

	return fsm
}

func (fsm *BaseFiniteStateMachine) State() interface{} {
	return fsm.state
}

func (fsm *BaseFiniteStateMachine) PreviousState() interface{} {
	return fsm.previousState
}

func (fsm *BaseFiniteStateMachine) SetState(condition interface{}, params interface{}, state interface{}) {
	if fsm.settingStates {
		panic("Impossible")
	}

	fsm.settingStates = true

	fsm.previousState = fsm.state

	fsm.state = state

	fsm.settingStates = false
}

type StateTransition struct {
	State     interface{}
	Condition interface{}
	NewState  interface{}
	Action    func(fsm FiniteStateMachine, condition interface{}, params interface{})
}

func NewStateTransition(state interface{}, condition interface{}, newState interface{}, action func(fsm FiniteStateMachine, condition interface{}, params interface{})) *StateTransition {
	var stateTransition = &StateTransition{
		State:     state,
		Condition: condition,
		NewState:  newState,
		Action:    action,
	}

	return stateTransition
}

type StateTransitions struct {
	fsmFactory          *FiniteStateMachineFactory
	state               interface{}
	perStateTransitions map[interface{}]*StateTransition
	onCompletedCallback func(fsm FiniteStateMachine)
}

func NewStateTransitions(fsmFactory *FiniteStateMachineFactory, state interface{}) *StateTransitions {
	var stateTransitions = &StateTransitions{
		fsmFactory:          fsmFactory,
		state:               state,
		perStateTransitions: make(map[interface{}]*StateTransition),
	}

	return stateTransitions
}

func (stateTransitions *StateTransitions) SetOnCompletedCallback(onCompletedCallback func(fsm FiniteStateMachine)) *StateTransitions {
	stateTransitions.onCompletedCallback = onCompletedCallback
	return stateTransitions
}

func (stateTransitions *StateTransitions) OnCondition(condition interface{}, action func(fsm FiniteStateMachine, condition interface{}, params interface{}), newState interface{}) *StateTransitions {
	stateTransitions.perStateTransitions[condition] = NewStateTransition(stateTransitions.state, condition, newState, action)
	return stateTransitions
}

func (stateTransitions *StateTransitions) fireTransition(fsm FiniteStateMachine, condition interface{}, params interface{}) {
	var stateTransition = stateTransitions.perStateTransitions[condition]

	stateTransition.Action(fsm, condition, params)

	var newState = stateTransition.NewState

	fsm.SetState(condition, params, newState)

	if stateTransitions.onCompletedCallback != nil {
		stateTransitions.onCompletedCallback(fsm)
	}
}

type FiniteStateMachineFactory struct {
	transitions map[interface{}]*StateTransitions
}

func NewFiniteStateMachineFactory() *FiniteStateMachineFactory {
	var fsmFactory = &FiniteStateMachineFactory{
		transitions: make(map[interface{}]*StateTransitions),
	}

	return fsmFactory
}

func (fsmFactory *FiniteStateMachineFactory) InState(state interface{}) *StateTransitions {
	if _, ok := fsmFactory.transitions[state]; !ok {
		fsmFactory.transitions[state] = NewStateTransitions(fsmFactory, state)
	}

	return fsmFactory.transitions[state]
}

func (fsmFactory *FiniteStateMachineFactory) FireTransition(fsm FiniteStateMachine, condition interface{}, params interface{}) {
	fsmFactory.transitions[fsm.State()].fireTransition(fsm, condition, params)
}
