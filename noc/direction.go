package noc

import "fmt"

type Direction string

const (
	DIRECTION_UNKNOWN = Direction("UNKNOWN")

	DIRECTION_LOCAL = Direction("LOCAL")

	DIRECTION_NORTH = Direction("NORTH")

	DIRECTION_EAST = Direction("EAST")

	DIRECTION_SOUTH = Direction("SOUTH")

	DIRECTION_WEST = Direction("WEST")
)

func (direction Direction) GetReflexDirection() Direction {
	switch direction {
	case DIRECTION_LOCAL:
		return DIRECTION_LOCAL
	case DIRECTION_NORTH:
		return DIRECTION_SOUTH
	case DIRECTION_EAST:
		return DIRECTION_WEST
	case DIRECTION_SOUTH:
		return DIRECTION_NORTH
	case DIRECTION_WEST:
		return DIRECTION_EAST
	default:
		panic(fmt.Sprintf("Cannot get the reflex direction of %d", direction))
	}
}
