package mem

type Memory interface {
	ReadByteAt(virtualAddress uint32) byte
	ReadHalfWordAt(virtualAddress uint32) uint16
	ReadWordAt(virtualAddress uint32) uint32
	ReadDoubleWordAt(virtualAddress uint32) uint64
	ReadBlockAt(virtualAddress uint32, size uint32) []byte
	ReadStringAt(virtualAddress uint32, size uint32) string
	WriteByteAt(virtualAddress uint32, data byte)
	WriteHalfWordAt(virtualAddress uint32, data uint16)
	WriteWordAt(virtualAddress uint32, data uint32)
	WriteDoubleWordAt(virtualAddress uint32, data uint64)
	WriteStringAt(virtualAddress uint32, data string)
	WriteBlockAt(virtualAddress uint32, size uint32, data []byte)
}

type MemoryReader interface {
	ReadByte() byte
	ReadHalfWord() uint16
	ReadWord() uint32
	ReadDoubleWord() uint64
	ReadString(size uint32) string
	ReadBlock(size uint32) []byte
}

type MemoryWriter interface {
	WriteByte(data byte)
	WriteHalfWord(data uint16)
	WriteWord(data uint32)
	WriteDoubleWord(data uint64)
	WriteString(data string)
	WriteBlock(size uint32, data []byte)
}
