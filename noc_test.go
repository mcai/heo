package main

import (
	"fmt"
	"os/exec"
	"os"
	"sync"
	"github.com/mcai/heo/noc"
	"github.com/mcai/heo/simutil"
)

var (
	numNodes int = 8 * 8
	maxCycles int64 = int64(20000)
	maxPackets int64 = int64(-1)
	drainPackets = false
)

func NewExperiment(outputDirectoryPrefix string, traffic noc.TrafficType, dataPacketInjectionRate float64, routing noc.RoutingType, selection noc.SelectionType, antPacketInjectionRate float64, acoSelectionAlpha float64, reinforcementFactor float64) simutil.Experiment {
	var outputDirectory string

	switch {
	case selection == noc.SELECTION_ACO:
		outputDirectory = fmt.Sprintf("results/%s/t_%s/j_%f/r_%s/s_%s/aj_%f/a_%f/rf_%f/",
			outputDirectoryPrefix, traffic, dataPacketInjectionRate, routing, selection, antPacketInjectionRate, acoSelectionAlpha, reinforcementFactor)
	default:
		outputDirectory = fmt.Sprintf("results/%s/t_%s/j_%f/r_%s/s_%s/",
			outputDirectoryPrefix, traffic, dataPacketInjectionRate, routing, selection)
	}

	var config = noc.NewNoCConfig(
		outputDirectory,
		numNodes,
		maxCycles,
		maxPackets,
		drainPackets)

	config.DataPacketTraffic = traffic
	config.DataPacketInjectionRate = dataPacketInjectionRate
	config.Routing = routing
	config.Selection = selection

	if selection == noc.SELECTION_ACO {
		config.AntPacketInjectionRate = antPacketInjectionRate
		config.AcoSelectionAlpha = acoSelectionAlpha
		config.ReinforcementFactor = reinforcementFactor
	}

	return noc.NewNoCExperiment(config)
}

type NoCRoutingSolution struct {
	Routing   noc.RoutingType
	Selection noc.SelectionType
}

