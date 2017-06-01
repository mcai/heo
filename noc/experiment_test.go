package noc

import "testing"

func TestNoCExperiment(t *testing.T) {
	var numNodes = 64
	var maxCycles = int64(10000)
	var maxPackets = int64(-1)
	var drainPackets = true

	var config = NewNoCConfig("test_results/synthetic/aco", numNodes, maxCycles, maxPackets, drainPackets)

	config.Routing = ROUTING_ODD_EVEN
	config.Selection = SELECTION_ACO

	config.DataPacketTraffic = TRAFFIC_TRANSPOSE1
	config.DataPacketInjectionRate = 0.06

	config.AntPacketTraffic = TRAFFIC_UNIFORM
	config.AntPacketInjectionRate = 0.0002

	config.AcoSelectionAlpha = 0.45
	config.ReinforcementFactor = 0.001

	var experiment = NewNoCExperiment(config)

	experiment.Run()
}
