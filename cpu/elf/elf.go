package elf

import (
	"io/ioutil"
	"fmt"
	"encoding/binary"
	"bytes"
	"github.com/mcai/heo/cpu/mem"
)

type ElfFile struct {
	Data                 *mem.SimpleMemory
	Identification       *ElfIdentification
	Header               *ElfHeader
	SectionHeaders       []*ElfSectionHeader
	ProgramHeaders       []*ElfProgramHeader
	StringTable          *ElfStringTable
	Symbols              map[uint32]*Symbol
	LocalFunctionSymbols map[uint32]*Symbol
	LocalObjectSymbols   map[uint32]*Symbol
	CommonObjectSymbols  map[uint32]*Symbol
}

func NewElfFile(fileName string) *ElfFile {
	var elfFile = &ElfFile{
	}

	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(fmt.Sprintf("Cannot read ELF file (%s)", err))
	}

	elfFile.Data = mem.NewSimpleMemory(false, data)

	elfFile.Identification = NewElfIdentification(elfFile)

	if elfFile.Identification.Clz != ElfClass32 {
		panic("ElfClass32 is not supported")
	}

	if elfFile.Identification.Data == ElfData2Lsb {
		elfFile.Data.LittleEndian = true
		elfFile.Data.ByteOrder = binary.LittleEndian
	} else {
		elfFile.Data.LittleEndian = false
		elfFile.Data.ByteOrder = binary.BigEndian
	}

	elfFile.Header = NewElfHeader(elfFile)

	if elfFile.Header.Machine != EM_MIPS {
		panic("Non-MIPS ELF file is not supported")
	}

	for i := uint16(0); i < elfFile.Header.SectionHeaderTableEntryCount; i++ {
		elfFile.Data.ReadPosition = elfFile.Header.SectionHeaderTableOffset +
			uint32(i * elfFile.Header.SectionHeaderTableEntrySize)
		elfFile.SectionHeaders = append(elfFile.SectionHeaders, NewElfSectionHeader(elfFile))
	}

	elfFile.StringTable = NewElfStringTable(elfFile, elfFile.SectionHeaders[elfFile.Header.SectionHeaderStringTableIndex])

	elfFile.Data.ReadPosition = elfFile.Header.ProgramHeaderTableOffset

	for i := uint16(0); i < elfFile.Header.ProgramHeaderTableEntryCount; i++ {
		elfFile.ProgramHeaders = append(elfFile.ProgramHeaders, NewElfProgramHeader(elfFile))
	}

	elfFile.Symbols = make(map[uint32]*Symbol)
	elfFile.LocalFunctionSymbols = make(map[uint32]*Symbol)
	elfFile.LocalObjectSymbols = make(map[uint32]*Symbol)
	elfFile.CommonObjectSymbols = make(map[uint32]*Symbol)

	elfFile.loadSymbols()

	return elfFile
}

func (elfFile *ElfFile) loadSymbols() {
	for _, sectionHeader := range elfFile.SectionHeaders {
		if sectionHeader.HeaderType == SHT_SYMTAB {
			elfFile.loadSymbolsBySection(sectionHeader)
		}
	}

	elfFile.loadLocalFunctions()
	elfFile.loadLocalObjects()
	elfFile.loadCommonObjects()
}

func (elfFile *ElfFile) loadSymbolsBySection(elfSectionHeader *ElfSectionHeader) {
	var numSymbols = uint32(1)

	if elfSectionHeader.EntrySize != 0 {
		numSymbols = elfSectionHeader.Size / elfSectionHeader.EntrySize
	}

	var offset = elfSectionHeader.Offset

	for i := uint32(0); i < numSymbols; i++ {
		elfFile.Data.ReadPosition = offset

		var symbol = NewSymbol(elfFile, elfSectionHeader)

		elfFile.Symbols[symbol.Value] = symbol

		offset += elfSectionHeader.EntrySize
	}
}

func (elfFile *ElfFile) GetSymbolAt(address uint32) *Symbol {
	for _, symbol := range elfFile.Symbols {
		if symbol.Value == address {
			return symbol
		}
	}

	return nil
}

