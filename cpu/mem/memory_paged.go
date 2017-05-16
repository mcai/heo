package mem

import (
	"math"
	"encoding/binary"
)

type MemoryPage struct {
	Id              uint32
	PhysicalAddress uint32
	Buffer          []byte
}

func NewMemoryPage(memory *PagedMemory, id uint32) *MemoryPage {
	var page = &MemoryPage{
		Id:id,
		PhysicalAddress:id << memory.PageSizeInLog2(),
		Buffer:make([]byte, memory.PageSize()),
	}

	return page
}

func (page *MemoryPage) Clone() *MemoryPage {
	var newPage = &MemoryPage{
		Id:page.Id,
		PhysicalAddress:page.PhysicalAddress,
		Buffer:make([]byte, len(page.Buffer)),
	}

	copy(newPage.Buffer, page.Buffer)

	return newPage
}

func (page *MemoryPage) Access(displacement uint32, buffer *[]byte, offset uint32, size uint32, write bool) {
	if write {
		copy(page.Buffer[displacement:displacement + size], (*buffer)[offset:offset + size])
	} else {
		copy((*buffer)[offset:offset + size], page.Buffer[displacement:displacement + size])
	}

	//fmt.Printf("[mem] vaddr: 0x%08x, size: %d, write: %t\n", virtualAddress, size, write)
}

type PagedMemory struct {
	LittleEndian bool
	ByteOrder    binary.ByteOrder
	Pages        map[uint32]*MemoryPage
	Geometry     *Geometry
	NumPages     uint32
}

func NewPagedMemory(littleEndian bool) *PagedMemory {
	var memory = &PagedMemory{
		LittleEndian:littleEndian,
		Pages:make(map[uint32]*MemoryPage),
		Geometry:NewGeometry(math.MaxUint32, 1, 1 << 12),
	}

	if littleEndian {
		memory.ByteOrder = binary.LittleEndian
	} else {
		memory.ByteOrder = binary.BigEndian
	}

	return memory
}

func (memory *PagedMemory) Clone() *PagedMemory {
	var newMemory = NewPagedMemory(memory.LittleEndian)

	newMemory.NumPages = memory.NumPages

	for index, page := range memory.Pages {
		newMemory.Pages[index] = page.Clone()
	}

	return newMemory
}

func (memory *PagedMemory) ReadByteAt(virtualAddress uint32) byte {
	var buffer = make([]byte, 1)
	memory.access(virtualAddress, 1, &buffer, false, true)
	return buffer[0]
}

func (memory *PagedMemory) ReadHalfWordAt(virtualAddress uint32) uint16 {
	var buffer = make([]byte, 2)
	memory.access(virtualAddress, 2, &buffer, false, true)
	return memory.ByteOrder.Uint16(buffer)
}

func (memory *PagedMemory) ReadWordAt(virtualAddress uint32) uint32 {
	var buffer = make([]byte, 4)
	memory.access(virtualAddress, 4, &buffer, false, true)
	return memory.ByteOrder.Uint32(buffer)
}

func (memory *PagedMemory) ReadDoubleWordAt(virtualAddress uint32) uint64 {
	var buffer = make([]byte, 8)
	memory.access(virtualAddress, 8, &buffer, false, true)
	return memory.ByteOrder.Uint64(buffer)
}

func (memory *PagedMemory) ReadBlockAt(virtualAddress uint32, size uint32) []byte {
	var buffer = make([]byte, size)
	memory.access(virtualAddress, size, &buffer, false, true)
	return buffer
}

func (memory *PagedMemory) ReadStringAt(virtualAddress uint32, size uint32) string {
	var data = memory.ReadBlockAt(virtualAddress, size)

	var str []byte

	for i := 0; data[i] != byte('\x00'); i++ {
		str = append(str, data[i])
	}

	return string(str)
}

func (memory *PagedMemory) WriteByteAt(virtualAddress uint32, data byte) {
	var buffer = make([]byte, 1)
	buffer[0] = data
	memory.access(virtualAddress, 1, &buffer, true, true)

	//fmt.Printf("[mem::WriteByteAt] vaddr: 0x%08x, data: 0x%08x\n", virtualAddress, data)
}

func (memory *PagedMemory) WriteHalfWordAt(virtualAddress uint32, data uint16) {
	var buffer = make([]byte, 2)
	memory.ByteOrder.PutUint16(buffer, data)
	memory.access(virtualAddress, 2, &buffer, true, true)

	//fmt.Printf("[mem::WriteHalfWordAt] vaddr: 0x%08x, data: 0x%08x\n", virtualAddress, data)
}

func (memory *PagedMemory) WriteWordAt(virtualAddress uint32, data uint32) {
	var buffer = make([]byte, 4)
	memory.ByteOrder.PutUint32(buffer, data)
	memory.access(virtualAddress, 4, &buffer, true, true)

	//fmt.Printf("[mem::WriteWordAt] vaddr: 0x%08x, data: 0x%08x\n", virtualAddress, data)
}

func (memory *PagedMemory) WriteDoubleWordAt(virtualAddress uint32, data uint64) {
	var buffer = make([]byte, 8)
	memory.ByteOrder.PutUint64(buffer, data)
	memory.access(virtualAddress, 8, &buffer, true, true)

	//fmt.Printf("[mem::WriteDoubleWordAt] vaddr: 0x%08x, data: 0x%08x\n", virtualAddress, data)
}

