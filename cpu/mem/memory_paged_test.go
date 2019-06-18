package mem

import (
	"fmt"
	"testing"
)

func TestPagedMemory(t *testing.T) {
	var memory Memory = NewPagedMemory(true)

	memory.WriteStringAt(12, "你好 world.")

	fmt.Printf("%s\n", memory.ReadStringAt(12, uint32(len([]byte("Hello world.")))))

	memory.WriteUInt32At(1, 12)

	fmt.Println(memory.ReadUInt32At(1))
}
