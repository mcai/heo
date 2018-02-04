package simutil

type Counter struct {
	initialValue uint32
	value        uint32
}

func NewCounter(initialValue uint32) *Counter {
	var counter = &Counter{
		initialValue: initialValue,
		value:        initialValue,
	}

	return counter
}

func (counter *Counter) InitialValue() uint32 {
	return counter.initialValue
}

func (counter *Counter) Value() uint32 {
	return counter.value
}

func (counter *Counter) Increment() {
	counter.value++
}

func (counter *Counter) Decrement() {
	counter.value--
}

func (counter *Counter) Reset() {
	counter.value = counter.initialValue
}
