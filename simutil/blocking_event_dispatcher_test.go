package simutil

import (
	"testing"
	"reflect"
	"fmt"
)

type HelloEvent struct {
	Name string
}

func NewHelloEvent(name string) *HelloEvent {
	var helloEvent = &HelloEvent{
		Name: name,
	}

	return helloEvent
}

func TestBlockingEventDispatcher(t *testing.T) {
	var blockingEventDispatcher = NewBlockingEventDispatcher()

	var helloEvent = NewHelloEvent("Test")

	blockingEventDispatcher.AddListener(reflect.TypeOf((*HelloEvent)(nil)), func(event interface{}) {
		fmt.Printf("hello event fired with name = %s\n", event.(*HelloEvent).Name)
	})

	blockingEventDispatcher.Dispatch(helloEvent)
}
