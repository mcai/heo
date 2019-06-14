package noc

import (
	"fmt"
	"math"
)

type BaseNetwork struct {
	driver NetworkDriver
	config *NoCConfig

	currentPacketId int64
	NumNodes        int
	Nodes           []*Node
	Width           int
	AcceptPacket    bool
	trafficSources  []TrafficSource

	NumPacketsReceived    int64
	NumPacketsTransmitted int64

	totalPacketDelays int64
	MaxPacketDelay    int

	totalPacketHops int64
	MaxPacketHops   int

	NumPayloadPacketsReceived    int64
	NumPayloadPacketsTransmitted int64

	totalPayloadPacketDelays int64
	MaxPayloadPacketDelay    int

	totalPayloadPacketHops int64
	MaxPayloadPacketHops   int

	numFlitPerStateDelaySamples map[FlitState]int64
	totalFlitPerStateDelays     map[FlitState]int64
	MaxFlitPerStateDelay        map[FlitState]int
}

func NewBaseNetwork(driver NetworkDriver, config *NoCConfig) *BaseNetwork {
	var baseNetwork = &BaseNetwork{
		driver:       driver,
		config:       config,
		NumNodes:     config.NumNodes,
		Width:        int(math.Sqrt(float64(config.NumNodes))),
		AcceptPacket: true,

		numFlitPerStateDelaySamples: make(map[FlitState]int64),
		totalFlitPerStateDelays:     make(map[FlitState]int64),
		MaxFlitPerStateDelay:        make(map[FlitState]int),
	}

	for i := 0; i < baseNetwork.NumNodes; i++ {
		var node = NewNode(baseNetwork, i)
		baseNetwork.Nodes = append(baseNetwork.Nodes, node)
	}

	switch selection := config.Selection; selection {
	case SelectionAco:
		switch antPacketTraffic := config.AntPacketTraffic; antPacketTraffic {
		case TrafficUniform:
			baseNetwork.AddTrafficSource(NewUniformTrafficSource(baseNetwork, config.AntPacketInjectionRate, int64(-1), func(src int, dest int) Packet {
				return NewAntPacket(baseNetwork, src, dest, config.AntPacketSize, func() {}, true)
			}))
		case TrafficTranspose1:
			baseNetwork.AddTrafficSource(NewTranspose1TrafficSource(baseNetwork, config.AntPacketInjectionRate, int64(-1), func(src int, dest int) Packet {
				return NewAntPacket(baseNetwork, src, dest, config.AntPacketSize, func() {}, true)
			}))
		case TrafficTranspose2:
			baseNetwork.AddTrafficSource(NewTranspose2TrafficSource(baseNetwork, config.AntPacketInjectionRate, int64(-1), func(src int, dest int) Packet {
				return NewAntPacket(baseNetwork, src, dest, config.AntPacketSize, func() {}, true)
			}))
		default:
			panic(fmt.Sprintf("ant packet traffic %s is not supported", antPacketTraffic))
		}
	}

	driver.CycleAccurateEventQueue().AddPerCycleEvent(func() {
		for _, node := range baseNetwork.Nodes {
			node.Router.AdvanceOneCycle()
		}
	})

	return baseNetwork
}

func (baseNetwork *BaseNetwork) CurrentPacketId() int64  {
	return baseNetwork.currentPacketId
}

func (baseNetwork *BaseNetwork) SetCurrentPacketId(currentPacketId int64)  {
	baseNetwork.currentPacketId = currentPacketId
}

func (baseNetwork *BaseNetwork) Driver() NetworkDriver {
	return baseNetwork.driver
}

func (baseNetwork *BaseNetwork) Config() *NoCConfig {
	return baseNetwork.config
}

func (baseNetwork *BaseNetwork) GetX(id int) int {
	return id % baseNetwork.Width
}

func (baseNetwork *BaseNetwork) GetY(id int) int {
	return (id - id%baseNetwork.Width) / baseNetwork.Width
}

func (baseNetwork *BaseNetwork) TrafficSources() []TrafficSource {
	return baseNetwork.trafficSources
}

func (baseNetwork *BaseNetwork) AddTrafficSource(trafficSource TrafficSource) {
	baseNetwork.trafficSources = append(baseNetwork.trafficSources, trafficSource)

	baseNetwork.Driver().CycleAccurateEventQueue().AddPerCycleEvent(func() {
		trafficSource.AdvanceOneCycle()
	})
}

func (baseNetwork *BaseNetwork) Receive(packet Packet) bool {
	if !baseNetwork.Nodes[packet.Src()].Router.InjectPacket(packet) {
		baseNetwork.Driver().CycleAccurateEventQueue().Schedule(func() {
			baseNetwork.Receive(packet)
		}, 1)
		return false
	}

	baseNetwork.LogPacketReceived(packet)

	return true
}

