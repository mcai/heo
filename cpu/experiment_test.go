package cpu

import (
	"github.com/mcai/heo/noc"
	"testing"
)

func TestMstBaseline(t *testing.T) {
	var config = NewCPUConfig("../test_results/real/mst_baseline")

	config.ContextMappings = append(config.ContextMappings,
		NewContextMapping(0, "../benchmarks/Olden_Custom1/mst/baseline/mst.mips", "200"))

	config.NumCores = 2
	config.NumThreadsPerCore = 2
	config.MaxFastForwardDynamicInsts = int64(-1)
	config.MaxMeasurementDynamicInsts = int64(0)

	config.TraceL2Requests = true

	//config.NetworkType = noc.NetworkType_BASE;
	config.NetworkType = noc.NetworkType_FIXED_LATENCY

	var experiment = NewCPUExperiment(config)
	experiment.Run()
}

func TestMstHelperThreaded(t *testing.T) {
	var config = NewCPUConfig("../test_results/real/mst_ht")

	config.ContextMappings = append(config.ContextMappings,
		NewContextMapping(0, "../benchmarks/Olden_Custom1/mst/ht/mst.mips", "200"))

	config.NumCores = 2
	config.NumThreadsPerCore = 2
	config.MaxFastForwardDynamicInsts = int64(-1)
	config.MaxMeasurementDynamicInsts = int64(0)

	config.TraceL2Requests = true

	//config.NetworkType = noc.NetworkType_BASE;
	config.NetworkType = noc.NetworkType_FIXED_LATENCY

	var experiment = NewCPUExperiment(config)
	experiment.Run()
}