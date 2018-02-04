package noc

type SwitchArbiter struct {
	OutputPort              *OutputPort
	InputVirtualChannelRing *InputVirtualChannelRing
}

func NewSwitchArbiter(outputPort *OutputPort) *SwitchArbiter {
	var arbiter = &SwitchArbiter{
		OutputPort:              outputPort,
		InputVirtualChannelRing: NewInputVirtualChannelRing(outputPort.Router),
	}

	return arbiter
}

func (arbiter *SwitchArbiter) Next() *InputVirtualChannel {
	for i := 0; i < len(arbiter.InputVirtualChannelRing.GetChannels()); i++ {
		var inputVirtualChannel = arbiter.InputVirtualChannelRing.Next()
		if inputVirtualChannel.OutputVirtualChannel != nil && inputVirtualChannel.OutputVirtualChannel.OutputPort == arbiter.OutputPort {
			var flit = inputVirtualChannel.InputBuffer.Peek()
			if flit != nil && (flit.Head && flit.GetState() == FLIT_STATE_VIRTUAL_CHANNEL_ALLOCATION || !flit.Head && flit.GetState() == FLIT_STATE_INPUT_BUFFER) {
				return inputVirtualChannel
			}
		}
	}

	return nil
}
