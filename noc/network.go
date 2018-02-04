package noc

import (
	"math"
	"github.com/mcai/heo/simutil"
	"fmt"
)

type NetworkDriver interface {
	CycleAccurateEventQueue() *simutil.CycleAccurateEventQueue
}

type Network struct {
	Driver NetworkDriver
	Config *NoCConfig

	CurrentPacketId   int64
	NumNodes          int
	Nodes             []*Node
	Width             int
	AcceptPacket      bool
	trafficGenerators []TrafficGenerator

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

func NewNetwork(driver NetworkDriver, config *NoCConfig) *Network {
	var network = &Network{
		Driver:       driver,
		Config:       config,
		NumNodes:     config.NumNodes,
		Width:        int(math.Sqrt(float64(config.NumNodes))),
		AcceptPacket: true,

		numFlitPerStateDelaySamples: make(map[FlitState]int64),
		totalFlitPerStateDelays:     make(map[FlitState]int64),
		MaxFlitPerStateDelay:        make(map[FlitState]int),
	}

	for i := 0; i < network.NumNodes; i++ {
		var node = NewNode(network, i)
		network.Nodes = append(network.Nodes, node)
	}
	switch selection := config.Selection; selection {
	case SelectionAco:
		switch antPacketTraffic := config.AntPacketTraffic; antPacketTraffic {
		case TrafficUniform:
			network.AddTrafficGenerator(NewUniformTrafficGenerator(network, config.AntPacketInjectionRate, int64(-1), func(src int, dest int) Packet {
				return NewAntPacket(network, src, dest, config.AntPacketSize, func() {}, true)
			}))
		case TrafficTranspose1:
			network.AddTrafficGenerator(NewTranspose1TrafficGenerator(network, config.AntPacketInjectionRate, int64(-1), func(src int, dest int) Packet {
				return NewAntPacket(network, src, dest, config.AntPacketSize, func() {}, true)
			}))
		case TrafficTranspose2:
			network.AddTrafficGenerator(NewTranspose2TrafficGenerator(network, config.AntPacketInjectionRate, int64(-1), func(src int, dest int) Packet {
				return NewAntPacket(network, src, dest, config.AntPacketSize, func() {}, true)
			}))
		default:
			panic(fmt.Sprintf("ant packet traffic %s is not supported", antPacketTraffic))
		}
	}

	driver.CycleAccurateEventQueue().AddPerCycleEvent(func() {
		for _, node := range network.Nodes {
			node.Router.AdvanceOneCycle()
		}
	})

	return network
}

func (network *Network) GetX(id int) int {
	return id % network.Width
}

func (network *Network) GetY(id int) int {
	return (id - id%network.Width) / network.Width
}

func (network *Network) TrafficGenerators() []TrafficGenerator {
	return network.trafficGenerators
}

func (network *Network) AddTrafficGenerator(trafficGenerator TrafficGenerator) {
	network.trafficGenerators = append(network.trafficGenerators, trafficGenerator)

	network.Driver.CycleAccurateEventQueue().AddPerCycleEvent(func() {
		trafficGenerator.AdvanceOneCycle()
	})
}

func (network *Network) Receive(packet Packet) bool {
	if !network.Nodes[packet.Src()].Router.InjectPacket(packet) {
		network.Driver.CycleAccurateEventQueue().Schedule(func() {
			network.Receive(packet)
		}, 1)
		return false
	}

	network.LogPacketReceived(packet)

	return true
}

func (network *Network) LogPacketReceived(packet Packet) {
	network.NumPacketsReceived++

	if packet.HasPayload() {
		network.NumPayloadPacketsReceived++
	}
}

func (network *Network) LogPacketTransmitted(packet Packet) {
	network.NumPacketsTransmitted++

	if packet.HasPayload() {
		network.NumPayloadPacketsTransmitted++
	}

	network.totalPacketDelays += int64(Delay(packet))
	network.totalPacketHops += int64(Hops(packet))

	if packet.HasPayload() {
		network.totalPayloadPacketDelays += int64(Delay(packet))
		network.totalPayloadPacketHops += int64(Hops(packet))
	}

	network.MaxPacketDelay = int(math.Max(float64(network.MaxPacketDelay), float64(Delay(packet))))
	network.MaxPacketHops = int(math.Max(float64(network.MaxPacketHops), float64(Hops(packet))))

	if packet.HasPayload() {
		network.MaxPayloadPacketDelay = int(math.Max(float64(network.MaxPayloadPacketDelay), float64(Delay(packet))))
		network.MaxPayloadPacketHops = int(math.Max(float64(network.MaxPayloadPacketHops), float64(Hops(packet))))
	}
}

func (network *Network) LogFlitPerStateDelay(state FlitState, delay int) {
	if _, exists := network.numFlitPerStateDelaySamples[state]; !exists {
		network.numFlitPerStateDelaySamples[state] = int64(0)
	}

	network.numFlitPerStateDelaySamples[state]++

	if _, exists := network.totalFlitPerStateDelays[state]; !exists {
		network.totalFlitPerStateDelays[state] = int64(0)
	}

	network.totalFlitPerStateDelays[state] += int64(delay)

	if _, exists := network.MaxFlitPerStateDelay[state]; !exists {
		network.MaxFlitPerStateDelay[state] = 0
	}

	network.MaxFlitPerStateDelay[state] = int(math.Max(float64(network.MaxFlitPerStateDelay[state]), float64(delay)))
}

func (network *Network) Throughput() float64 {
	if network.Driver.CycleAccurateEventQueue().CurrentCycle == 0 {
		return float64(0)
	}

	return float64(network.NumPacketsTransmitted) / float64(network.Driver.CycleAccurateEventQueue().CurrentCycle) / float64(network.NumNodes)
}

func (network *Network) AveragePacketDelay() float64 {
	if network.NumPacketsTransmitted > 0 {
		return float64(network.totalPacketDelays) / float64(network.NumPacketsTransmitted)
	} else {
		return 0.0
	}
}

func (network *Network) AveragePacketHops() float64 {
	if network.NumPacketsTransmitted > 0 {
		return float64(network.totalPacketHops) / float64(network.NumPacketsTransmitted)
	} else {
		return 0.0
	}
}

func (network *Network) PayloadThroughput() float64 {
	if network.Driver.CycleAccurateEventQueue().CurrentCycle == 0 {
		return float64(0)
	}

	return float64(network.NumPayloadPacketsTransmitted) / float64(network.Driver.CycleAccurateEventQueue().CurrentCycle) / float64(network.NumNodes)
}

func (network *Network) AveragePayloadPacketDelay() float64 {
	if network.NumPayloadPacketsTransmitted > 0 {
		return float64(network.totalPayloadPacketDelays) / float64(network.NumPayloadPacketsTransmitted)
	} else {
		return 0.0
	}
}

func (network *Network) AveragePayloadPacketHops() float64 {
	if network.NumPayloadPacketsTransmitted > 0 {
		return float64(network.totalPayloadPacketHops) / float64(network.NumPayloadPacketsTransmitted)
	} else {
		return 0.0
	}
}

func (network *Network) AverageFlitPerStateDelay(state FlitState) float64 {
	if _, exists := network.numFlitPerStateDelaySamples[state]; exists {
		if network.numFlitPerStateDelaySamples[state] > 0 {
			return float64(network.totalFlitPerStateDelays[state]) / float64(network.numFlitPerStateDelaySamples[state])
		}
	}

	return 0.0
}
