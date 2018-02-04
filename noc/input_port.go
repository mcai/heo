package noc

type InputPort struct {
	Router          *Router
	Direction       Direction
	VirtualChannels []*InputVirtualChannel
}

func NewInputPort(router *Router, direction Direction) *InputPort {
	var inputPort = &InputPort{
		Router:    router,
		Direction: direction,
	}

	for i := 0; i < router.Node.Network.Config.NumVirtualChannels; i++ {
		inputPort.VirtualChannels = append(inputPort.VirtualChannels, NewInputVirtualChannel(inputPort, i))
	}

	return inputPort
}
