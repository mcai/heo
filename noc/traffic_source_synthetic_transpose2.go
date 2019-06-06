package noc

type Transpose2TrafficSource struct {
	*BaseSyntheticTrafficSource
}

func NewTranspose2TrafficSource(network *Network, packetInjectionRate float64, maxPackets int64, newPacket func(src int, dest int) Packet) *Transpose2TrafficSource {
	var source = &Transpose2TrafficSource{
		BaseSyntheticTrafficSource: NewBaseSyntheticTrafficSource(network, packetInjectionRate, maxPackets, newPacket),
	}

	return source
}

func (source *Transpose2TrafficSource) AdvanceOneCycle() {
	source.BaseSyntheticTrafficSource.AdvanceOneCycle(func(src int) int {
		var srcX, srcY = source.Network.GetX(src), source.Network.GetY(src)
		var destX, destY = srcY, srcX

		return destY*source.Network.Width + destX
	})
}
