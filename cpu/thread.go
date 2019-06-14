package cpu

import (
	"github.com/mcai/heo/cpu/uncore"
)

type Thread interface {
	Core() Core
	Num() int32
	Id() int32
	Context() *Context
	SetContext(context *Context)
	FastForwardOneCycle()

	Itlb() *uncore.TranslationLookasideBuffer
	Dtlb() *uncore.TranslationLookasideBuffer

	NumDynamicInsts() int64

	InstructionsPerCycle() float64
	CyclesPerInstruction() float64

	ResetStats()
}

type BaseThread struct {
	core            Core
	num             int32
	id              int32
	context         *Context
	numDynamicInsts int64
}

func NewBaseThread(core Core, num int32) *BaseThread {
	var thread = &BaseThread{
		core: core,
		num:  num,
		id:   core.Num()*core.Processor().Experiment.CPUConfig.NumThreadsPerCore + num,
	}

	return thread
}

func (thread *BaseThread) Core() Core {
	return thread.core
}

func (thread *BaseThread) Num() int32 {
	return thread.num
}

func (thread *BaseThread) Id() int32 {
	return thread.id
}

func (thread *BaseThread) Context() *Context {
	return thread.context
}

func (thread *BaseThread) SetContext(context *Context) {
	thread.context = context
}

func (thread *BaseThread) NumDynamicInsts() int64 {
	return thread.numDynamicInsts
}

func (thread *BaseThread) ResetStats() {
	thread.numDynamicInsts = 0
}

func (thread *BaseThread) FastForwardOneCycle() {
	if thread.Context() != nil && thread.Context().State == ContextState_RUNNING {
		var staticInst *StaticInst

		for {
			staticInst = thread.Context().DecodeNextStaticInst()
			staticInst.Execute(thread.Context())

			if staticInst.Mnemonic.Name != Mnemonic_NOP {
				thread.numDynamicInsts++
			}

			if !(thread.Context() != nil &&
				thread.Context().State == ContextState_RUNNING &&
				staticInst.Mnemonic.Name == Mnemonic_NOP) {
				break
			}
		}
	}
}

func (thread *BaseThread) InstructionsPerCycle() float64 {
	if thread.Core().Processor().Experiment.CycleAccurateEventQueue().CurrentCycle == 0 {
		return float64(0)
	}
	return float64(thread.numDynamicInsts) / float64(thread.Core().Processor().Experiment.CycleAccurateEventQueue().CurrentCycle)
}

func (thread *BaseThread) CyclesPerInstruction() float64 {
	if thread.numDynamicInsts == 0 {
		return float64(0)
	}

	return float64(thread.Core().Processor().Experiment.CycleAccurateEventQueue().CurrentCycle) / float64(thread.numDynamicInsts)
}
