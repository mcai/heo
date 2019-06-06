package noc

import (
	"fmt"
	"math"
)

type DataPacket struct {
	network              Network
	id                   int64
	beginCycle, endCycle int64
	src, dest            int
	size                 int
	onCompletedCallback  func()
	memory               []*PacketMemoryEntry
	flits                []*Flit
	hasPayload           bool
}

func NewDataPacket(network Network, src int, dest int, size int, hasPayload bool, onCompletedCallback func()) *DataPacket {
	var packet = &DataPacket{
		network:             network,
		id:                  network.CurrentPacketId(),
		beginCycle:          network.Driver().CycleAccurateEventQueue().CurrentCycle,
		endCycle:            -1,
		src:                 src,
		dest:                dest,
		size:                size,
		onCompletedCallback: onCompletedCallback,
		hasPayload:          hasPayload,
	}

	network.SetCurrentPacketId(network.CurrentPacketId() + 1)

	var numFlits = int(math.Ceil(float64(packet.size) / float64(network.Config().LinkWidth)))
	if numFlits > network.Config().MaxInputBufferSize {
		panic(fmt.Sprintf("Number of flits (%d) in a packet cannot be greater than max input buffer size (%d)", numFlits, network.Config().MaxInputBufferSize))
	}

	return packet
}

func (packet *DataPacket) Network() Network {
	return packet.network
}

func (packet *DataPacket) Id() int64 {
	return packet.id
}

func (packet *DataPacket) BeginCycle() int64 {
	return packet.beginCycle
}

func (packet *DataPacket) EndCycle() int64 {
	return packet.endCycle
}

func (packet *DataPacket) SetEndCycle(endCycle int64) {
	packet.endCycle = endCycle
}

func (packet *DataPacket) Src() int {
	return packet.src
}

func (packet *DataPacket) Dest() int {
	return packet.dest
}

func (packet *DataPacket) Size() int {
	return packet.size
}

func (packet *DataPacket) OnCompletedCallback() func() {
	return packet.onCompletedCallback
}

func (packet *DataPacket) Memory() []*PacketMemoryEntry {
	return packet.memory
}

func (packet *DataPacket) Flits() []*Flit {
	return packet.flits
}

func (packet *DataPacket) SetFlits(flits []*Flit) {
	packet.flits = flits
}

func (packet *DataPacket) HasPayload() bool {
	return packet.hasPayload
}

func (packet *DataPacket) HandleDestArrived(inputVirtualChannel *InputVirtualChannel) {
	packet.Memorize(inputVirtualChannel.InputPort.Router.Node)

	packet.endCycle = inputVirtualChannel.InputPort.Router.Node.Network.driver.CycleAccurateEventQueue().CurrentCycle

	inputVirtualChannel.InputPort.Router.Node.Network.LogPacketTransmitted(packet)

	if packet.onCompletedCallback != nil {
		packet.onCompletedCallback()
	}
}

func (packet *DataPacket) DoRouteComputation(inputVirtualChannel *InputVirtualChannel) Direction {
	var parent = -1

	if len(packet.memory) > 0 {
		parent = packet.memory[len(packet.memory)-1].NodeId
	}

	packet.Memorize(inputVirtualChannel.InputPort.Router.Node)

	var directions = inputVirtualChannel.InputPort.Router.Node.RoutingAlgorithm.NextHop(packet, parent)

	return inputVirtualChannel.InputPort.Router.Node.SelectionAlgorithm.Select(packet, inputVirtualChannel.Num, directions)
}

func (packet *DataPacket) Memorize(node *Node) {
	for _, entry := range packet.memory {
		if entry.NodeId == node.Id {
			panic(fmt.Sprintf("packet#%d(src=%d, dest=%d): %d", packet.id, packet.src, packet.dest, node.Id))
		}
	}

	packet.memory = append(packet.memory, &PacketMemoryEntry{
		NodeId:    node.Id,
		Timestamp: packet.network.Driver().CycleAccurateEventQueue().CurrentCycle,
	})
}

func (packet *DataPacket) DumpMemory() {
	for i, entry := range packet.memory {
		fmt.Printf("packet#%d.memory[%d]=%d\n", packet.id, i, entry.NodeId)
	}

	fmt.Println()
}
