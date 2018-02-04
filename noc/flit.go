package noc

import "fmt"

type Flit struct {
	Packet             Packet
	Num                int
	Head               bool
	Tail               bool
	node               *Node
	state              FlitState
	prevStateTimestamp int64
	Timestamp          int64
}

func NewFlit(packet Packet, num int, head bool, tail bool) *Flit {
	var flit = &Flit{
		Packet:             packet,
		Num:                num,
		Head:               head,
		Tail:               tail,
		state:              FLIT_STATE_UNKNOWN,
		prevStateTimestamp: packet.Network().Driver.CycleAccurateEventQueue().CurrentCycle,
		Timestamp:          packet.Network().Driver.CycleAccurateEventQueue().CurrentCycle,
	}

	packet.SetFlits(append(packet.Flits(), flit))

	return flit
}

func (flit *Flit) GetNode() *Node {
	return flit.node
}

func (flit *Flit) GetState() FlitState {
	return flit.state
}

func (flit *Flit) SetNodeAndState(node *Node, state FlitState) {
	if state == flit.state {
		panic(fmt.Sprintf("Flit is already in the %s state", state))
	}

	if flit.state != FLIT_STATE_UNKNOWN {
		flit.Packet.Network().LogFlitPerStateDelay(flit.state, int(flit.Packet.Network().Driver.CycleAccurateEventQueue().CurrentCycle-flit.prevStateTimestamp))

		if flit.GetNumInflightFlits()[flit.state] == 0 {
			panic("Impossible")
		}

		flit.GetNumInflightFlits()[flit.state] = flit.GetNumInflightFlits()[flit.state] - 1
	}

	flit.node = node
	flit.state = state

	if flit.state != FLIT_STATE_DESTINATION_ARRIVED {
		flit.GetNumInflightFlits()[flit.state] = flit.GetNumInflightFlits()[flit.state] + 1
	}

	flit.prevStateTimestamp = flit.Packet.Network().Driver.CycleAccurateEventQueue().CurrentCycle
}

func (flit *Flit) GetNumInflightFlits() map[FlitState]int {
	if flit.Head {
		return flit.node.Router.NumInflightHeadFlits
	} else {
		return flit.node.Router.NumInflightNonHeadFlits
	}
}
