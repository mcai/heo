package noc

type InputVirtualChannelRing struct {
	Router   *Router
	Channels []*InputVirtualChannel
	currentChannelIndex int
}

func NewInputVirtualChannelRing(router *Router) *InputVirtualChannelRing {
	var inputVirtualChannelRing = &InputVirtualChannelRing{
		Router: router,
		currentChannelIndex: 0,
	}

	return inputVirtualChannelRing
}

func (inputVirtualChannelRing *InputVirtualChannelRing) GetChannels() []*InputVirtualChannel {
	if inputVirtualChannelRing.Channels == nil {
		inputVirtualChannelRing.Channels = inputVirtualChannelRing.Router.GetInputVirtualChannels()
	}

	return inputVirtualChannelRing.Channels
}

func (inputVirtualChannelRing *InputVirtualChannelRing) Next() *InputVirtualChannel {
	var next = inputVirtualChannelRing.Channels[inputVirtualChannelRing.currentChannelIndex]

	inputVirtualChannelRing.currentChannelIndex =
		(inputVirtualChannelRing.currentChannelIndex + 1) % len(inputVirtualChannelRing.Channels)

	return next
}
