package main

import (
	"github.com/mcai/heo/noc"
	"github.com/mcai/heo/simutil"
)

func NewTraceDrivenExperiment(outputDirectory string, numNodes int, maxCycles int64, traffic noc.TrafficType, dataPacketInjectionRate float64, routing noc.RoutingType, selection noc.SelectionType, antPacketInjectionRate float64, acoSelectionAlpha float64, reinforcementFactor float64, traceFileNames string) simutil.Experiment {
	var config = noc.NewNoCConfig(
		outputDirectory,
		numNodes,
		maxCycles,
		-1,
		false)

	config.DataPacketTraffic = traffic
	config.DataPacketInjectionRate = dataPacketInjectionRate
	config.Routing = routing
	config.Selection = selection

	if selection == noc.SELECTION_ACO {
		config.AntPacketInjectionRate = antPacketInjectionRate
		config.AcoSelectionAlpha = acoSelectionAlpha
		config.ReinforcementFactor = reinforcementFactor
	}

	config.TraceFileName = traceFileNames

	return noc.NewNoCExperiment(config)
}
