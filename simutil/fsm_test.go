package simutil

import (
	"fmt"
	"testing"
)

func TestFiniteStateMachineFactory(t *testing.T) {
	var fsmFactory = NewFiniteStateMachineFactory()

	fsmFactory.InState(0).OnCondition(
		"hello",
		func(fsm FiniteStateMachine, condition interface{}, params interface{}) {
			fmt.Printf("params[a] = %s\n", params.(map[string]string)["a"])
		},
		1)

	var fsm = NewBaseFiniteStateMachine(0)

	fsmFactory.FireTransition(fsm, "hello", map[string]string{
		"a": "testA",
	})
}
