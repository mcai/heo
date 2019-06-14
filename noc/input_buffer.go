package noc

import "container/list"

type InputBuffer struct {
	Flits               *list.List
	Size                int
}

func NewInputBuffer(inputVirtualChannel *InputVirtualChannel) *InputBuffer {
	var inputBuffer = &InputBuffer{
		Flits:               list.New(),
		Size:                inputVirtualChannel.InputPort.Router.Node.Network.Config().MaxInputBufferSize,
	}

	return inputBuffer
}

func (inputBuffer *InputBuffer) Push(flit *Flit) {
	if inputBuffer.Full() {
		panic("Input buffer is full")
	}

	inputBuffer.Flits.PushBack(flit)
}

func (inputBuffer *InputBuffer) Peek() *Flit {
	if inputBuffer.Flits.Len() > 0 {
		return inputBuffer.Flits.Front().Value.(*Flit)
	} else {
		return nil
	}
}

func (inputBuffer *InputBuffer) Pop() {
	var e = inputBuffer.Flits.Front()
	inputBuffer.Flits.Remove(e)
}

func (inputBuffer *InputBuffer) Full() bool {
	return inputBuffer.Size <= inputBuffer.Flits.Len()
}

func (inputBuffer *InputBuffer) Count() int {
	return inputBuffer.Flits.Len()
}

func (inputBuffer *InputBuffer) FreeSlots() int {
	return inputBuffer.Size - inputBuffer.Flits.Len()
}
