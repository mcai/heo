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

func (memory *SimpleMemory) ReadUInt8At(virtualAddress uint32) uint8 {
	var buffer = make([]byte, 1)
	memory.access(virtualAddress, 1, &buffer, false)
	return buffer[0]
}

func (memory *SimpleMemory) ReadUInt16At(virtualAddress uint32) uint16 {
	var buffer = make([]byte, 2)
	memory.access(virtualAddress, 2, &buffer, false)
	return memory.ByteOrder.Uint16(buffer)
}

func (memory *SimpleMemory) ReadUInt32At(virtualAddress uint32) uint32 {
	var buffer = make([]byte, 4)
	memory.access(virtualAddress, 4, &buffer, false)
	return memory.ByteOrder.Uint32(buffer)
}

func (memory *SimpleMemory) ReadUInt64At(virtualAddress uint32) uint64 {
	var buffer = make([]byte, 8)
	memory.access(virtualAddress, 8, &buffer, false)
	return memory.ByteOrder.Uint64(buffer)
}

func (memory *SimpleMemory) ReadBlockAt(virtualAddress uint32, size uint32) []uint8 {
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

func (memory *SimpleMemory) WriteUInt8At(virtualAddress uint32, data uint8) {
	var buffer = make([]byte, 1)
	buffer[0] = data
	memory.access(virtualAddress, 1, &buffer, true)
}

func (memory *SimpleMemory) WriteUInt16At(virtualAddress uint32, data uint16) {
	var buffer = make([]byte, 2)
	memory.ByteOrder.PutUint16(buffer, data)
	memory.access(virtualAddress, 2, &buffer, true)
}

func (memory *SimpleMemory) WriteUInt32At(virtualAddress uint32, data uint32) {
	var buffer = make([]byte, 4)
	memory.ByteOrder.PutUint32(buffer, data)
	memory.access(virtualAddress, 4, &buffer, true)
}

func (memory *SimpleMemory) WriteUInt64At(virtualAddress uint32, data uint64) {
	var buffer = make([]byte, 8)
	memory.ByteOrder.PutUint64(buffer, data)
	memory.access(virtualAddress, 8, &buffer, true)
}

func (memory *SimpleMemory) WriteStringAt(virtualAddress uint32, data string) {
	var buffer = []byte(data + "\x00")
	memory.access(virtualAddress, uint32(len(buffer)), &buffer, true)
}

func (memory *SimpleMemory) WriteBlockAt(virtualAddress uint32, size uint32, data []uint8) {
	memory.access(virtualAddress, size, &data, true)
}

func (memory *SimpleMemory) ReadUInt8() uint8 {
	var data = memory.ReadUInt8At(memory.ReadPosition)
	memory.ReadPosition++
	return data
}

func (memory *SimpleMemory) ReadUInt16() uint16 {
	var data = memory.ReadUInt16At(memory.ReadPosition)
	memory.ReadPosition += 2
	return data
}

func (memory *SimpleMemory) ReadUInt32() uint32 {
	var data = memory.ReadUInt32At(memory.ReadPosition)
	memory.ReadPosition += 4
	return data
}

func (memory *SimpleMemory) ReadUInt64() uint64 {
	var data = memory.ReadUInt64At(memory.ReadPosition)
	memory.ReadPosition += 8
	return data
}

func (memory *SimpleMemory) ReadString(size uint32) string {
	var data = memory.ReadStringAt(memory.ReadPosition, size)
	memory.ReadPosition += size
	return data
}

func (memory *SimpleMemory) ReadBlock(size uint32) []uint8 {
	var data = memory.ReadBlockAt(memory.ReadPosition, size)
	memory.ReadPosition += size
	return data
}

func (memory *SimpleMemory) WriteUInt8(data uint8) {
	memory.WriteUInt8At(memory.WritePosition, data)
	memory.WritePosition++
}

func (memory *SimpleMemory) WriteUInt16(data uint16) {
	memory.WriteUInt16At(memory.WritePosition, data)
	memory.WritePosition += 2
}

func (memory *SimpleMemory) WriteUInt32(data uint32) {
	memory.WriteUInt32At(memory.WritePosition, data)
	memory.WritePosition += 4
}

func (memory *SimpleMemory) WriteUInt64(data uint64) {
	memory.WriteUInt64At(memory.WritePosition, data)
	memory.WritePosition += 8
}

func (memory *SimpleMemory) WriteString(data string) {
	memory.WriteStringAt(memory.WritePosition, data)
	memory.WritePosition += uint32(len([]byte(data)))
}

func (memory *SimpleMemory) WriteBlock(size uint32, data []uint8) {
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
