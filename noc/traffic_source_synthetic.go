package noc

import "math/rand"

type BaseSyntheticTrafficSource struct {
	Network             *BaseNetwork
	PacketInjectionRate float64
	MaxPackets          int64
	NewPacket           func(src int, dest int) Packet
}

func NewBaseSyntheticTrafficSource(network *BaseNetwork, packetInjectionRate float64, maxPackets int64, newPacket func(src int, dest int) Packet) *BaseSyntheticTrafficSource {
	var baseSyntheticTrafficSource = &BaseSyntheticTrafficSource{
		Network:             network,
		PacketInjectionRate: packetInjectionRate,
		MaxPackets:          maxPackets,
		NewPacket:           newPacket,
	}

	return baseSyntheticTrafficSource
}

func (source *BaseSyntheticTrafficSource) AdvanceOneCycle(dest func(src int) int) {
	for _, node := range source.Network.Nodes {
		if !source.Network.AcceptPacket || source.MaxPackets != -1 && source.Network.NumPacketsReceived > source.MaxPackets {
			break
		}

		if rand.Float64() <= source.PacketInjectionRate {
			var src = node.Id
			var dest = dest(src)

			if src != dest {
				source.Network.Driver().CycleAccurateEventQueue().Schedule(func() {
					source.Network.Receive(source.NewPacket(src, dest))
				}, 1)
			}
		}
	}
}
