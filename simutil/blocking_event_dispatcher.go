package simutil

import (
	"reflect"
)

type BlockingEventDispatcher struct {
	Listeners map[reflect.Type]([](func(event interface{})))
}

func NewBlockingEventDispatcher() *BlockingEventDispatcher {
	var dispatcher = &BlockingEventDispatcher{
		Listeners: make(map[reflect.Type]([](func(event interface{})))),
	}

	return dispatcher
}

func (dispatcher *BlockingEventDispatcher) Dispatch(event interface{}) {
	var t = reflect.TypeOf(event)

	if listeners, ok := dispatcher.Listeners[t]; ok {
		for _, listener := range listeners {
			listener(event)
		}
	}
}

func (dispatcher *BlockingEventDispatcher) AddListener(t reflect.Type, listener func(event interface{})) {
	if _, ok := dispatcher.Listeners[t]; !ok {
		dispatcher.Listeners[t] = make([](func(event interface{})), 0)
	}

	dispatcher.Listeners[t] = append(dispatcher.Listeners[t], listener)
}
