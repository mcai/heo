package mem

import (
	"fmt"
	"testing"
)

func TestCircularByteBuffer(t *testing.T) {
	var buffer = NewCircularByteBuffer(1024)

	var str = "你好 测试."
	var str2 = "Hello test."

	buffer.Write([]byte(str))
	buffer.Write([]byte(str2))

	var str3 = buffer.Read(uint32(len([]byte(str))))
	var str4 = buffer.Read(uint32(len([]byte(str2))))

	fmt.Printf("str: %s, str2: %s, str3: %s, str4: %s\n", str, str2, str3, str4)
}