func (memory *PagedMemory) WriteStringAt(virtualAddress uint32, data string) {
	var buffer = []byte(data + "\x00")
	memory.access(virtualAddress, uint32(len(buffer)), &buffer, true, true)

	//fmt.Printf("[mem::WriteStringAt] vaddr: 0x%08x, size: %d, data: %s\n", virtualAddress, uint32(len(buffer)), data)
}

func (memory *PagedMemory) WriteBlockAt(virtualAddress uint32, size uint32, data []byte) {
	memory.access(virtualAddress, size, &data, true, true)

	//fmt.Printf("[mem::WriteBlockAt] vaddr: 0x%08x, size: %d\n", virtualAddress, size)
}

func (memory *PagedMemory) Zero(virtualAddress uint32, size uint32) {
	memory.WriteBlockAt(virtualAddress, size, make([]byte, size))
}

func (memory *PagedMemory) Map(virtualAddress uint32, size uint32) uint32 {
	var tagStart, tagEnd uint32

	tagStart = memory.GetTag(virtualAddress)
	tagEnd = tagStart

	var pageSize = memory.PageSize()

	for pageCount := (memory.GetTag(virtualAddress + size - 1) - tagStart) / pageSize + 1; ; {
		if tagEnd == 0 {
			panic("Unimplemented")
			return 0 //TODO
			//return uint32(-1)
		}

		if memory.GetPage(tagEnd) != nil {
			tagEnd += pageSize
			tagStart = tagEnd
			continue
		}

		if (tagEnd - tagStart) / pageSize + 1 == pageCount {
			break
		}

		tagEnd += pageSize
	}

	for tag := tagStart; tag <= tagEnd; tag += pageSize {
		if memory.GetPage(tag) != nil {
			panic("Impossible")
		}
		memory.addPage(tag)
	}

	return tagStart
}

func (memory *PagedMemory) Unmap(virtualAddress uint32, size uint32) {
	var tagStart = memory.GetTag(virtualAddress)
	var tagEnd = memory.GetTag(virtualAddress + size - 1)

	var pageSize = memory.PageSize()

	for tag := tagStart; tag <= tagEnd; tag += pageSize {
		memory.removePage(tag)
	}
}

func (memory *PagedMemory) Remap(oldAddr uint32, oldSize uint32, newSize uint32) uint32 {
	var start = memory.Map(0, newSize)

	if int32(start) != -1 {
		panic("Not supported")
	}

	return start
}

func (memory *PagedMemory) access(virtualAddress uint32, size uint32, buffer *[]byte, write bool, createNewPageIfNecessary bool) {
	var offset = uint32(0)

	var pageSize = memory.PageSize()

	for size > 0 {
		var chunkSize = uint32(math.Min(float64(size), float64(pageSize - memory.GetDisplacement(virtualAddress))))
		memory.accessPageBoundary(virtualAddress, chunkSize, buffer, offset, write, createNewPageIfNecessary)

		size -= chunkSize
		offset += chunkSize
		virtualAddress += chunkSize
	}
}

func (memory *PagedMemory) accessPageBoundary(virtualAddress uint32, size uint32, buffer *[]byte, offset uint32, write bool, createNewPageIfNecessary bool) {
	var page = memory.GetPage(virtualAddress)

	if page == nil && createNewPageIfNecessary {
		page = memory.addPage(memory.GetTag(virtualAddress))
	}

	if page != nil {
		var displacement = memory.GetDisplacement(virtualAddress)

		page.Access(displacement, buffer, offset, size, write)
	}
}

func (memory *PagedMemory) GetPage(virtualAddress uint32) *MemoryPage {
	var index = memory.GetIndex(virtualAddress)

	if page, ok := memory.Pages[index]; ok {
		return page
	}

	return nil
}

func (memory *PagedMemory) addPage(virtualAddress uint32) *MemoryPage {
	var index = memory.GetIndex(virtualAddress)

	var page = NewMemoryPage(memory, memory.NumPages)

	memory.NumPages++

	memory.Pages[index] = page

	return page
}

func (memory *PagedMemory) removePage(virtualAddress uint32) {
	var index = memory.GetIndex(virtualAddress)

	delete(memory.Pages, index)
}

func (memory *PagedMemory) GetPhysicalAddress(virtualAddress uint32) uint32 {
	return memory.GetPage(virtualAddress).PhysicalAddress + memory.GetDisplacement(virtualAddress)
}

func (memory *PagedMemory) GetDisplacement(virtualAddress uint32) uint32 {
	return memory.Geometry.GetDisplacement(virtualAddress)
}

func (memory *PagedMemory) GetTag(virtualAddress uint32) uint32 {
	return memory.Geometry.GetTag(virtualAddress)
}

func (memory *PagedMemory) GetIndex(virtualAddress uint32) uint32 {
	return memory.Geometry.GetLineId(virtualAddress)
}

func (memory *PagedMemory) PageSizeInLog2() uint32 {
	return memory.Geometry.LineSizeInLog2
}

func (memory *PagedMemory) PageSize() uint32 {
	return memory.Geometry.LineSize
}