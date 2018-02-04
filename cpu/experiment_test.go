package cpu

import (
	"testing"
)

func TestCPUExperiment(t *testing.T) {
	var config = NewCPUConfig("test_results/real/mst_ht_100")

	config.ContextMappings = append(config.ContextMappings,
		NewContextMapping(0, "/home/itecgo/Projects/Archimulator/benchmarks/Olden_Custom1/mst/ht/mst.mips", "100"))

	config.MaxMeasurementDynamicInsts = 1000000

	var experiment = NewCPUExperiment(config)

	experiment.Run(false)
}
