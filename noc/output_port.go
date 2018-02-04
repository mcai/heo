package noc

type OutputPort struct {
	Router          *Router
	Direction       Direction
	VirtualChannels []*OutputVirtualChannel
	Arbiter         *SwitchArbiter
}

func NewOutputPort(router *Router, direction Direction) *OutputPort {
	var outputPort = &OutputPort{
		Router:    router,
		Direction: direction,
	}

	for i := 0; i < router.Node.Network.Config.NumVirtualChannels; i++ {
		var outputVirtualChannel = NewOutputVirtualChannel(outputPort, i)
		outputPort.VirtualChannels = append(outputPort.VirtualChannels, outputVirtualChannel)
	}

	outputPort.Arbiter = NewSwitchArbiter(outputPort)

	return outputPort
}
