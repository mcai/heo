package noc

import "github.com/mcai/heo/simutil"

type TrafficType string

const (
	TrafficUniform    = TrafficType("Uniform")
	TrafficTranspose1 = TrafficType("Transpose1")
	TrafficTranspose2 = TrafficType("Transpose2")
	TrafficTrace      = TrafficType("Trace")
)

var TRAFFICS = []TrafficType{
	TrafficUniform,
	TrafficTranspose1,
	TrafficTranspose2,
	TrafficTrace,
}

type RoutingType string

const (
	RoutingXY            = RoutingType("XY")
	RoutingNegativeFirst = RoutingType("NegativeFirst")
	RoutingWestFirst     = RoutingType("WestFirst")
	RoutingNorthLast     = RoutingType("NorthLast")
	RoutingOddEven       = RoutingType("OddEven")
)

var ROUTINGS = []RoutingType{
	RoutingXY,
	RoutingNegativeFirst,
	RoutingWestFirst,
	RoutingNorthLast,
	RoutingOddEven,
}

type SelectionType string

const (
	SelectionRandom      = SelectionType("Random")
	SelectionBufferLevel = SelectionType("BufferLevel")
	SelectionAco         = SelectionType("ACO")
)

var SELECTIONS = []SelectionType{
	SelectionRandom,
	SelectionBufferLevel,
	SelectionAco,
}

type NoCConfig struct {
	OutputDirectory string

	NumNodes int

	MaxCycles int64

	MaxPackets int64

	DrainPackets bool

	Routing RoutingType

	Selection SelectionType

	MaxInjectionBufferSize int

	MaxInputBufferSize int

	NumVirtualChannels int

	LinkWidth int
	LinkDelay int

	DataPacketTraffic       TrafficType
	DataPacketInjectionRate float64
	DataPacketSize          int

	AntPacketTraffic       TrafficType
	AntPacketInjectionRate float64
	AntPacketSize          int

	AcoSelectionAlpha   float64
	ReinforcementFactor float64

	TraceFileName string
}

func NewNoCConfig(outputDirectory string, numNodes int, maxCycles int64, maxPackets int64, drainPackets bool) *NoCConfig {
	var nocConfig = &NoCConfig{
		OutputDirectory: outputDirectory,

		NumNodes: numNodes,

		MaxCycles: maxCycles,

		MaxPackets: maxPackets,

		DrainPackets: drainPackets,

		Routing: RoutingOddEven,

		Selection: SelectionBufferLevel,

		MaxInjectionBufferSize: 32,

		MaxInputBufferSize: 4,

		NumVirtualChannels: 4,

		LinkWidth: 4,
		LinkDelay: 1,

		DataPacketTraffic:       TrafficTranspose1,
		DataPacketInjectionRate: 0.01,
		DataPacketSize:          16,

		AntPacketTraffic:       TrafficUniform,
		AntPacketInjectionRate: 0.01,
		AntPacketSize:          4,

		AcoSelectionAlpha:   0.5,
		ReinforcementFactor: 0.05,
	}

	return nocConfig
}

func (nocConfig *NoCConfig) Dump(outputDirectory string) {
	simutil.WriteJsonFile(nocConfig, outputDirectory, simutil.NOC_CONFIG_JSON_FILE_NAME)
}
