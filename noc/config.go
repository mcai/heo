package noc

import "github.com/mcai/heo/simutil"

type TrafficType string

const (
	TRAFFIC_UNIFORM = TrafficType("Uniform")
	TRAFFIC_TRANSPOSE1 = TrafficType("Transpose1")
	TRAFFIC_TRANSPOSE2 = TrafficType("Transpose2")
	TRAFFIC_TRACE = TrafficType("Trace")
)

var TRAFFICS = []TrafficType{
	TRAFFIC_UNIFORM,
	TRAFFIC_TRANSPOSE1,
	TRAFFIC_TRANSPOSE2,
	TRAFFIC_TRACE,
}

type RoutingType string

const (
	ROUTING_XY = RoutingType("XY")
	ROUTING_NEGATIVE_FIRST = RoutingType("NegativeFirst")
	ROUTING_WEST_FIRST = RoutingType("WestFirst")
	ROUTING_NORTH_LAST = RoutingType("NorthLast")
	ROUTING_ODD_EVEN = RoutingType("OddEven")
)

var ROUTINGS = []RoutingType{
	ROUTING_XY,
	ROUTING_NEGATIVE_FIRST,
	ROUTING_WEST_FIRST,
	ROUTING_NORTH_LAST,
	ROUTING_ODD_EVEN,
}

type SelectionType string

const (
	SELECTION_RANDOM = SelectionType("Random")
	SELECTION_BUFFER_LEVEL = SelectionType("BufferLevel")
	SELECTION_ACO = SelectionType("ACO")
)

var SELECTIONS = []SelectionType{
	SELECTION_RANDOM,
	SELECTION_BUFFER_LEVEL,
	SELECTION_ACO,
}

type NoCConfig struct {
	OutputDirectory         string

	NumNodes                int

	MaxCycles               int64

	MaxPackets              int64

	DrainPackets            bool

	Routing                 RoutingType

	Selection               SelectionType

	MaxInjectionBufferSize  int

	MaxInputBufferSize      int

	NumVirtualChannels      int

	LinkWidth               int
	LinkDelay               int

	DataPacketTraffic       TrafficType
	DataPacketInjectionRate float64
	DataPacketSize          int

	AntPacketTraffic        TrafficType
	AntPacketInjectionRate  float64
	AntPacketSize           int

	AcoSelectionAlpha       float64
	ReinforcementFactor     float64

	TraceFileName           string
}

func NewNoCConfig(outputDirectory string, numNodes int, maxCycles int64, maxPackets int64, drainPackets bool) *NoCConfig {
	var nocConfig = &NoCConfig{
		OutputDirectory:outputDirectory,

		NumNodes:numNodes,

		MaxCycles:maxCycles,

		MaxPackets:maxPackets,

		DrainPackets:drainPackets,

		Routing:ROUTING_ODD_EVEN,

		Selection:SELECTION_BUFFER_LEVEL,

		MaxInjectionBufferSize:32,

		MaxInputBufferSize:4,

		NumVirtualChannels:4,

		LinkWidth:4,
		LinkDelay:1,

		DataPacketTraffic:TRAFFIC_TRANSPOSE1,
		DataPacketInjectionRate:0.01,
		DataPacketSize:16,

		AntPacketTraffic:TRAFFIC_UNIFORM,
		AntPacketInjectionRate:0.01,
		AntPacketSize:4,

		AcoSelectionAlpha:0.5,
		ReinforcementFactor:0.05,
	}

	return nocConfig
}

func (nocConfig *NoCConfig) Dump(outputDirectory string) {
	simutil.WriteJsonFile(nocConfig, outputDirectory, simutil.NOC_CONFIG_JSON_FILE_NAME)
}
