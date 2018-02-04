package cpu

import "github.com/mcai/heo/cpu/uncore"

type Core interface {
	Processor() *Processor
	Threads() []Thread
	AddThread(thread Thread)
	Num() int32

	FastForwardOneCycle()

	FUPool() *FUPool

	WaitingInstructionQueue() []GeneralReorderBufferEntry
	SetWaitingInstructionQueue(waitingInstructionQueue []GeneralReorderBufferEntry)

	ReadyInstructionQueue() []GeneralReorderBufferEntry
	SetReadyInstructionQueue(readyInstructionQueue []GeneralReorderBufferEntry)

	ReadyLoadQueue() []GeneralReorderBufferEntry
	SetReadyLoadQueue(readyLoadQueue []GeneralReorderBufferEntry)

	WaitingStoreQueue() []GeneralReorderBufferEntry
	SetWaitingStoreQueue(waitingStoreQueue []GeneralReorderBufferEntry)

	ReadyStoreQueue() []GeneralReorderBufferEntry
	SetReadyStoreQueue(readyStoreQueue []GeneralReorderBufferEntry)

	OoOEventQueue() []GeneralReorderBufferEntry
	SetOoOEventQueue(oooEventQueue []GeneralReorderBufferEntry)

	L1IController() *uncore.L1IController
	L1DController() *uncore.L1DController

	CanIfetch(thread Thread, virtualAddress uint32) bool
	CanLoad(thread Thread, virtualAddress uint32) bool
	CanStore(thread Thread, virtualAddress uint32) bool

	Ifetch(thread Thread, virtualAddress uint32, virtualPc uint32, onCompletedCallback func())
	Load(thread Thread, virtualAddress uint32, virtualPc uint32, onCompletedCallback func())
	Store(thread Thread, virtualAddress uint32, virtualPc uint32, onCompletedCallback func())

	RemoveFromQueues(entryToRemove GeneralReorderBufferEntry)

	NumDynamicInsts() int64

	ResetStats()
}

type BaseCore struct {
	processor *Processor
	threads   []Thread
	num       int32
}

func NewBaseCore(processor *Processor, num int32) *BaseCore {
	var core = &BaseCore{
		processor: processor,
		num:       num,
	}

	return core
}

func (core *BaseCore) Processor() *Processor {
	return core.processor
}

func (core *BaseCore) Threads() []Thread {
	return core.threads
}

func (core *BaseCore) AddThread(thread Thread) {
	core.threads = append(core.threads, thread)
}

func (core *BaseCore) Num() int32 {
	return core.num
}

func (core *BaseCore) FastForwardOneCycle() {
	for _, thread := range core.Threads() {
		thread.FastForwardOneCycle()
	}
}

func (core *BaseCore) NumDynamicInsts() int64 {
	var numDynamicInsts = int64(0)

	for _, thread := range core.Threads() {
		numDynamicInsts += thread.NumDynamicInsts()
	}

	return numDynamicInsts
}

func (core *BaseCore) ResetStats() {
	for _, thread := range core.Threads() {
		thread.ResetStats()
	}
}
