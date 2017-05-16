package noc

import (
	"fmt"
	"testing"
	"github.com/mcai/heo/simutil"
)

var (
	numNodes int = 8 * 8
	maxCycles int64 = int64(20000)
	maxPackets int64 = int64(-1)
	drainPackets = false
)

func NewExperiment(outputDirectoryPrefix string, traffic TrafficType, dataPacketInjectionRate float64, routing RoutingType, selection SelectionType, antPacketInjectionRate float64, acoSelectionAlpha float64, reinforcementFactor float64) *NoCExperiment {
	var outputDirectory string

	switch {
	case selection == SELECTION_ACO:
		outputDirectory = fmt.Sprintf("results/%s/t_%s/j_%f/r_%s/s_%s/aj_%f/a_%f/rf_%f/",
			outputDirectoryPrefix, traffic, dataPacketInjectionRate, routing, selection, antPacketInjectionRate, acoSelectionAlpha, reinforcementFactor)
	default:
		outputDirectory = fmt.Sprintf("results/%s/t_%s/j_%f/r_%s/s_%s/",
			outputDirectoryPrefix, traffic, dataPacketInjectionRate, routing, selection)
	}

	var config = NewNoCConfig(
		outputDirectory,
		numNodes,
		maxCycles,
		maxPackets,
		drainPackets)

	config.DataPacketTraffic = traffic
	config.DataPacketInjectionRate = dataPacketInjectionRate
	config.Routing = routing
	config.Selection = selection

	if selection == SELECTION_ACO {
		config.AntPacketInjectionRate = antPacketInjectionRate
		config.AcoSelectionAlpha = acoSelectionAlpha
		config.ReinforcementFactor = reinforcementFactor
	}

	return NewNoCExperiment(config)
}

type NoCRoutingSolution struct {
	Routing   RoutingType
	Selection SelectionType
}

func TestTrafficsAndDataPacketInjectionRates(t *testing.T) {
	var dataPacketInjectionRates = []float64{
		0.015,
		0.030,
		0.045,
		0.060,
		0.075,
		0.090,
		0.105,
		0.120,
	}

	var antPacketInjectionRate = 0.0002

	var acoSelectionAlpha = 0.45
	var reinforcementFactor = 0.001

	var outputDirectoryPrefix = "trafficsAndDataPacketInjectionRates"

	var nocExperimentsPerTraffic = make(map[TrafficType]([]*NoCExperiment))

	for _, traffic := range TRAFFICS {
		for _, dataPacketInjectionRate := range dataPacketInjectionRates {
			var nocRoutingSolutions = []NoCRoutingSolution{
				{ROUTING_XY, SELECTION_BUFFER_LEVEL},
				{ROUTING_NEGATIVE_FIRST, SELECTION_BUFFER_LEVEL},
				{ROUTING_WEST_FIRST, SELECTION_BUFFER_LEVEL},
				{ROUTING_NORTH_LAST, SELECTION_BUFFER_LEVEL},
				{ROUTING_NORTH_LAST, SELECTION_ACO},
				{ROUTING_ODD_EVEN, SELECTION_RANDOM},
				{ROUTING_ODD_EVEN, SELECTION_BUFFER_LEVEL},
				{ROUTING_ODD_EVEN, SELECTION_ACO},
			}

			for _, nocRoutingSolution := range nocRoutingSolutions {
				nocExperimentsPerTraffic[traffic] = append(
					nocExperimentsPerTraffic[traffic],
					NewExperiment(
						outputDirectoryPrefix,
						traffic,
						dataPacketInjectionRate,
						nocRoutingSolution.Routing,
						nocRoutingSolution.Selection,
						antPacketInjectionRate,
						acoSelectionAlpha,
						reinforcementFactor))
			}
		}
	}

	var experiments []simutil.Experiment

	for _, traffic := range TRAFFICS {
		for _, nocExperiment := range nocExperimentsPerTraffic[traffic] {
			experiments = append(experiments, nocExperiment)
		}
	}

	simutil.RunExperiments(experiments, true)

	for _, traffic := range TRAFFICS {
		var outputDirectory = fmt.Sprintf("results/%s/t_%s", outputDirectoryPrefix, traffic)

		WriteCSVFile(outputDirectory, "result.csv", nocExperimentsPerTraffic[traffic], GetCSVFields())
	}
}

func TestAntPacketInjectionRates(t *testing.T) {
	var traffic = TRAFFIC_TRANSPOSE1
	var dataPacketInjectionRate = 0.060

	var antPacketInjectionRates = []float64{
		0.0001,
		0.0002,
		0.0004,
		0.0008,
		0.0016,
		0.0032,
		0.0064,
		0.0128,
		0.0256,
		0.0512,
		0.1024,
	}

	var acoSelectionAlpha = 0.45
	var reinforcementFactor = 0.001

	var outputDirectoryPrefix = "antPacketInjectionRates"

	var nocExperiments []*NoCExperiment

	for _, antPacketInjectionRate := range antPacketInjectionRates {
		nocExperiments = append(
			nocExperiments,
			NewExperiment(
				outputDirectoryPrefix,
				traffic,
				dataPacketInjectionRate,
				ROUTING_ODD_EVEN,
				SELECTION_ACO,
				antPacketInjectionRate,
				acoSelectionAlpha,
				reinforcementFactor))
	}

	var experiments []simutil.Experiment

	for _, nocExperiment := range nocExperiments {
		experiments = append(experiments, nocExperiment)
	}

	simutil.RunExperiments(experiments, true)

	var outputDirectory = fmt.Sprintf("results/%s", outputDirectoryPrefix)

	WriteCSVFile(outputDirectory, "result.csv", nocExperiments, GetCSVFields())
}

func TestAcoSelectionAlphasAndReinforcementFactors(t *testing.T) {
	var traffic = TRAFFIC_TRANSPOSE1
	var dataPacketInjectionRate = 0.060

	var antPacketInjectionRate = 0.0002

	var acoSelectionAlphas = []float64{
		0.30,
		0.35,
		0.40,
		0.45,
		0.50,
		0.55,
		0.60,
		0.65,
		0.70,
	}
	var reinforcementFactors = []float64{
		0.0005,
		0.001,
		0.002,
		0.004,
		0.008,
		0.016,
		0.032,
		0.064,
	}

	var outputDirectoryPrefix = "acoSelectionAlphasAndReinforcementFactors"

	var nocExperiments []*NoCExperiment

	for _, acoSelectionAlpha := range acoSelectionAlphas {
		for _, reinforcementFactor := range reinforcementFactors {
			nocExperiments = append(
				nocExperiments,
				NewExperiment(
					outputDirectoryPrefix,
					traffic,
					dataPacketInjectionRate,
					ROUTING_ODD_EVEN,
					SELECTION_ACO,
					antPacketInjectionRate,
					acoSelectionAlpha,
					reinforcementFactor))
		}
	}

	var experiments []simutil.Experiment

	for _, nocExperiment := range nocExperiments {
		experiments = append(experiments, nocExperiment)
	}

	simutil.RunExperiments(experiments, true)

	var outputDirectory = fmt.Sprintf("results/%s", outputDirectoryPrefix)

	WriteCSVFile(outputDirectory, "result.csv", nocExperiments, GetCSVFields())
}