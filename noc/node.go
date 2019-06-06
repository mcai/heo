package noc

import (
	"fmt"
)

type Node struct {
	Network            *BaseNetwork
	Id                 int
	X, Y               int
	Neighbors          map[Direction]int
	Router             *Router
	RoutingAlgorithm   RoutingAlgorithm
	SelectionAlgorithm SelectionAlgorithm
}

func NewNode(network *BaseNetwork, id int) *Node {
	var node = &Node{
		Network:   network,
		Id:        id,
		X:         network.GetX(id),
		Y:         network.GetY(id),
		Neighbors: make(map[Direction]int),
	}

	if id/network.Width > 0 {
		node.Neighbors[DIRECTION_NORTH] = id - network.Width
	}

	if (id % network.Width) != network.Width-1 {
		node.Neighbors[DIRECTION_EAST] = id + 1
	}

	if id/network.Width < network.Width-1 {
		node.Neighbors[DIRECTION_SOUTH] = id + network.Width
	}

	if id%network.Width != 0 {
		node.Neighbors[DIRECTION_WEST] = id - 1
	}

	node.Router = NewRouter(node)

	switch routing := network.Config().Routing; routing {
	case RoutingXY:
		node.RoutingAlgorithm = NewXYRoutingAlgorithm(node)
	case RoutingNegativeFirst:
		node.RoutingAlgorithm = NewNegativeFirstRoutingAlgorithm(node)
	case RoutingWestFirst:
		node.RoutingAlgorithm = NewWestFirstRoutingAlgorithm(node)
	case RoutingNorthLast:
		node.RoutingAlgorithm = NewNorthLastRoutingAlgorithm(node)
	case RoutingOddEven:
		node.RoutingAlgorithm = NewOddEvenRoutingAlgorithm(node)
	default:
		panic(fmt.Sprintf("Not supported: %s", routing))
	}

	switch selection := network.Config().Selection; selection {
	case SelectionRandom:
		node.SelectionAlgorithm = NewRandomSelectionAlgorithm(node)
	case SelectionBufferLevel:
		node.SelectionAlgorithm = NewBufferLevelSelectionAlgorithm(node)
	case SelectionAco:
		node.SelectionAlgorithm = NewACOSelectionAlgorithm(node)
	default:
		panic(fmt.Sprintf("Not supported: %s", selection))
	}

	return node
}

func (node *Node) DumpNeighbors() {
	for direction, neighbor := range node.Neighbors {
		fmt.Printf("node#%d.neighbors[%s]=%d\n", node.Id, direction, neighbor)
	}

	fmt.Println()
}