func (elfFile *ElfFile) loadLocalFunctions() {
	for _, symbol := range elfFile.Symbols {
		if symbol.GetSymbolType() == STT_FUNC {
			var idx = symbol.SectionHeaderTableIndex
			if idx > SHN_LOPROC && idx < SHN_HIPROC {
				if len(symbol.GetName(elfFile)) > 0 {
					elfFile.LocalFunctionSymbols[symbol.Value] = symbol
				}
			} else if idx >= 0 && elfFile.SectionHeaders[idx].HeaderType != SHT_NULL {
				elfFile.LocalFunctionSymbols[symbol.Value] = symbol
			}
		}
	}
}

func (elfFile *ElfFile) loadLocalObjects() {
	for _, symbol := range elfFile.Symbols {
		if symbol.GetSymbolType() == STT_OBJECT {
			var idx = symbol.SectionHeaderTableIndex
			if idx > SHN_LOPROC && idx < SHN_HIPROC {
				if len(symbol.GetName(elfFile)) > 0 {
					elfFile.LocalObjectSymbols[symbol.Value] = symbol
				}
			} else if idx >= 0 && elfFile.SectionHeaders[idx].HeaderType != SHT_NULL {
				elfFile.LocalObjectSymbols[symbol.Value] = symbol
			}
		}
	}
}

func (elfFile *ElfFile) loadCommonObjects() {
	for _, symbol := range elfFile.Symbols {
		if symbol.GetBind() == STB_GLOBAL && symbol.GetSymbolType() == STT_OBJECT {
			var idx = symbol.SectionHeaderTableIndex
			if idx == SHN_COMMON {
				elfFile.CommonObjectSymbols[symbol.Value] = symbol
			}
		}
	}
}

func (elfFile *ElfFile) Dump() {
	fmt.Printf("Clz: %s, data: %s\n", elfFile.Identification.Clz, elfFile.Identification.Data)

	for i, sectionHeader := range elfFile.SectionHeaders {
		fmt.Printf("sectionHeader[%d].Type = 0x%08x\n", i, sectionHeader.HeaderType)
	}

	for i, programHeader := range elfFile.ProgramHeaders {
		fmt.Printf("programHeader[%d].Type = 0x%08x\n", i, programHeader.HeaderType)
	}

	for i, symbol := range elfFile.Symbols {
		fmt.Printf("symbol[0x%08x].Type = 0x%08x\n", i, symbol.GetSymbolType())
		fmt.Printf("symbol[0x%08x].Name = %s\n", i, symbol.GetName(elfFile))
	}
}

type ElfClass string

const (
	ElfClassNone ElfClass = "ElfClassNone"
	ElfClass32 ElfClass = "ElfClass32"
	ElfClass64 ElfClass = "ElfClass64"
)

type ElfData string

const (
	ElfDataNone ElfData = "ElfDataNone"
	ElfData2Lsb ElfData = "ElfData2Lsb"
	ElfData2Msb ElfData = "ElfData2Msb"
)

type ElfIdentification struct {
	Clz  ElfClass
	Data ElfData
}

func NewElfIdentification(elfFile *ElfFile) *ElfIdentification {
	var elfIdentification = &ElfIdentification{
	}

	var eIdent = elfFile.Data.ReadBlock(16)

	if !(eIdent[0] == 0x7f && eIdent[1] == byte('E') && eIdent[2] == byte('L') && eIdent[3] == byte('F')) {
		panic("Not ELF file")
	}

	switch eIdent[4] {
	case 1:
		elfIdentification.Clz = ElfClass32
	case 2:
		elfIdentification.Clz = ElfClass64
	default:
		elfIdentification.Clz = ElfClassNone
	}

	switch eIdent[5] {
	case 1:
		elfIdentification.Data = ElfData2Lsb
	case 2:
		elfIdentification.Data = ElfData2Msb
	default:
		elfIdentification.Data = ElfDataNone
	}

	return elfIdentification
}

const (
	EM_MIPS uint16 = 8
)

type ElfHeader struct {
	HeaderType                    uint16
	Machine                       uint16
	Version                       uint32
	Entry                         uint32
	ProgramHeaderTableOffset      uint32
	SectionHeaderTableOffset      uint32
	Flags                         uint32
	ElfHeaderSize                 uint16
	ProgramHeaderTableEntrySize   uint16
	ProgramHeaderTableEntryCount  uint16
	SectionHeaderTableEntrySize   uint16
	SectionHeaderTableEntryCount  uint16
	SectionHeaderStringTableIndex uint16
}

