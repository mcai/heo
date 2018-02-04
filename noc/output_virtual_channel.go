package noc

type OutputVirtualChannel struct {
	OutputPort          *OutputPort
	Num                 int
	InputVirtualChannel *InputVirtualChannel
	Credits             int
	Arbiter             *VirtualChannelArbiter
}

func NewOutputVirtualChannel(outputPort *OutputPort, num int) *OutputVirtualChannel {
	var outputVirtualChannel = &OutputVirtualChannel{
		OutputPort: outputPort,
		Num:        num,
		Credits:    10,
	}

	outputVirtualChannel.Arbiter = NewVirtualChannelArbiter(outputVirtualChannel)

	return outputVirtualChannel
}
