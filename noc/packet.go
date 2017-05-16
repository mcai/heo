package noc

type Packet interface {
	Network() *Network
	Id() int64
	BeginCycle() int64
	EndCycle() int64
	SetEndCycle(endCycle int64)
	Src() int
	Dest() int
	Size() int
	OnCompletedCallback() func()
	Memory() []*PacketMemoryEntry
	Flits() []*Flit
	SetFlits(flits []*Flit)
	HasPayload() bool
	HandleDestArrived(inputVirtualChannel *InputVirtualChannel)
	DoRouteComputation(inputVirtualChannel *InputVirtualChannel) Direction
}

func Delay(packet Packet) int {
	if packet.EndCycle() == -1 {
		return -1
	} else {
		return int(packet.EndCycle() - packet.BeginCycle())
	}
}

func Hops(packet Packet) int {
	return len(packet.Memory())
}

type PacketMemoryEntry struct {
	NodeId    int
	Timestamp int64
}
