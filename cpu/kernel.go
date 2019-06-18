package cpu

import (
	"github.com/mcai/heo/cpu/mem"
	"github.com/mcai/heo/cpu/regs"
)

type Kernel struct {
	Experiment *CPUExperiment

	Pipes            []*Pipe
	SystemEvents     []SystemEvent
	SignalActions    []*SignalAction
	Contexts         []*Context
	Processes        []*Process
	SyscallEmulation *SyscallEmulation

	CurrentPid          int32
	CurrentProcessId    int32
	CurrentMemoryId     int32
	CurrentMemoryPageId int32
	CurrentContextId    int32
	CurrentFd           int32
}

func NewKernel(experiment *CPUExperiment) *Kernel {
	var kernel = &Kernel{
		Experiment:       experiment,
		SyscallEmulation: NewSyscallEmulation(),
		CurrentPid:       1000,
		CurrentMemoryId:  0,
		CurrentContextId: 0,
		CurrentFd:        100,
	}

	for i := 0; i < MAX_SIGNAL; i++ {
		kernel.SignalActions = append(kernel.SignalActions, NewSignalAction())
	}

	return kernel
}

func (kernel *Kernel) LoadContexts() {
	for _, contextMapping := range kernel.Experiment.CPUConfig.ContextMappings {
		var context = LoadContext(kernel, contextMapping)

		if !kernel.Map(context, func(candidateThreadId int32) bool {
			return candidateThreadId == contextMapping.ThreadId
		}) {
			panic("Impossible")
		}

		kernel.Contexts = append(kernel.Contexts, context)
	}
}

func (kernel *Kernel) GetContextFromProcessId(processId int32) *Context {
	for _, context := range kernel.Contexts {
		if context.ProcessId == processId {
			return context
		}
	}

	return nil
}

func (kernel *Kernel) Map(contextToMap *Context, predicate func(candidateThreadId int32) bool) bool {
	if contextToMap.ThreadId != -1 {
		panic("Impossible")
	}

	for coreNum := int32(0); coreNum < kernel.Experiment.CPUConfig.NumCores; coreNum++ {
		for threadNum := int32(0); threadNum < kernel.Experiment.CPUConfig.NumThreadsPerCore; threadNum++ {
			var threadId = coreNum*kernel.Experiment.CPUConfig.NumThreadsPerCore + threadNum

			var hasMapped = false

			for _, context := range kernel.Contexts {
				if context.ThreadId == threadId {
					hasMapped = true
					break
				}
			}

			if !hasMapped && predicate(threadId) {
				contextToMap.ThreadId = threadId
				return true
			}
		}
	}

	return false
}

func (kernel *Kernel) ProcessSystemEvents() {
	var systemEventsToReserve []SystemEvent

	for _, e := range kernel.SystemEvents {
		if (e.Context().State == ContextState_RUNNING || e.Context().State == ContextState_BLOCKED) && !e.Context().Speculative && e.NeedProcess() {
			e.Process()
		} else {
			systemEventsToReserve = append(systemEventsToReserve, e)
		}
	}

	kernel.SystemEvents = systemEventsToReserve
}

func (kernel *Kernel) ProcessSignals() {
	for _, context := range kernel.Contexts {
		if (context.State == ContextState_RUNNING || context.State == ContextState_BLOCKED) && !context.Speculative {
			for signal := uint32(1); signal <= MAX_SIGNAL; signal++ {
				if kernel.MustProcessSignal(context, signal) {
					kernel.RunSignalHandler(context, signal)
				}
			}
		}
	}
}

func (kernel *Kernel) CreatePipe() []int32 {
	var fileDescriptors = make([]int32, 2)

	fileDescriptors[0] = kernel.CurrentFd

	kernel.CurrentFd++

	fileDescriptors[1] = kernel.CurrentFd

	kernel.CurrentFd++

	kernel.Pipes = append(kernel.Pipes, NewPipe(fileDescriptors))

	return fileDescriptors
}

func (kernel *Kernel) getBuffer(fileDescriptor int32, index uint32) *mem.CircularByteBuffer {
	for _, pipe := range kernel.Pipes {
		if pipe.FileDescriptors[index] == fileDescriptor {
			return pipe.Buffer
		}
	}

	return nil
}

func (kernel *Kernel) GetReadBuffer(fileDescriptor int32) *mem.CircularByteBuffer {
	return kernel.getBuffer(fileDescriptor, 0)
}

func (kernel *Kernel) GetWriteBuffer(fileDescriptor int32) *mem.CircularByteBuffer {
	return kernel.getBuffer(fileDescriptor, 1)
}

func (kernel *Kernel) RunSignalHandler(context *Context, signal uint32) {
	if kernel.SignalActions[signal-1].Handler == 0 {
		panic("Impossible")
	}

	context.SignalMasks.Pending.Clear(signal)

	var oldRegs = context.Regs().Clone()

	context.Regs().Gpr[regs.REGISTER_A0] = signal
	context.Regs().Gpr[regs.REGISTER_T9] = kernel.SignalActions[signal-1].Handler
	context.Regs().Gpr[regs.REGISTER_RA] = 0xffffffff
	context.Regs().Npc = kernel.SignalActions[signal-1].Handler
	context.Regs().Nnpc = context.Regs().Npc + 4

	for context.State == ContextState_RUNNING && context.Regs().Npc != 0xffffffff {
		context.DecodeNextStaticInst().Execute(context)
	}

	context.SetRegs(oldRegs)
}

func (kernel *Kernel) MustProcessSignal(context *Context, signal uint32) bool {
	return context.SignalMasks.Pending.Contains(signal) && !context.SignalMasks.Blocked.Contains(signal)
}

func (kernel *Kernel) AdvanceOneCycle() {
	if kernel.Experiment.CycleAccurateEventQueue().CurrentCycle%1000 == 0 {
		kernel.ProcessSystemEvents()
		kernel.ProcessSignals()
	}
}

func (kernel *Kernel) ResetStats() {
}
