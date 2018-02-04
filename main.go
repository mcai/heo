package main

import (
	"flag"
	"github.com/mcai/heo/noc"
)

func main() {
	var outputDirectory string

	var traffic string
	var traceFileName string

	var maxCycles int64
	var numNodes int

	var routing string
	var selection string

	var dataPacketInjectionRate float64
	var antPacketInjectionRate float64

	var acoSelectionAlpha float64
	var reinforcementFactor float64

	flag.StringVar(&outputDirectory, "d", "results", "output directory")

	flag.StringVar(&traffic, "t", "Uniform", "traffic")
	flag.StringVar(&traceFileName, "tf", "", "NoC simulation trace file name")

	flag.Int64Var(&maxCycles, "c", 1000, "maximum number of cycles to simulate")
	flag.IntVar(&numNodes, "n", 16, "number of NoC nodes")

	flag.StringVar(&routing, "r", "OddEven", "NoC routing algorithm")
	flag.StringVar(&selection, "s", "BufferLevel", "NoC selection algorithm")

	flag.Float64Var(&dataPacketInjectionRate, "di", 0.015, "data packet injection rate")
	flag.Float64Var(&antPacketInjectionRate, "ai", 0.0002, "ant packet injection rate in the ACO NoC selection algorithm")

	flag.Float64Var(&acoSelectionAlpha, "sa", 0.45, "ACO selection alpha in the ACO NoC selection algorithm")
	flag.Float64Var(&reinforcementFactor, "rf", 0.001, "reinforcement factor in the ACO NoC selection algorithm")

	flag.Parse()

	var config = noc.NewNoCConfig(
		outputDirectory,
		numNodes,
		maxCycles,
		-1,
		false)

	config.DataPacketTraffic = noc.TrafficType(traffic)
	config.DataPacketInjectionRate = dataPacketInjectionRate
	config.Routing = noc.RoutingType(routing)
	config.Selection = noc.SelectionType(selection)

	if noc.SelectionType(selection) == noc.SelectionAco {
		config.AntPacketInjectionRate = antPacketInjectionRate
		config.AcoSelectionAlpha = acoSelectionAlpha
		config.ReinforcementFactor = reinforcementFactor
	}

	config.TraceFileName = traceFileName

	noc.NewNoCExperiment(config).Run()
}