func NewElfHeader(elfFile *ElfFile) *ElfHeader {
	var header = &ElfHeader{
	}

	header.HeaderType = elfFile.Data.ReadHalfWord()

	header.Machine = elfFile.Data.ReadHalfWord()
	header.Version = elfFile.Data.ReadWord()
	header.Entry = elfFile.Data.ReadWord()
	header.ProgramHeaderTableOffset = elfFile.Data.ReadWord()
	header.SectionHeaderTableOffset = elfFile.Data.ReadWord()
	header.Flags = elfFile.Data.ReadWord()

	header.ElfHeaderSize = elfFile.Data.ReadHalfWord()
	header.ProgramHeaderTableEntrySize = elfFile.Data.ReadHalfWord()
	header.ProgramHeaderTableEntryCount = elfFile.Data.ReadHalfWord()
	header.SectionHeaderTableEntrySize = elfFile.Data.ReadHalfWord()
	header.SectionHeaderTableEntryCount = elfFile.Data.ReadHalfWord()
	header.SectionHeaderStringTableIndex = elfFile.Data.ReadHalfWord()

	return header
}

type ElfSectionHeaderType uint32

const (
	SHT_NULL ElfSectionHeaderType = 0
	SHT_PROGBITS ElfSectionHeaderType = 1
	SHT_SYMTAB ElfSectionHeaderType = 2
	SHT_STRTAB ElfSectionHeaderType = 3
	SHT_RELA ElfSectionHeaderType = 4
	SHT_HASH ElfSectionHeaderType = 5
	SHT_DYNAMIC ElfSectionHeaderType = 6
	SHT_NOTE ElfSectionHeaderType = 7
	SHT_NOBITS ElfSectionHeaderType = 8
	SHT_REL ElfSectionHeaderType = 9
	SHT_SHLIB ElfSectionHeaderType = 10
	SHT_DYNSYM ElfSectionHeaderType = 11
)

type ElfSectionHeaderFlag uint32

const (
	SHF_WRITE ElfSectionHeaderFlag = 0x1
	SHF_ALLOC ElfSectionHeaderFlag = 0x2
	SHF_EXECINSTR ElfSectionHeaderFlag = 0x4
)

type ElfSectionHeader struct {
	NameIndex        uint32
	HeaderType       ElfSectionHeaderType
	Flags            uint32
	Address          uint32
	Offset           uint32
	Size             uint32
	Link             uint32
	Info             uint32
	AddressAlignment uint32
	EntrySize        uint32
	name             string
}

func NewElfSectionHeader(elfFile *ElfFile) *ElfSectionHeader {
	var elfSectionHeader = &ElfSectionHeader{
	}

	elfSectionHeader.NameIndex = elfFile.Data.ReadWord()
	elfSectionHeader.HeaderType = ElfSectionHeaderType(elfFile.Data.ReadWord())
	elfSectionHeader.Flags = elfFile.Data.ReadWord()
	elfSectionHeader.Address = elfFile.Data.ReadWord()
	elfSectionHeader.Offset = elfFile.Data.ReadWord()
	elfSectionHeader.Size = elfFile.Data.ReadWord()
	elfSectionHeader.Link = elfFile.Data.ReadWord()
	elfSectionHeader.Info = elfFile.Data.ReadWord()
	elfSectionHeader.AddressAlignment = elfFile.Data.ReadWord()
	elfSectionHeader.EntrySize = elfFile.Data.ReadWord()

	return elfSectionHeader
}

func (elfSectionHeader *ElfSectionHeader) ReadContent(elfFile *ElfFile) []byte {
	return elfFile.Data.ReadBlockAt(elfSectionHeader.Offset, elfSectionHeader.Size)
}

func (elfSectionHeader *ElfSectionHeader) GetName(elfFile *ElfFile) string {
	if elfSectionHeader.name == "" {
		elfSectionHeader.name = elfFile.StringTable.GetString(elfSectionHeader.NameIndex)
	}

	return elfSectionHeader.name
}