func (baseNetwork *BaseNetwork) LogPacketReceived(packet Packet) {
	baseNetwork.NumPacketsReceived++

	if packet.HasPayload() {
		baseNetwork.NumPayloadPacketsReceived++
	}
}

func (baseNetwork *BaseNetwork) LogPacketTransmitted(packet Packet) {
	baseNetwork.NumPacketsTransmitted++

	if packet.HasPayload() {
		baseNetwork.NumPayloadPacketsTransmitted++
	}

	baseNetwork.totalPacketDelays += int64(Delay(packet))
	baseNetwork.totalPacketHops += int64(Hops(packet))

	if packet.HasPayload() {
		baseNetwork.totalPayloadPacketDelays += int64(Delay(packet))
		baseNetwork.totalPayloadPacketHops += int64(Hops(packet))
	}

	baseNetwork.MaxPacketDelay = int(math.Max(float64(baseNetwork.MaxPacketDelay), float64(Delay(packet))))
	baseNetwork.MaxPacketHops = int(math.Max(float64(baseNetwork.MaxPacketHops), float64(Hops(packet))))

	if packet.HasPayload() {
		baseNetwork.MaxPayloadPacketDelay = int(math.Max(float64(baseNetwork.MaxPayloadPacketDelay), float64(Delay(packet))))
		baseNetwork.MaxPayloadPacketHops = int(math.Max(float64(baseNetwork.MaxPayloadPacketHops), float64(Hops(packet))))
	}
}

func (baseNetwork *BaseNetwork) LogFlitPerStateDelay(state FlitState, delay int) {
	if _, exists := baseNetwork.numFlitPerStateDelaySamples[state]; !exists {
		baseNetwork.numFlitPerStateDelaySamples[state] = int64(0)
	}

	baseNetwork.numFlitPerStateDelaySamples[state]++

	if _, exists := baseNetwork.totalFlitPerStateDelays[state]; !exists {
		baseNetwork.totalFlitPerStateDelays[state] = int64(0)
	}

	baseNetwork.totalFlitPerStateDelays[state] += int64(delay)

	if _, exists := baseNetwork.MaxFlitPerStateDelay[state]; !exists {
		baseNetwork.MaxFlitPerStateDelay[state] = 0
	}

	baseNetwork.MaxFlitPerStateDelay[state] = int(math.Max(float64(baseNetwork.MaxFlitPerStateDelay[state]), float64(delay)))
}

func (baseNetwork *BaseNetwork) Throughput() float64 {
	if baseNetwork.Driver().CycleAccurateEventQueue().CurrentCycle == 0 {
		return 0.0
	}

	return float64(baseNetwork.NumPacketsTransmitted) / float64(baseNetwork.Driver().CycleAccurateEventQueue().CurrentCycle) / float64(baseNetwork.NumNodes)
}

func (baseNetwork *BaseNetwork) AveragePacketDelay() float64 {
	if baseNetwork.NumPacketsTransmitted == 0 {
		return 0.0
	}

	return float64(baseNetwork.totalPacketDelays) / float64(baseNetwork.NumPacketsTransmitted)
}

func (baseNetwork *BaseNetwork) AveragePacketHops() float64 {
	if baseNetwork.NumPacketsTransmitted == 0 {
		return 0.0
	}

	return float64(baseNetwork.totalPacketHops) / float64(baseNetwork.NumPacketsTransmitted)
}

func (baseNetwork *BaseNetwork) PayloadThroughput() float64 {
	if baseNetwork.Driver().CycleAccurateEventQueue().CurrentCycle == 0 {
		return 0.0
	}

	return float64(baseNetwork.NumPayloadPacketsTransmitted) / float64(baseNetwork.Driver().CycleAccurateEventQueue().CurrentCycle) / float64(baseNetwork.NumNodes)
}

func (baseNetwork *BaseNetwork) AveragePayloadPacketDelay() float64 {
	if baseNetwork.NumPayloadPacketsTransmitted == 0 {
		return 0.0
	}

	return float64(baseNetwork.totalPayloadPacketDelays) / float64(baseNetwork.NumPayloadPacketsTransmitted)
}

func (baseNetwork *BaseNetwork) AveragePayloadPacketHops() float64 {
	if baseNetwork.NumPayloadPacketsTransmitted == 0 {
		return 0.0
	}

	return float64(baseNetwork.totalPayloadPacketHops) / float64(baseNetwork.NumPayloadPacketsTransmitted)
}

func (baseNetwork *BaseNetwork) AverageFlitPerStateDelay(state FlitState) float64 {
	if _, exists := baseNetwork.numFlitPerStateDelaySamples[state]; exists {
		if baseNetwork.numFlitPerStateDelaySamples[state] > 0 {
			return float64(baseNetwork.totalFlitPerStateDelays[state]) / float64(baseNetwork.numFlitPerStateDelaySamples[state])
		}
	}

	return 0.0
}
