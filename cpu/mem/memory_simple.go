package mem

import "encoding/binary"

type SimpleMemory struct {
	LittleEndian  bool
	ByteOrder     binary.ByteOrder
	Data          []byte
	ReadPosition  uint32
	WritePosition uint32
}

func NewSimpleMemory(littleEndian bool, data []byte) *SimpleMemory {
	var memory = &SimpleMemory{
		LittleEndian: littleEndian,
		Data:         data,
	}

	if littleEndian {
		memory.ByteOrder = binary.LittleEndian
	} else {
		memory.ByteOrder = binary.BigEndian
	}

	return memory
}

func (memory *SimpleMemory) ReadByteAt(virtualAddress uint32) byte {
	var buffer = make([]byte, 1)
	memory.access(virtualAddress, 1, &buffer, false)
	return buffer[0]
}

func (memory *SimpleMemory) ReadHalfWordAt(virtualAddress uint32) uint16 {
	var buffer = make([]byte, 2)
	memory.access(virtualAddress, 2, &buffer, false)
	return memory.ByteOrder.Uint16(buffer)
}

func (memory *SimpleMemory) ReadWordAt(virtualAddress uint32) uint32 {
	var buffer = make([]byte, 4)
	memory.access(virtualAddress, 4, &buffer, false)
	return memory.ByteOrder.Uint32(buffer)
}

func (memory *SimpleMemory) ReadDoubleWordAt(virtualAddress uint32) uint64 {
	var buffer = make([]byte, 8)
	memory.access(virtualAddress, 8, &buffer, false)
	return memory.ByteOrder.Uint64(buffer)
}

func (memory *SimpleMemory) ReadBlockAt(virtualAddress uint32, size uint32) []byte {
	var buffer = make([]byte, size)
	memory.access(virtualAddress, size, &buffer, false)
	return buffer
}

func (memory *SimpleMemory) ReadStringAt(virtualAddress uint32, size uint32) string {
	var data = memory.ReadBlockAt(virtualAddress, size)

	var str []byte

	for i := 0; data[i] != byte('\x00'); i++ {
		str = append(str, data[i])
	}

	return string(str)
}

func (memory *SimpleMemory) WriteByteAt(virtualAddress uint32, data byte) {
	var buffer = make([]byte, 1)
	buffer[0] = data
	memory.access(virtualAddress, 1, &buffer, true)
}

func (memory *SimpleMemory) WriteHalfWordAt(virtualAddress uint32, data uint16) {
	var buffer = make([]byte, 2)
	memory.ByteOrder.PutUint16(buffer, data)
	memory.access(virtualAddress, 2, &buffer, true)
}

func (memory *SimpleMemory) WriteWordAt(virtualAddress uint32, data uint32) {
	var buffer = make([]byte, 4)
	memory.ByteOrder.PutUint32(buffer, data)
	memory.access(virtualAddress, 4, &buffer, true)
}

func (memory *SimpleMemory) WriteDoubleWordAt(virtualAddress uint32, data uint64) {
	var buffer = make([]byte, 8)
	memory.ByteOrder.PutUint64(buffer, data)
	memory.access(virtualAddress, 8, &buffer, true)
}

func (memory *SimpleMemory) WriteStringAt(virtualAddress uint32, data string) {
	var buffer = []byte(data + "\x00")
	memory.access(virtualAddress, uint32(len(buffer)), &buffer, true)
}

func (memory *SimpleMemory) WriteBlockAt(virtualAddress uint32, size uint32, data []byte) {
	memory.access(virtualAddress, size, &data, true)
}

func (memory *SimpleMemory) ReadByte() byte {
	var data = memory.ReadByteAt(memory.ReadPosition)
	memory.ReadPosition++
	return data
}

func (memory *SimpleMemory) ReadHalfWord() uint16 {
	var data = memory.ReadHalfWordAt(memory.ReadPosition)
	memory.ReadPosition += 2
	return data
}

func (memory *SimpleMemory) ReadWord() uint32 {
	var data = memory.ReadWordAt(memory.ReadPosition)
	memory.ReadPosition += 4
	return data
}

func (memory *SimpleMemory) ReadDoubleWord() uint64 {
	var data = memory.ReadDoubleWordAt(memory.ReadPosition)
	memory.ReadPosition += 8
	return data
}

func (memory *SimpleMemory) ReadString(size uint32) string {
	var data = memory.ReadStringAt(memory.ReadPosition, size)
	memory.ReadPosition += size
	return data
}

func (memory *SimpleMemory) ReadBlock(size uint32) []byte {
	var data = memory.ReadBlockAt(memory.ReadPosition, size)
	memory.ReadPosition += size
	return data
}

func (memory *SimpleMemory) WriteByte(data byte) {
	memory.WriteByteAt(memory.WritePosition, data)
	memory.WritePosition++
}

func (memory *SimpleMemory) WriteHalfWord(data uint16) {
	memory.WriteHalfWordAt(memory.WritePosition, data)
	memory.WritePosition += 2
}

func (memory *SimpleMemory) WriteWord(data uint32) {
	memory.WriteWordAt(memory.WritePosition, data)
	memory.WritePosition += 4
}

func (memory *SimpleMemory) WriteDoubleWord(data uint64) {
	memory.WriteDoubleWordAt(memory.WritePosition, data)
	memory.WritePosition += 8
}

func (memory *SimpleMemory) WriteString(data string) {
	memory.WriteStringAt(memory.WritePosition, data)
	memory.WritePosition += uint32(len([]byte(data)))
}

func (memory *SimpleMemory) WriteBlock(size uint32, data []byte) {
	memory.WriteBlockAt(memory.WritePosition, size, data)
	memory.WritePosition += size
}

func (memory *SimpleMemory) access(virtualAddress uint32, size uint32, buffer *[]byte, write bool) {
	if write {
		copy(memory.Data[virtualAddress:virtualAddress+size], (*buffer)[0:size])
	} else {
		copy((*buffer)[0:size], memory.Data[virtualAddress:virtualAddress+size])
	}
}
