package noc

import "math/rand"

type UniformTrafficSource struct {
	*BaseSyntheticTrafficSource
}

func NewUniformTrafficSource(network *Network, packetInjectionRate float64, maxPackets int64, newPacket func(src int, dest int) Packet) *UniformTrafficSource {
	var source = &UniformTrafficSource{
		BaseSyntheticTrafficSource: NewBaseSyntheticTrafficSource(network, packetInjectionRate, maxPackets, newPacket),
	}

	return source
}

func (source *UniformTrafficSource) AdvanceOneCycle() {
	source.BaseSyntheticTrafficSource.AdvanceOneCycle(func(src int) int {
		for {
			var i = rand.Intn(source.Network.NumNodes)
			if i != src {
				return i
			}
		}
	})
}
