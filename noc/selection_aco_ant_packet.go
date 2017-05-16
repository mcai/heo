package noc

type AntPacket struct {
	*DataPacket
	Forward bool
}

func NewAntPacket(network *Network, src int, dest int, size int, onCompletedCallback func(), forward bool) *AntPacket {
	var packet = &AntPacket{
		DataPacket:NewDataPacket(network, src, dest, size, false, onCompletedCallback),
		Forward:forward,
	}

	return packet
}

func (packet *AntPacket) HandleDestArrived(inputVirtualChannel *InputVirtualChannel) {
	var selectionAlgorithm = inputVirtualChannel.InputPort.Router.Node.SelectionAlgorithm.(*ACOSelectionAlgorithm)

	if packet.Forward {
		packet.Memorize(inputVirtualChannel.InputPort.Router.Node)
		selectionAlgorithm.CreateAndSendBackwardAntPacket(packet)
	} else {
		selectionAlgorithm.UpdatePheromoneTable(packet, inputVirtualChannel)
	}

	packet.endCycle = inputVirtualChannel.InputPort.Router.Node.Network.Driver.CycleAccurateEventQueue().CurrentCycle

	inputVirtualChannel.InputPort.Router.Node.Network.LogPacketTransmitted(packet)

	if packet.OnCompletedCallback != nil {
		packet.OnCompletedCallback()
	}
}

func (packet *AntPacket) DoRouteComputation(inputVirtualChannel *InputVirtualChannel) Direction {
	var selectionAlgorithm = inputVirtualChannel.InputPort.Router.Node.SelectionAlgorithm.(*ACOSelectionAlgorithm)

	if packet.Forward {
		return packet.DataPacket.DoRouteComputation(inputVirtualChannel)
	} else {
		if inputVirtualChannel.InputPort.Router.Node.Id != packet.src {
			selectionAlgorithm.UpdatePheromoneTable(packet, inputVirtualChannel)
		}

		return selectionAlgorithm.BackwardAntPacket(packet)
	}
}
