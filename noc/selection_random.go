package noc

import "math/rand"

type RandomSelectionAlgorithm struct {
	Node *Node
}

func NewRandomSelectionAlgorithm(node *Node) *RandomSelectionAlgorithm {
	var selectionAlgorithm = &RandomSelectionAlgorithm{
		Node:node,
	}

	return selectionAlgorithm
}

func (selectionAlgorithm *RandomSelectionAlgorithm) Select(packet Packet, ivc int, directions []Direction) Direction {
	return directions[rand.Intn(len(directions))]
}
