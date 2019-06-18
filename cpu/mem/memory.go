package mem

type Memory interface {
	ReadUInt8At(virtualAddress uint32) uint8
	ReadUInt16At(virtualAddress uint32) uint16
	ReadUInt32At(virtualAddress uint32) uint32
	ReadUInt64At(virtualAddress uint32) uint64
	ReadBlockAt(virtualAddress uint32, size uint32) []uint8
	ReadStringAt(virtualAddress uint32, size uint32) string
	WriteUInt8At(virtualAddress uint32, data uint8)
	WriteUInt16At(virtualAddress uint32, data uint16)
	WriteUInt32At(virtualAddress uint32, data uint32)
	WriteUInt64At(virtualAddress uint32, data uint64)
	WriteStringAt(virtualAddress uint32, data string)
	WriteBlockAt(virtualAddress uint32, size uint32, data []uint8)
}

type MemoryReader interface {
	ReadUInt8() uint8
	ReadUInt16() uint16
	ReadUInt32() uint32
	ReadUInt64() uint64
	ReadString(size uint32) string
	ReadBlock(size uint32) []uint8
}

type MemoryWriter interface {
	WriteUInt8(data uint8)
	WriteUInt16(data uint16)
	WriteUInt32(data uint32)
	WriteUInt64(data uint64)
	WriteString(data string)
	WriteBlock(size uint32, data []uint8)
}
