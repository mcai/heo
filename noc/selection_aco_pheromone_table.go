package noc

type PheromoneTable struct {
	Node       *Node
	Pheromones map[int]map[Direction]*Pheromone
}

func NewPheromoneTable(node *Node) *PheromoneTable {
	var pheromoneTable = &PheromoneTable{
		Node:       node,
		Pheromones: make(map[int]map[Direction]*Pheromone),
	}

	return pheromoneTable
}

func (pheromoneTable *PheromoneTable) Append(dest int, direction Direction, pheromoneValue float64) {
	var pheromone = NewPheromone(pheromoneTable, dest, direction, pheromoneValue)

	if _, exists := pheromoneTable.Pheromones[dest]; !exists {
		pheromoneTable.Pheromones[dest] = make(map[Direction]*Pheromone)
	}

	pheromoneTable.Pheromones[dest][direction] = pheromone
}

func (pheromoneTable *PheromoneTable) Update(dest int, direction Direction) {
	for _, pheromone := range pheromoneTable.Pheromones[dest] {
		if pheromone.Direction == direction {
			pheromone.Value += pheromoneTable.Node.Network.Config().ReinforcementFactor * (1 - pheromone.Value)
		} else {
			pheromone.Value -= pheromoneTable.Node.Network.Config().ReinforcementFactor * pheromone.Value
		}
	}
}
