package mem

import (
	"testing"
	"fmt"
)

func TestPagedMemory(t *testing.T) {
	var memory Memory = NewPagedMemory(true)

	memory.WriteStringAt(12, "你好 world.")

	fmt.Printf("%s\n", memory.ReadStringAt(12, uint32(len([]byte("Hello world.")))))

	memory.WriteWordAt(1, 12)

	fmt.Println(memory.ReadWordAt(1))
}
