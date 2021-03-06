package cpu

import (
	"github.com/mcai/heo/cpu/uncore"
	"github.com/mcai/heo/noc"
	"github.com/mcai/heo/simutil"
	"time"
)

type CPUExperiment struct {
	CPUConfig    *CPUConfig
	UncoreConfig *uncore.UncoreConfig
	NocConfig    *noc.NoCConfig

	cycleAccurateEventQueue *simutil.CycleAccurateEventQueue
	blockingEventDispatcher *simutil.BlockingEventDispatcher

	ISA *ISA

	Kernel    *Kernel
	Processor *Processor

	MemoryHierarchy uncore.MemoryHierarchy
	OoO             *OoO

	BeginTime, EndTime time.Time

	Stats   simutil.Stats
	statMap map[string]interface{}

	L2PrefetchRequestProfiler *L2PrefetchRequestProfiler

	L2CacheRequestTracer *L2CacheRequestTracer
}

type CPUExperimentStoppedEvent struct {
}

func NewCPUExperiment(config *CPUConfig) *CPUExperiment {
	var experiment = &CPUExperiment{
		CPUConfig:    config,
		UncoreConfig: uncore.NewUncoreConfig(config.NumCores, config.NumThreadsPerCore),
		NocConfig:    noc.NewNoCConfig(config.OutputDirectory, -1, -1, config.NetworkType, -1, false),
	}

	experiment.ISA = NewISA()

	experiment.Kernel = NewKernel(experiment)

	experiment.cycleAccurateEventQueue = simutil.NewCycleAccurateEventQueue()
	experiment.blockingEventDispatcher = simutil.NewBlockingEventDispatcher()

	experiment.Processor = NewProcessor(experiment)

	experiment.MemoryHierarchy = uncore.NewBaseMemoryHierarchy(experiment, experiment.UncoreConfig, experiment.NocConfig)
	experiment.OoO = NewOoO(experiment)

	experiment.Kernel.LoadContexts()

	experiment.Processor.UpdateContextToThreadAssignments()

	experiment.L2PrefetchRequestProfiler = NewL2PrefetchRequestProfiler(experiment)

	if config.TraceL2Requests {
		experiment.L2CacheRequestTracer = NewL2CacheRequestTracer(experiment, config.OutputDirectory + "/l2_requests_trace.txt")
	}

	return experiment
}

func (experiment *CPUExperiment) CycleAccurateEventQueue() *simutil.CycleAccurateEventQueue {
	return experiment.cycleAccurateEventQueue
}

func (experiment *CPUExperiment) BlockingEventDispatcher() *simutil.BlockingEventDispatcher {
	return experiment.blockingEventDispatcher
}

func (experiment *CPUExperiment) Run() {
	experiment.dumpConfigs()

	experiment.BeginTime = time.Now()

	experiment.doFastForward()

	experiment.EndTime = time.Now()

	experiment.dumpStats("fastforward")

	experiment.ResetStats()

	experiment.BeginTime = time.Now()

	experiment.doMeasurement()

	experiment.EndTime = time.Now()

	experiment.blockingEventDispatcher.Dispatch(&CPUExperimentStoppedEvent{})

	experiment.dumpStats("measurement")
}

func (experiment *CPUExperiment) dumpConfigs() {
	experiment.CPUConfig.Dump(experiment.CPUConfig.OutputDirectory)
	experiment.MemoryHierarchy.Config().Dump(experiment.CPUConfig.OutputDirectory)
	experiment.MemoryHierarchy.Network().Config().Dump(experiment.CPUConfig.OutputDirectory)
}

func (experiment *CPUExperiment) canDoFastForwardOneCycle() bool {
	return experiment.CPUConfig.MaxFastForwardDynamicInsts == -1 ||
		experiment.Processor.Cores[0].Threads()[0].NumDynamicInsts() < experiment.CPUConfig.MaxFastForwardDynamicInsts
}

func (experiment *CPUExperiment) canDoMeasurementOneCycle() bool {
	return experiment.CPUConfig.MaxMeasurementDynamicInsts == -1 ||
		experiment.Processor.Cores[0].Threads()[0].NumDynamicInsts() < experiment.CPUConfig.MaxMeasurementDynamicInsts
}

func (experiment *CPUExperiment) advanceOneCycle() {
	experiment.Kernel.AdvanceOneCycle()
	experiment.Processor.UpdateContextToThreadAssignments()

	experiment.cycleAccurateEventQueue.AdvanceOneCycle()
}

func (experiment *CPUExperiment) doFastForward() {
	for len(experiment.Kernel.Contexts) > 0 && experiment.canDoFastForwardOneCycle() {
		for _, core := range experiment.Processor.Cores {
			core.FastForwardOneCycle()
		}

		experiment.advanceOneCycle()
	}
}

func (experiment *CPUExperiment) doMeasurement() {
	for len(experiment.Kernel.Contexts) > 0 && experiment.canDoMeasurementOneCycle() {
		for _, core := range experiment.Processor.Cores {
			core.(*OoOCore).MeasurementOneCycle()
		}

		experiment.advanceOneCycle()
	}
}

func (experiment *CPUExperiment) SimulationTime() time.Duration {
	return experiment.EndTime.Sub(experiment.BeginTime)
}

func (experiment *CPUExperiment) CyclesPerSecond() float64 {
	return float64(experiment.CycleAccurateEventQueue().CurrentCycle) / experiment.EndTime.Sub(experiment.BeginTime).Seconds()
}

func (experiment *CPUExperiment) InstructionsPerSecond() float64 {
	return float64(experiment.Processor.NumDynamicInsts()) / experiment.EndTime.Sub(experiment.BeginTime).Seconds()
}
