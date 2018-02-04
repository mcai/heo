package noc

type NegativeFirstRoutingAlgorithm struct {
	Node               *Node
	XYRoutingAlgorithm *XYRoutingAlgorithm
}

func NewNegativeFirstRoutingAlgorithm(node *Node) *NegativeFirstRoutingAlgorithm {
	var routingAlgorithm = &NegativeFirstRoutingAlgorithm{
		Node:               node,
		XYRoutingAlgorithm: NewXYRoutingAlgorithm(node),
	}

	return routingAlgorithm
}

func (routingAlgorithm *NegativeFirstRoutingAlgorithm) NextHop(packet Packet, parent int) []Direction {
	var directions []Direction

	var destX = routingAlgorithm.Node.Network.GetX(packet.Dest())
	var destY = routingAlgorithm.Node.Network.GetY(packet.Dest())

	var x = routingAlgorithm.Node.X
	var y = routingAlgorithm.Node.Y

	if (destX <= x && destY <= y) || (destX >= x && destY >= y) {
		return routingAlgorithm.XYRoutingAlgorithm.NextHop(packet, parent)
	}

	if destX > x && destY < y {
		directions = append(directions, DIRECTION_NORTH)
		directions = append(directions, DIRECTION_EAST)
	} else {
		directions = append(directions, DIRECTION_SOUTH)
		directions = append(directions, DIRECTION_WEST)
	}

	return directions
}
