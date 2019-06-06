package cpu

import (
	"testing"
)

func TestMstBaseline(t *testing.T) {
	var config = NewCPUConfig("../test_results/real/mst_baseline")

	config.ContextMappings = append(config.ContextMappings,
		NewContextMapping(0, "../benchmarks/Olden_Custom1/mst/baseline/mst.mips", "100"))

	config.NumCores = 2
	config.NumThreadsPerCore = 2
	config.MaxFastForwardDynamicInsts = int64(10000)
	config.MaxMeasurementDynamicInsts = int64(100000000)

	var experiment = NewCPUExperiment(config)
	experiment.Run()
}

func TestMstHelperThreaded(t *testing.T) {
	var config = NewCPUConfig("../test_results/real/mst_ht")

	config.ContextMappings = append(config.ContextMappings,
		NewContextMapping(0, "../benchmarks/Olden_Custom1/mst/ht/mst.mips", "100"))

	config.NumCores = 2
	config.NumThreadsPerCore = 2
	config.MaxFastForwardDynamicInsts = int64(10000)
	config.MaxMeasurementDynamicInsts = int64(100000000)

	var experiment = NewCPUExperiment(config)
	experiment.Run()
}