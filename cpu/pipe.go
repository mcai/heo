package cpu

import "github.com/mcai/heo/cpu/mem"

type Pipe struct {
	FileDescriptors []int32
	Buffer          *mem.CircularByteBuffer
}

func NewPipe(fileDescriptors []int32) *Pipe {
	var pipe = &Pipe{
		FileDescriptors:fileDescriptors,
		Buffer:mem.NewCircularByteBuffer(1024),
	}

	return pipe
}
