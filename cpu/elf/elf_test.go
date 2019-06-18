package elf

import (
	"testing"
)

func TestElfFile(t *testing.T) {
	var elfFile = NewElfFile(
		"Data/Benchmarks/Olden_Custom1/mst/baseline/mst.mips")

	elfFile.Dump()
}
