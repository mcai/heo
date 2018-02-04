package noc

type OddEvenRoutingAlgorithm struct {
	Node *Node
}

func NewOddEvenRoutingAlgorithm(node *Node) *OddEvenRoutingAlgorithm {
	var routingAlgorithm = &OddEvenRoutingAlgorithm{
		Node: node,
	}

	return routingAlgorithm
}

func (routingAlgorithm *OddEvenRoutingAlgorithm) NextHop(packet Packet, parent int) []Direction {
	var directions []Direction

	var c0 = routingAlgorithm.Node.X
	var c1 = routingAlgorithm.Node.Y

	var s0 = routingAlgorithm.Node.Network.GetX(packet.Src())

	var d0 = routingAlgorithm.Node.Network.GetX(packet.Dest())
	var d1 = routingAlgorithm.Node.Network.GetY(packet.Dest())

	var e0 = d0 - c0
	var e1 = -(d1 - c1)

	if e0 == 0 {
		if e1 > 0 {
			directions = append(directions, DIRECTION_NORTH)
		} else {
			directions = append(directions, DIRECTION_SOUTH)
		}
	} else {
		if e0 > 0 {
			if e1 == 0 {
				directions = append(directions, DIRECTION_EAST)
			} else {
				if c0%2 == 1 || c0 == s0 {
					if e1 > 0 {
						directions = append(directions, DIRECTION_NORTH)
					} else {
						directions = append(directions, DIRECTION_SOUTH)
					}
				}

				if d0%2 == 1 || e0 != 1 {
					directions = append(directions, DIRECTION_EAST)
				}
			}
		} else {
			directions = append(directions, DIRECTION_WEST)
			if c0%2 == 0 {
				if e1 > 0 {
					directions = append(directions, DIRECTION_NORTH)
				}
				if e1 < 0 {
					directions = append(directions, DIRECTION_SOUTH)
				}
			}
		}
	}

	return directions
}
