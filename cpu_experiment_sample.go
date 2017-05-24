package main

import (
	"github.com/mcai/heo/cpu"
	"github.com/mcai/heo/simutil"
)

var (
	numCores = int32(2)
	numThreadsPerCore = int32(2)

	//maxFastForwardDynamicInsts = int64(100000000)
	//maxMeasurementDynamicInsts = int64(100000000)

	maxFastForwardDynamicInsts = int64(0)
	maxMeasurementDynamicInsts = int64(-1)

	//maxFastForwardDynamicInsts = int64(-1)
	//maxMeasurementDynamicInsts = int64(0)
)

func runExecutionDriven() {
	var experiments = []simutil.Experiment{
		mstBaseline(),
		mstHelperThreaded(),
	}

	simutil.RunExperiments(experiments, false)
}

func mstBaseline() simutil.Experiment {
	var config = cpu.NewCPUConfig("test_results/real/mst_baseline")

	config.ContextMappings = append(config.ContextMappings,
		cpu.NewContextMapping(0, "benchmarks/Olden_Custom1/mst/baseline/mst.mips", "1000"))

	config.NumCores = numCores
	config.NumThreadsPerCore = numThreadsPerCore
	config.MaxMeasurementDynamicInsts = maxMeasurementDynamicInsts
	config.MaxFastForwardDynamicInsts = maxFastForwardDynamicInsts

	return cpu.NewCPUExperiment(config)
}

func mstHelperThreaded() simutil.Experiment {
	var config = cpu.NewCPUConfig("test_results/real/mst_ht")

	config.ContextMappings = append(config.ContextMappings,
		cpu.NewContextMapping(0, "benchmarks/Olden_Custom1/mst/ht/mst.mips", "1000"))

	config.NumCores = numCores
	config.NumThreadsPerCore = numThreadsPerCore
	config.MaxMeasurementDynamicInsts = maxMeasurementDynamicInsts
	config.MaxFastForwardDynamicInsts = maxFastForwardDynamicInsts

	return cpu.NewCPUExperiment(config)
}
