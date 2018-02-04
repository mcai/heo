package noc

type NorthLastRoutingAlgorithm struct {
	Node               *Node
	XYRoutingAlgorithm *XYRoutingAlgorithm
}

func NewNorthLastRoutingAlgorithm(node *Node) *NorthLastRoutingAlgorithm {
	var routingAlgorithm = &NorthLastRoutingAlgorithm{
		Node:               node,
		XYRoutingAlgorithm: NewXYRoutingAlgorithm(node),
	}

	return routingAlgorithm
}

func (routingAlgorithm *NorthLastRoutingAlgorithm) NextHop(packet Packet, parent int) []Direction {
	var directions []Direction

	var destX = routingAlgorithm.Node.Network.GetX(packet.Dest())
	var destY = routingAlgorithm.Node.Network.GetY(packet.Dest())

	var x = routingAlgorithm.Node.X
	var y = routingAlgorithm.Node.Y

	if destX == x || destY <= y {
		return routingAlgorithm.XYRoutingAlgorithm.NextHop(packet, parent)
	}

	if destX < x {
		directions = append(directions, DIRECTION_SOUTH)
		directions = append(directions, DIRECTION_WEST)
	} else {
		directions = append(directions, DIRECTION_SOUTH)
		directions = append(directions, DIRECTION_EAST)
	}

	return directions
}
