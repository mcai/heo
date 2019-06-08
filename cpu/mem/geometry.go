package mem

import "math"

type Geometry struct {
	Size           uint32
	Assoc          uint32
	LineSize       uint32
	LineSizeInLog2 uint32
	NumSets        uint32
	NumSetsInLog2  uint32
	NumLines       uint32
}

func NewGeometry(size uint32, assoc uint32, lineSize uint32) *Geometry {
	var geometry = &Geometry{
		Size:           size,
		Assoc:          assoc,
		LineSize:       lineSize,
		LineSizeInLog2: uint32(math.Log2(float64(lineSize))),
		NumSets:        size / assoc / lineSize,
		NumSetsInLog2:  uint32(math.Log2(float64(size / assoc / lineSize))),
		NumLines:       size / lineSize,
	}

	return geometry
}

func (geometry *Geometry) GetDisplacement(address uint32) uint32 {
	return address & (geometry.LineSize - 1)
}

func (geometry *Geometry) GetTag(address uint32) uint32 {
	return address & ^(geometry.LineSize - 1)
}

func (geometry *Geometry) GetLineId(address uint32) uint32 {
	return address >> geometry.LineSizeInLog2
}

func (geometry *Geometry) GetSet(address uint32) uint32 {
	return geometry.GetLineId(address) % geometry.NumSets
}

func (geometry *Geometry) IsAligned(address uint32) bool {
	return geometry.GetDisplacement(address) == 0
}
