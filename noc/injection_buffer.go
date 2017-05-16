package noc

import "container/list"

type InjectionBuffer struct {
	Router  *Router
	Packets *list.List
	Size    int
}

func NewInjectionBuffer(router *Router) *InjectionBuffer {
	var injectionBuffer = &InjectionBuffer{
		Router:router,
		Packets:list.New(),
		Size: router.Node.Network.Config.MaxInjectionBufferSize,
	}

	return injectionBuffer
}

func (injectionBuffer *InjectionBuffer) Push(packet Packet) {
	if injectionBuffer.Full() {
		panic("Injection buffer is full")
	}

	injectionBuffer.Packets.PushBack(packet)
}

func (injectionBuffer *InjectionBuffer) Peek() Packet {
	if injectionBuffer.Packets.Len() > 0 {
		return injectionBuffer.Packets.Front().Value.(Packet)
	} else {
		return nil
	}
}

func (injectionBuffer *InjectionBuffer) Pop() {
	var e = injectionBuffer.Packets.Front()
	injectionBuffer.Packets.Remove(e)
}

func (injectionBuffer *InjectionBuffer) Full() bool {
	return injectionBuffer.Size <= injectionBuffer.Packets.Len()
}

func (injectionBuffer *InjectionBuffer) Count() int {
	return injectionBuffer.Packets.Len()
}

func (injectionBuffer *InjectionBuffer) FreeSlots() int {
	return injectionBuffer.Size - injectionBuffer.Packets.Len()
}
