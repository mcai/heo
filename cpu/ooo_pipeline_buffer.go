package cpu

type PipelineBuffer struct {
	Size    uint32
	Entries []interface{}
}

func NewPipelineBuffer(size uint32) *PipelineBuffer {
	var buffer = &PipelineBuffer{
		Size: size,
	}

	return buffer
}

func (buffer *PipelineBuffer) Count() uint32 {
	return uint32(len(buffer.Entries))
}

func (buffer *PipelineBuffer) Full() bool {
	return buffer.Count() >= buffer.Size
}

func (buffer *PipelineBuffer) Empty() bool {
	return buffer.Count() == 0
}
