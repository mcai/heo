package noc

type WestFirstRoutingAlgorithm struct {
	Node               *Node
	XYRoutingAlgorithm *XYRoutingAlgorithm
}

func NewWestFirstRoutingAlgorithm(node *Node) *WestFirstRoutingAlgorithm {
	var routingAlgorithm = &WestFirstRoutingAlgorithm{
		Node:               node,
		XYRoutingAlgorithm: NewXYRoutingAlgorithm(node),
	}

	return routingAlgorithm
}

func (routingAlgorithm *WestFirstRoutingAlgorithm) NextHop(packet Packet, parent int) []Direction {
	var directions []Direction

	var destX = routingAlgorithm.Node.Network.GetX(packet.Dest())
	var destY = routingAlgorithm.Node.Network.GetY(packet.Dest())

	var x = routingAlgorithm.Node.X
	var y = routingAlgorithm.Node.Y

	if destX <= x || destY == y {
		return routingAlgorithm.XYRoutingAlgorithm.NextHop(packet, parent)
	}

	if destY < y {
		directions = append(directions, DIRECTION_NORTH)
		directions = append(directions, DIRECTION_EAST)
	} else {
		directions = append(directions, DIRECTION_SOUTH)
		directions = append(directions, DIRECTION_EAST)
	}

	return directions
}
