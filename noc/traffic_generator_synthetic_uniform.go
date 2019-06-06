package noc

import "math/rand"

type UniformTrafficGenerator struct {
	*BaseSyntheticTrafficGenerator
}

func NewUniformTrafficGenerator(network *Network, packetInjectionRate float64, maxPackets int64, newPacket func(src int, dest int) Packet) *UniformTrafficGenerator {
	var generator = &UniformTrafficGenerator{
		BaseSyntheticTrafficGenerator: NewBaseSyntheticTrafficGenerator(network, packetInjectionRate, maxPackets, newPacket),
	}

	return generator
}

func (generator *UniformTrafficGenerator) AdvanceOneCycle() {
	generator.BaseSyntheticTrafficGenerator.AdvanceOneCycle(func(src int) int {
		for {
			var i = rand.Intn(generator.Network.NumNodes)
			if i != src {
				return i
			}
		}
	})
}
