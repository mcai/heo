package cpu

type Processor struct {
	Experiment              *CPUExperiment
	Cores                   []Core
	ContextToThreadMappings map[*Context]Thread
}

func NewProcessor(experiment *CPUExperiment) *Processor {
	var processor = &Processor{
		Experiment:              experiment,
		ContextToThreadMappings: make(map[*Context]Thread),
	}

	for i := int32(0); i < experiment.CPUConfig.NumCores; i++ {
		var core = NewOoOCore(processor, i)

		for j := int32(0); j < experiment.CPUConfig.NumThreadsPerCore; j++ {
			var thread = NewOoOThread(core, j)
			core.AddThread(thread)
		}

		processor.Cores = append(processor.Cores, core)
	}

	return processor
}

func (processor *Processor) UpdateContextToThreadAssignments() {
	var contextsToReserve []*Context

	for _, context := range processor.Experiment.Kernel.Contexts {
		if context.ThreadId != -1 && processor.ContextToThreadMappings[context] == nil {
			if context.State == ContextState_IDLE {
				context.State = ContextState_RUNNING
			}

			var coreNum = context.ThreadId / processor.Experiment.CPUConfig.NumThreadsPerCore
			var threadNum = context.ThreadId % processor.Experiment.CPUConfig.NumThreadsPerCore

			var candidateThread = processor.Cores[coreNum].Threads()[threadNum]

			processor.ContextToThreadMappings[context] = candidateThread

			candidateThread.SetContext(context)

			if oooThread := candidateThread.(*OoOThread); oooThread != nil {
				oooThread.UpdateFetchNpcAndNnpcFromRegs()
			}

			contextsToReserve = append(contextsToReserve, context)
		} else if context.State == ContextState_FINISHED {
			var thread = processor.ContextToThreadMappings[context].(*OoOThread)

			if thread != nil && thread.IsLastDecodedDynamicInstCommitted() && thread.ReorderBuffer.Empty() {
				processor.kill(context)
			}
		} else {
			contextsToReserve = append(contextsToReserve, context)
		}
	}

	processor.Experiment.Kernel.Contexts = contextsToReserve
}

func (processor *Processor) NumZombies() int32 {
	var numZombies = int32(0)

	for _, context := range processor.Experiment.Kernel.Contexts {
		if context.State == ContextState_FINISHED {
			numZombies++
		}
	}

	return numZombies
}

func (processor *Processor) kill(context *Context) {
	if context.State != ContextState_FINISHED {
		panic("Impossible")
	}

	for _, c := range processor.Experiment.Kernel.Contexts {
		if c.Parent == context {
			processor.kill(c)
		}
	}

	if context.Parent == nil {
		context.Process.CloseProgram()
	}

	processor.ContextToThreadMappings[context].SetContext(nil)

	context.ThreadId = -1
}

func (processor *Processor) NumDynamicInsts() int64 {
	var numDynamicInsts = int64(0)

	for _, core := range processor.Cores {
		numDynamicInsts += core.NumDynamicInsts()
	}

	return numDynamicInsts
}

func (processor *Processor) InstructionsPerCycle() float64 {
	if processor.Experiment.CycleAccurateEventQueue().CurrentCycle == 0 {
		return float64(0)
	}

	return float64(processor.NumDynamicInsts()) / float64(processor.Experiment.CycleAccurateEventQueue().CurrentCycle)
}

func (processor *Processor) CyclesPerInstructions() float64 {
	if processor.NumDynamicInsts() == 0 {
		return float64(0)
	}

	return float64(processor.Experiment.CycleAccurateEventQueue().CurrentCycle) / float64(processor.NumDynamicInsts())
}

func (processor *Processor) ResetStats() {
	for _, core := range processor.Cores {
		core.ResetStats()
	}
}
