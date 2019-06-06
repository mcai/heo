package noc

type Transpose1TrafficSource struct {
	*BaseSyntheticTrafficSource
}

func NewTranspose1TrafficSource(network *Network, packetInjectionRate float64, maxPackets int64, newPacket func(src int, dest int) Packet) *Transpose1TrafficSource {
	var source = &Transpose1TrafficSource{
		BaseSyntheticTrafficSource: NewBaseSyntheticTrafficSource(network, packetInjectionRate, maxPackets, newPacket),
	}

	return source
}

func (source *Transpose1TrafficSource) AdvanceOneCycle() {
	source.BaseSyntheticTrafficSource.AdvanceOneCycle(func(src int) int {
		var srcX, srcY = source.Network.GetX(src), source.Network.GetY(src)
		var destX, destY = source.Network.Width-1-srcY, source.Network.Width-1-srcX

		return destY*source.Network.Width + destX
	})
}
