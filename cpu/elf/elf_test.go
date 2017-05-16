package elf

import (
	"testing"
)

func TestElfFile(t *testing.T) {
	var elfFile = NewElfFile(
		"/home/itecgo/Projects/Archimulator/benchmarks/Olden_Custom1/mst/baseline/mst.mips")

	elfFile.Dump()
}