func TestTrafficsAndDataPacketInjectionRates() map[noc.TrafficType]([]simutil.Experiment) {
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

	var experiments = make(map[noc.TrafficType]([]simutil.Experiment))

	for _, traffic := range noc.TRAFFICS {
		for _, dataPacketInjectionRate := range dataPacketInjectionRates {
			var nocRoutingSolutions = []NoCRoutingSolution{
				{noc.ROUTING_XY, noc.SELECTION_BUFFER_LEVEL},
				{noc.ROUTING_NEGATIVE_FIRST, noc.SELECTION_BUFFER_LEVEL},
				{noc.ROUTING_WEST_FIRST, noc.SELECTION_BUFFER_LEVEL},
				{noc.ROUTING_NORTH_LAST, noc.SELECTION_BUFFER_LEVEL},
				{noc.ROUTING_NORTH_LAST, noc.SELECTION_ACO},
				{noc.ROUTING_ODD_EVEN, noc.SELECTION_RANDOM},
				{noc.ROUTING_ODD_EVEN, noc.SELECTION_BUFFER_LEVEL},
				{noc.ROUTING_ODD_EVEN, noc.SELECTION_ACO},
			}

			for _, nocRoutingSolution := range nocRoutingSolutions {
				experiments[traffic] = append(
					experiments[traffic],
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

	return experiments
}

func TestAntPacketInjectionRates() []simutil.Experiment {
	var traffic = noc.TRAFFIC_TRANSPOSE1
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

	var experiments []simutil.Experiment

	for _, antPacketInjectionRate := range antPacketInjectionRates {
		experiments = append(
			experiments,
			NewExperiment(
				outputDirectoryPrefix,
				traffic,
				dataPacketInjectionRate,
				noc.ROUTING_ODD_EVEN,
				noc.SELECTION_ACO,
				antPacketInjectionRate,
				acoSelectionAlpha,
				reinforcementFactor))
	}

	return experiments
}

func TestAcoSelectionAlphasAndReinforcementFactors() []simutil.Experiment {
	var traffic = noc.TRAFFIC_TRANSPOSE1
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

	var experiments []simutil.Experiment

	for _, acoSelectionAlpha := range acoSelectionAlphas {
		for _, reinforcementFactor := range reinforcementFactors {
			experiments = append(
				experiments,
				NewExperiment(
					outputDirectoryPrefix,
					traffic,
					dataPacketInjectionRate,
					noc.ROUTING_ODD_EVEN,
					noc.SELECTION_ACO,
					antPacketInjectionRate,
					acoSelectionAlpha,
					reinforcementFactor))
		}
	}

	return experiments
}

var (
	TrafficsAndDataPacketInjectionRates = TestTrafficsAndDataPacketInjectionRates()
	AntPacketInjectionRates = TestAntPacketInjectionRates()
	AcoSelectionAlphasAndReinforcementFactors = TestAcoSelectionAlphasAndReinforcementFactors()
)

func run() {
	var experiments []simutil.Experiment

	for _, traffic := range noc.TRAFFICS {
		for _, experiment := range TrafficsAndDataPacketInjectionRates[traffic] {
			experiments = append(experiments, experiment)
		}
	}

	for _, experiment := range AntPacketInjectionRates {
		experiments = append(experiments, experiment)
	}

	for _, experiment := range AcoSelectionAlphasAndReinforcementFactors {
		experiments = append(experiments, experiment)
	}

	simutil.RunExperiments(experiments, true)
}

func analyze() {
	var wg = &sync.WaitGroup{}

	for _, traffic := range noc.TRAFFICS {
		outputDirectory := fmt.Sprintf("results/trafficsAndDataPacketInjectionRates/t_%s", traffic)

		noc.WriteCSVFile(outputDirectory, "result.csv", TrafficsAndDataPacketInjectionRates[traffic], noc.GetCSVFields())

		wg.Add(1)
		go generatePlot(
			wg,
			outputDirectory,
			"result.csv",
			"throughput.pdf",
			"Data_Packet_Injection_Rate_(packets/cycle/node)",
			"NoC_Routing_Solution",
			"Payload_Throughput_(packets/cycle/node)",
			90,
			BAR_PLOT,
		)

		wg.Add(1)
		go generatePlot(
			wg,
			outputDirectory,
			"result.csv",
			"average_packet_delay.pdf",
			"Data_Packet_Injection_Rate_(packets/cycle/node)",
			"NoC_Routing_Solution",
			"Avg._Payload_Packet_Delay_(cycles)",
			90,
			BAR_PLOT,
		)
	}

	outputDirectory := "results/antPacketInjectionRates"

	noc.WriteCSVFile(outputDirectory, "result.csv", AntPacketInjectionRates, noc.GetCSVFields())

	wg.Add(1)
	go generatePlot(
		wg,
		outputDirectory,
		"result.csv",
		"throughput.pdf",
		"Ant_Packet_Injection_Rate_(packets/cycle/node)",
		"",
		"Payload_Throughput_(packets/cycle/node)",
		90,
		BAR_PLOT,
	)

	wg.Add(1)
	go generatePlot(
		wg,
		outputDirectory,
		"result.csv",
		"average_packet_delay.pdf",
		"Ant_Packet_Injection_Rate_(packets/cycle/node)",
		"",
		"Avg._Payload_Packet_Delay_(cycles)",
		90,
		BAR_PLOT,
	)

	outputDirectory = "results/acoSelectionAlphasAndReinforcementFactors"

	noc.WriteCSVFile(outputDirectory, "result.csv", AcoSelectionAlphasAndReinforcementFactors, noc.GetCSVFields())

	wg.Add(1)
	go generatePlot(
		wg,
		outputDirectory,
		"result.csv",
		"throughput.pdf",
		"Alpha",
		"Reinforcement_Factor",
		"Payload_Throughput_(packets/cycle/node)",
		90,
		BAR_PLOT,
	)

	wg.Add(1)
	go generatePlot(
		wg,
		outputDirectory,
		"result.csv",
		"average_packet_delay.pdf",
		"Alpha",
		"Reinforcement_Factor",
		"Avg._Payload_Packet_Delay_(cycles)",
		90,
		BAR_PLOT,
	)

	wg.Wait()
}

const (
	BAR_PLOT = 0
	LINE_PLOT = 1
)

type PlotType int

func generatePlot(wg *sync.WaitGroup, outputDirectory string, csvFileName string, plotFileName string, x string, hue string, y string, xticklabelsRotation int, plotType PlotType) {
	var cmd *exec.Cmd

	switch hue {
	case "":
		cmd = exec.Command(
			"./plots.sh",
			"--csv_file_name", outputDirectory + "/" + csvFileName,
			"--plot_file_name", outputDirectory + "/" + plotFileName,
			"--x", x,
			"--y", y,
			"--xticklabels_rotation", fmt.Sprintf("%d", xticklabelsRotation),
			"--plot_type", fmt.Sprintf("%d", plotType),
		)
	default:
		cmd = exec.Command(
			"./plots.sh",
			"--csv_file_name", outputDirectory + "/" + csvFileName,
			"--plot_file_name", outputDirectory + "/" + plotFileName,
			"--x", x,
			"--hue", hue,
			"--y", y,
			"--xticklabels_rotation", fmt.Sprintf("%d", xticklabelsRotation),
			"--plot_type", fmt.Sprintf("%d", plotType),
		)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	wg.Done()
}

func testNoc() {
	run()
	analyze()
}
