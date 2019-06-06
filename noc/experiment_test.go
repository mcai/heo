package noc

import "testing"

func TestNoCExperiment(t *testing.T) {
	var numNodes = 64
	var maxCycles = int64(10000)
	var maxPackets = int64(-1)
	var drainPackets = true

	var config = NewNoCConfig("test_results/synthetic/aco", numNodes, maxCycles, NetworkType_FIXED_LATENCY, maxPackets, drainPackets)

	config.Routing = RoutingOddEven
	config.Selection = SelectionAco

	config.DataPacketTraffic = TrafficTranspose1
	config.DataPacketInjectionRate = 0.06

	config.AntPacketTraffic = TrafficUniform
	config.AntPacketInjectionRate = 0.0002

	config.AcoSelectionAlpha = 0.45
	config.ReinforcementFactor = 0.001

	var experiment = NewNoCExperiment(config)

	experiment.Run()
}