type ElfProgramHeader struct {
	HeaderType      uint32
	Offset          uint32
	VirtualAddress  uint32
	PhysicalAddress uint32
	SizeInFile      uint32
	SizeInMemory    uint32
	Flags           uint32
	Alignment       uint32
}

func NewElfProgramHeader(elfFile *ElfFile) *ElfProgramHeader {
	var elfProgramHeader = &ElfProgramHeader{
	}

	elfProgramHeader.HeaderType = elfFile.Data.ReadWord()
	elfProgramHeader.Offset = elfFile.Data.ReadWord()
	elfProgramHeader.VirtualAddress = elfFile.Data.ReadWord()
	elfProgramHeader.PhysicalAddress = elfFile.Data.ReadWord()
	elfProgramHeader.SizeInFile = elfFile.Data.ReadWord()
	elfProgramHeader.SizeInMemory = elfFile.Data.ReadWord()
	elfProgramHeader.Flags = elfFile.Data.ReadWord()
	elfProgramHeader.Alignment = elfFile.Data.ReadWord()

	return elfProgramHeader
}

func (elfProgramHeader *ElfProgramHeader) ReadContent(elfFile *ElfFile) []byte {
	return elfFile.Data.ReadBlockAt(elfProgramHeader.Offset, elfProgramHeader.SizeInFile)
}

type ElfStringTable struct {
	Data []byte
}

func NewElfStringTable(elfFile *ElfFile, sectionHeader *ElfSectionHeader) *ElfStringTable {
	if sectionHeader.HeaderType != SHT_STRTAB {
		panic("Section is not a string table")
	}

	var elfStringTable = &ElfStringTable{
	}

	elfStringTable.Data = sectionHeader.ReadContent(elfFile)

	return elfStringTable
}

func (elfStringTable *ElfStringTable) GetString(index uint32) string {
	var buf bytes.Buffer

	for i := index; elfStringTable.Data[i] != byte('\x00'); i++ {
		buf.WriteByte(elfStringTable.Data[i])
	}

	return buf.String()
}

const (
	STB_LOCAL = 0

	STB_GLOBAL = 1

	STB_WEAK = 2

	STT_NOTYPE = 0

	STT_OBJECT = 1

	STT_FUNC = 2

	STT_SECTION = 3

	STT_FILE = 4
)

const (
	SHN_UNDEF = 0

	SHN_LORESERVE = 0xff00

	SHN_LOPROC = 0xff00

	SHN_HIPROC = 0xff1f

	SHN_LOOS = 0xff20

	SHN_HIOS = 0xff3f

	SHN_ABS = 0xfff1

	SHN_COMMON = 0xfff2

	SHN_XINDEX = 0xffff

	SHN_HIRESERVE = 0xffff
)

type Symbol struct {
	NameIndex               uint32
	Value                   uint32
	Size                    uint32
	Info                    byte
	Other                   byte
	SectionHeaderTableIndex uint16
	name                    string
	SymbolSectionHeader     *ElfSectionHeader
}

func NewSymbol(elfFile *ElfFile, symbolSectionHeader *ElfSectionHeader) *Symbol {
	var symbol = &Symbol{
		SymbolSectionHeader:symbolSectionHeader,
	}

	symbol.NameIndex = elfFile.Data.ReadWord()
	symbol.Value = elfFile.Data.ReadWord()
	symbol.Size = elfFile.Data.ReadWord()
	symbol.Info = elfFile.Data.ReadByte()
	symbol.Other = elfFile.Data.ReadByte()
	symbol.SectionHeaderTableIndex = elfFile.Data.ReadHalfWord()

	return symbol
}

func (symbol *Symbol) GetSymbolType() byte {
	return symbol.Info & 0xf
}

func (symbol *Symbol) GetBind() byte {
	return (symbol.Info >> 4) & 0xf
}

func (symbol *Symbol) GetName(elfFile *ElfFile) string {
	if symbol.name == "" {
		var elfSectionHeader = elfFile.SectionHeaders[symbol.SymbolSectionHeader.Link]
		symbol.name = NewElfStringTable(elfFile, elfSectionHeader).GetString(symbol.NameIndex)
	}

	return symbol.name
}