package mem

import (
	"bytes"
	"math"
)

type CircularByteBuffer struct {
	data  chan byte
	Size  uint32
	Count uint32
}

func NewCircularByteBuffer(size uint32) *CircularByteBuffer {
	var buffer = &CircularByteBuffer{
		data:make(chan byte, size),
		Size:size,
	}

	return buffer
}

func (buffer *CircularByteBuffer) Read(count uint32) []byte {
	var buf bytes.Buffer

	if count > buffer.Count {
		panic("Requested read is greater than the buffer")
	}

	count = uint32(math.Min(float64(count), float64(buffer.Count)))

	for i := uint32(0); i < count; i++ {
		buf.WriteByte(<-buffer.data)
	}

	buffer.Count -= count

	return buf.Bytes()
}

func (buffer *CircularByteBuffer) Write(src []byte) {
	if uint32(len(src)) > buffer.Size - buffer.Count {
		panic("Requested write is greater than the buffer")
	}

	for _, b := range src {
		buffer.data <- b
	}

	buffer.Count += uint32(len(src))
}

func (buffer *CircularByteBuffer) IsEmpty() bool {
	return buffer.Count == 0
}
