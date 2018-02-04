package simutil

type SaturatingCounter struct {
	minValue     uint32
	threshold    uint32
	maxValue     uint32
	value        uint32
	initialValue uint32
}

func NewSaturatingCounter(minValue uint32, threshold uint32, maxValue uint32, initialValue uint32) *SaturatingCounter {
	var saturatingCounter = &SaturatingCounter{
		minValue:     minValue,
		threshold:    threshold,
		maxValue:     maxValue,
		initialValue: initialValue,
	}

	return saturatingCounter
}

func (saturatingCounter *SaturatingCounter) MinValue() uint32 {
	return saturatingCounter.minValue
}

func (saturatingCounter *SaturatingCounter) Threshold() uint32 {
	return saturatingCounter.threshold
}

func (saturatingCounter *SaturatingCounter) MaxValue() uint32 {
	return saturatingCounter.maxValue
}

func (saturatingCounter *SaturatingCounter) Value() uint32 {
	return saturatingCounter.value
}

func (saturatingCounter *SaturatingCounter) InitialValue() uint32 {
	return saturatingCounter.initialValue
}

func (saturatingCounter *SaturatingCounter) Reset() {
	saturatingCounter.value = saturatingCounter.initialValue
}

func (saturatingCounter *SaturatingCounter) Update(direction bool) {
	if direction {
		saturatingCounter.increment()
	} else {
		saturatingCounter.decrement()
	}
}

func (saturatingCounter *SaturatingCounter) increment() {
	if saturatingCounter.value < saturatingCounter.maxValue {
		saturatingCounter.value++
	}
}

func (saturatingCounter *SaturatingCounter) decrement() {
	if saturatingCounter.value > saturatingCounter.minValue {
		saturatingCounter.value--
	}
}

func (saturatingCounter *SaturatingCounter) Taken() bool {
	return saturatingCounter.value >= saturatingCounter.threshold
}
