package simutil

import (
	"container/heap"
)

type CycleAccurateEvent struct {
	eventQueue *CycleAccurateEventQueue
	When       int64
	Action     func()
	Id         int64
}

type CycleAccurateEventQueue struct {
	Events                   []*CycleAccurateEvent
	PerCycleEvents           []func()
	CurrentCycle             int64
	currentEventId           int64

	lastEventDispatchedCycle int64
}

func (q CycleAccurateEventQueue) Len() int {
	return len(q.Events)
}

func (q CycleAccurateEventQueue) Less(i, j int) bool {
	var x = q.Events[i]
	var y = q.Events[j]
	return x.When < y.When || (x.When == y.When && x.Id < y.Id)
}

func (q CycleAccurateEventQueue) Swap(i, j int) {
	q.Events[i], q.Events[j] = q.Events[j], q.Events[i]
}

func (q *CycleAccurateEventQueue) Push(x interface{}) {
	item := x.(*CycleAccurateEvent)
	q.Events = append(q.Events, item)
}

func (q *CycleAccurateEventQueue) Pop() interface{} {
	old := q.Events
	n := len(old)
	item := old[n - 1]
	q.Events = old[0 : n - 1]
	return item
}

func NewCycleAccurateEventQueue() *CycleAccurateEventQueue {
	var q = &CycleAccurateEventQueue{
	}

	return q
}

func (q *CycleAccurateEventQueue) Schedule(action func(), delay int) {
	q.currentEventId++

	var event = &CycleAccurateEvent{
		eventQueue:q,
		When:q.CurrentCycle + int64(delay),
		Action:action,
		Id:q.currentEventId,
	}

	heap.Push(q, event)
}

func (q *CycleAccurateEventQueue) AddPerCycleEvent(action func()) {
	q.PerCycleEvents = append(q.PerCycleEvents, action)
}

func (q *CycleAccurateEventQueue) AdvanceOneCycle() {
	for q.Len() > 0 {
		var event = heap.Pop(q).(*CycleAccurateEvent)

		if event.When > q.CurrentCycle {
			heap.Push(q, event)
			break
		}

		event.Action()
		q.lastEventDispatchedCycle = q.CurrentCycle
	}

	for _, e := range q.PerCycleEvents {
		e()
	}

	q.CurrentCycle++
}