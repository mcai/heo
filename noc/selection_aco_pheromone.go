package noc

type Pheromone struct {
	PheromoneTable *PheromoneTable
	Dest           int
	Direction      Direction
	Value          float64
}

func NewPheromone(pheromoneTable *PheromoneTable, dest int, direction Direction, value float64) *Pheromone {
	var pheromone = &Pheromone{
		PheromoneTable:pheromoneTable,
		Dest:dest,
		Direction:direction,
		Value:value,
	}

	return pheromone
}