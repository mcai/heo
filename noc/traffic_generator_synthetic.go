package noc

import "math/rand"

type BaseSyntheticTrafficGenerator struct {
	Network             *Network
	PacketInjectionRate float64
	MaxPackets          int64
	NewPacket           func(src int, dest int) Packet
}

func NewBaseSyntheticTrafficGenerator(network *Network, packetInjectionRate float64, maxPackets int64, newPacket func(src int, dest int) Packet) *BaseSyntheticTrafficGenerator {
	var baseSyntheticTrafficGenerator = &BaseSyntheticTrafficGenerator{
		Network:             network,
		PacketInjectionRate: packetInjectionRate,
		MaxPackets:          maxPackets,
		NewPacket:           newPacket,
	}

	return baseSyntheticTrafficGenerator
}

func (generator *BaseSyntheticTrafficGenerator) AdvanceOneCycle(dest func(src int) int) {
	for _, node := range generator.Network.Nodes {
		if !generator.Network.AcceptPacket || generator.MaxPackets != -1 && generator.Network.NumPacketsReceived > generator.MaxPackets {
			break
		}

		if rand.Float64() <= generator.PacketInjectionRate {
			var src = node.Id
			var dest = dest(src)

			if src != dest {
				generator.Network.Driver.CycleAccurateEventQueue().Schedule(func() {
					generator.Network.Receive(generator.NewPacket(src, dest))
				}, 1)
			}
		}
	}
}
