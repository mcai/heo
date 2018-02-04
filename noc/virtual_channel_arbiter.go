package noc

type VirtualChannelArbiter struct {
	OutputVirtualChannel    *OutputVirtualChannel
	InputVirtualChannelRing *InputVirtualChannelRing
}

func NewVirtualChannelArbiter(outputVirtualChannel *OutputVirtualChannel) *VirtualChannelArbiter {
	var arbiter = &VirtualChannelArbiter{
		OutputVirtualChannel:    outputVirtualChannel,
		InputVirtualChannelRing: NewInputVirtualChannelRing(outputVirtualChannel.OutputPort.Router),
	}

	return arbiter
}

func (arbiter *VirtualChannelArbiter) Next() *InputVirtualChannel {
	for i := 0; i < len(arbiter.InputVirtualChannelRing.GetChannels()); i++ {
		var inputVirtualChannel = arbiter.InputVirtualChannelRing.Next()
		if inputVirtualChannel.Route == arbiter.OutputVirtualChannel.OutputPort.Direction {
			var flit = inputVirtualChannel.InputBuffer.Peek()
			if flit != nil && flit.Head && flit.GetState() == FLIT_STATE_ROUTE_COMPUTATION {
				return inputVirtualChannel
			}
		}
	}

	return nil
}
