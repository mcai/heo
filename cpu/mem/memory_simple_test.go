package mem

import (
	"fmt"
	"testing"
)

func TestSimpleMemory(t *testing.T) {
	var data = make([]byte, 1024)

	var memory = NewSimpleMemory(true, data)

	memory.WriteStringAt(12, "你好 world.")

	fmt.Printf("%s\n", memory.ReadStringAt(12, uint32(len([]byte("Hello world.")))))

	memory.WriteWordAt(1, 12)

	fmt.Println(memory.ReadWordAt(1))

	memory.WriteString("Hello 世界.")

	fmt.Printf("%s\n", memory.ReadString(uint32(len([]byte("Hello 世界.")))))
}
