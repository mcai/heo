package noc

type ACOSelectionAlgorithm struct {
	Node           *Node
	PheromoneTable *PheromoneTable
}

func NewACOSelectionAlgorithm(node *Node) *ACOSelectionAlgorithm {
	var selectionAlgorithm = &ACOSelectionAlgorithm{
		Node:           node,
		PheromoneTable: NewPheromoneTable(node),
	}

	var pheromoneValue = 1.0 / float64(len(node.Neighbors))

	for dest := 0; dest < node.Network.NumNodes; dest++ {
		if node.Id != dest {
			for i := range node.Neighbors {
				var direction = Direction(i)
				selectionAlgorithm.PheromoneTable.Append(dest, direction, pheromoneValue)
			}
		}
	}

	return selectionAlgorithm
}

func (selectionAlgorithm *ACOSelectionAlgorithm) CreateAndSendBackwardAntPacket(packet *AntPacket) {
	var newPacket = NewAntPacket(packet.Network(), packet.Dest(), packet.Src(), selectionAlgorithm.Node.Network.Config().AntPacketSize, func() {}, false)

	newPacket.memory = packet.memory

	selectionAlgorithm.Node.Network.driver.CycleAccurateEventQueue().Schedule(func() {
		selectionAlgorithm.Node.Network.Receive(newPacket)
	}, 1)
}

func (selectionAlgorithm *ACOSelectionAlgorithm) BackwardAntPacket(packet *AntPacket) Direction {
	var i int

	for i = len(packet.memory) - 1; i > 0; i-- {
		var entry = packet.memory[i]
		if entry.NodeId == selectionAlgorithm.Node.Id {
			break
		}
	}

	var prev = packet.memory[i-1].NodeId

	for direction, neighbor := range selectionAlgorithm.Node.Neighbors {
		if neighbor == prev {
			return direction
		}
	}

	panic("Impossible")
}

func (selectionAlgorithm *ACOSelectionAlgorithm) UpdatePheromoneTable(packet *AntPacket, inputVirtualChannel *InputVirtualChannel) {
	var i int
	for i = 0; i < len(packet.memory); i++ {
		var entry = packet.memory[i]
		if entry.NodeId == selectionAlgorithm.Node.Id {
			break
		}
	}

	for j := i + 1; j < len(packet.memory); j++ {
		var entry = packet.memory[j]
		var dest = entry.NodeId
		selectionAlgorithm.PheromoneTable.Update(dest, inputVirtualChannel.InputPort.Direction)
	}
}

func (selectionAlgorithm *ACOSelectionAlgorithm) Select(packet Packet, ivc int, directions []Direction) Direction {
	var maxProbability = -1.0
	var bestDirection = DIRECTION_UNKNOWN

	for _, direction := range directions {
		var neighbor = selectionAlgorithm.Node.Neighbors[direction]
		var neighborRouter = selectionAlgorithm.Node.Network.Nodes[neighbor].Router
		var pheromone = selectionAlgorithm.PheromoneTable.Pheromones[packet.Dest()][direction]
		var freeSlots = neighborRouter.FreeSlots(direction.GetReflexDirection(), ivc)

		var acoSelectionAlpha = selectionAlgorithm.Node.Network.Config().AcoSelectionAlpha
		var maxInputBufferSize = selectionAlgorithm.Node.Network.Config().MaxInputBufferSize
		var numNeighbors = len(selectionAlgorithm.Node.Neighbors)

		var probability = (pheromone.Value + acoSelectionAlpha*(float64(freeSlots)/float64(maxInputBufferSize))) /
			(1 + acoSelectionAlpha*float64(numNeighbors-1))

		if probability > maxProbability {
			maxProbability = probability
			bestDirection = direction
		}
	}

	return bestDirection
}
