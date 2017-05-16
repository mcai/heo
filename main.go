package main

import (
	"flag"
	"github.com/mcai/heo/noc"
	"github.com/mcai/heo/simutil"
)

func main() {
	var outputDirectory string
	var benchmark string
	var traceFileName string
	var numNodes int
	var routing string
	var selection string
	var maxCycles int64

	flag.StringVar(&outputDirectory, "d", "", "output directory")
	flag.StringVar(&benchmark, "b", "", "benchmark")
	flag.StringVar(&traceFileName, "f", "", "NOC simulation trace file name")
	flag.IntVar(&numNodes, "n", 16, "number of NOC nodes")
	flag.StringVar(&routing, "r", "OddEven", "NOC routing algorithm")
	flag.StringVar(&selection, "s", "BufferLevel", "NOC selection algorithm")
	flag.Int64Var(&maxCycles, "c", 1000, "Maximum number of cycles to simulate")

	flag.Parse()

	var dataPacketInjectionRate = 0.015
	var antPacketInjectionRate = 0.0002

	var acoSelectionAlpha = 0.45
	var reinforcementFactor = 0.001

	var experiment = NewTraceDrivenExperiment(
		outputDirectory,
		numNodes,
		maxCycles,
		noc.TRAFFIC_TRACE,
		dataPacketInjectionRate,
		noc.ROUTING_ODD_EVEN, noc.SELECTION_BUFFER_LEVEL,
		antPacketInjectionRate,
		acoSelectionAlpha,
		reinforcementFactor,
		traceFileName,
	)

	var experiments []simutil.Experiment

	experiments = append(experiments, experiment)

	simutil.RunExperiments(experiments, false)
}
