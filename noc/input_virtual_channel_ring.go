package noc

import "container/ring"

type InputVirtualChannelRing struct {
	Router   *Router
	Channels []*InputVirtualChannel
	ring     *ring.Ring
}

func NewInputVirtualChannelRing(router *Router) *InputVirtualChannelRing {
	var inputVirtualChannelRing = &InputVirtualChannelRing{
		Router:router,
	}

	return inputVirtualChannelRing
}

func (inputVirtualChannelRing *InputVirtualChannelRing) GetChannels() []*InputVirtualChannel {
	if inputVirtualChannelRing.Channels == nil {
		inputVirtualChannelRing.Channels = inputVirtualChannelRing.Router.GetInputVirtualChannels()
	}

	return inputVirtualChannelRing.Channels
}

func (inputVirtualChannelRing *InputVirtualChannelRing) GetRing() *ring.Ring {
	if inputVirtualChannelRing.ring == nil {
		var inputVirtualChannels = inputVirtualChannelRing.GetChannels()

		var r = ring.New(len(inputVirtualChannels))

		for _, inputVirtualChannel := range inputVirtualChannels {
			r.Value = inputVirtualChannel
			r = r.Next()
		}

		inputVirtualChannelRing.ring = r
	}

	return inputVirtualChannelRing.ring
}

func (inputVirtualChannelRing *InputVirtualChannelRing) Next() *InputVirtualChannel {
	inputVirtualChannelRing.ring = inputVirtualChannelRing.GetRing().Next()
	return inputVirtualChannelRing.GetRing().Value.(*InputVirtualChannel)
}