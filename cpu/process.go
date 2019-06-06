package cpu

import (
	"fmt"
	"github.com/mcai/heo/cpu/cpuutil"
	"github.com/mcai/heo/cpu/elf"
	"github.com/mcai/heo/cpu/mem"
	"math"
	"strings"
	"syscall"
)

const (
	TEXT_BASE = 0x00400000

	DATA_BASE = 0x10000000

	STACK_BASE = 0x7fffc000

	MAX_ENVIRON = 16 * 1024
)

type Process struct {
	Kernel *Kernel

	Id int32

	ContextMapping *ContextMapping

	Environments []string

	StdInFileDescriptor  int32
	StdOutFileDescriptor int32

	StackBase       uint32
	StackSize       uint32
	TextSize        uint32
	EnvironmentBase uint32
	HeapTop         uint32
	DataTop         uint32
	ProgramEntry    uint32

	LittleEndian bool

	memory            *mem.PagedMemory
	speculativeMemory *mem.PagedMemory

	Speculative bool

	pcToMachInsts          map[uint32]MachInst
	machInstsToStaticInsts map[MachInst]*StaticInst
}

func NewProcess(kernel *Kernel, contextMapping *ContextMapping) *Process {
	var process = &Process{
		Kernel: kernel,

		Id: kernel.CurrentProcessId,

		ContextMapping: contextMapping,

		StdInFileDescriptor:  0,
		StdOutFileDescriptor: 1,

		memory: mem.NewPagedMemory(false),
	}

	kernel.CurrentProcessId++

	kernel.Processes = append(kernel.Processes, process)

	process.LoadProgram(kernel, contextMapping)

	return process
}

func (process *Process) LoadProgram(kernel *Kernel, contextMapping *ContextMapping) {
	process.pcToMachInsts = make(map[uint32]MachInst)
	process.machInstsToStaticInsts = make(map[MachInst]*StaticInst)

	var cmdArgs = strings.Split(contextMapping.Executable+" "+contextMapping.Arguments, " ")

	var elfFileName = cmdArgs[0]

	var elfFile = elf.NewElfFile(elfFileName)

	for _, sectionHeader := range elfFile.SectionHeaders {
		if sectionHeader.GetName(elfFile) == ".dynamic" {
			panic("dynamic linking is not supported")
		}

		if sectionHeader.HeaderType == elf.SHT_PROGBITS || sectionHeader.HeaderType == elf.SHT_NOBITS {
			if sectionHeader.Size > 0 && (sectionHeader.Flags&uint32(elf.SHF_ALLOC)) != 0 {
				if sectionHeader.HeaderType == elf.SHT_NOBITS {
					process.memory.Zero(sectionHeader.Address, sectionHeader.Size)
				} else {
					process.memory.WriteBlockAt(sectionHeader.Address, sectionHeader.Size, sectionHeader.ReadContent(elfFile))

					if sectionHeader.Flags&uint32(elf.SHF_EXECINSTR) != 0 {
						for i := uint32(0); i < sectionHeader.Size; i += 4 {
							var pc = sectionHeader.Address + i
							process.predecode(pc)
						}
					}
				}

				if sectionHeader.Address >= DATA_BASE {
					process.DataTop = uint32(math.Max(float64(process.DataTop), float64(sectionHeader.Address+sectionHeader.Size-1)))
				}
			}
		}

		if sectionHeader.GetName(elfFile) == ".text" {
			process.TextSize = sectionHeader.Address + sectionHeader.Size - TEXT_BASE
		}
	}

	process.ProgramEntry = elfFile.Header.Entry
	process.HeapTop = cpuutil.RoundUp(process.DataTop, process.memory.PageSize())

	process.StackBase = STACK_BASE
	process.StackSize = MAX_ENVIRON
	process.EnvironmentBase = STACK_BASE - MAX_ENVIRON

	process.memory.Zero(process.StackBase-process.StackSize, process.StackSize)

	var stackPointer = process.EnvironmentBase
	process.memory.WriteWordAt(stackPointer, uint32(len(cmdArgs)))
	stackPointer += 4

	var argAddr = stackPointer
	stackPointer += (uint32(len(cmdArgs)) + 1) * 4

	var environmentAddr = stackPointer
	stackPointer += (uint32(len(process.Environments)) + 1) * 4

	for i := uint32(0); i < uint32(len(cmdArgs)); i++ {
		process.memory.WriteWordAt(argAddr+i*4, stackPointer)
		process.memory.WriteStringAt(stackPointer, cmdArgs[i])
		stackPointer += uint32(len([]byte(cmdArgs[i] + "\x00")))
	}
	process.memory.WriteWordAt(argAddr+uint32(len(cmdArgs))*4, 0)

	for i := uint32(0); i < uint32(len(process.Environments)); i++ {
		process.memory.WriteWordAt(environmentAddr+i*4, stackPointer)
		process.memory.WriteStringAt(stackPointer, process.Environments[i])
		stackPointer += uint32(len([]byte(process.Environments[i] + "\x00")))
	}
	process.memory.WriteWordAt(environmentAddr+uint32(len(process.Environments))*4, 0)

	if stackPointer > process.StackBase {
		panic("'environ' overflow, increment MAX_ENVIRON")
	}
}

func (process *Process) Memory() *mem.PagedMemory {
	if process.Speculative {
		return process.speculativeMemory
	} else {
		return process.memory
	}
}

func (process *Process) EnterSpeculativeState() {
	process.speculativeMemory = process.memory.Clone()

	process.Speculative = true
}

func (process *Process) ExitSpeculativeState() {
	process.speculativeMemory = nil

	process.Speculative = false
}

func (process *Process) TranslateFileDescriptor(fileDescriptor int32) int32 {
	if fileDescriptor == 1 || fileDescriptor == 2 {
		return process.StdOutFileDescriptor
	} else if fileDescriptor == 0 {
		return process.StdInFileDescriptor
	} else {
		return fileDescriptor
	}
}

func (process *Process) CloseProgram() {
	if process.StdInFileDescriptor != 0 {
		syscall.Close(int(process.StdInFileDescriptor))
	}

	if process.StdOutFileDescriptor > 2 {
		syscall.Close(int(process.StdOutFileDescriptor))
	}
}

func (process *Process) decode(machInst MachInst) *StaticInst {
	for _, mnemonic := range process.Kernel.Experiment.ISA.Mnemonics {
		if (uint32(machInst)&mnemonic.Mask) == mnemonic.Bits && (mnemonic.ExtraBitField == nil || machInst.ValueOf(mnemonic.ExtraBitField) == mnemonic.ExtraBitFieldValue) {
			return NewStaticInst(mnemonic, machInst)
		}
	}

	panic(fmt.Sprintf("Cannot decode machInst: 0x%08x", machInst))
}

func (process *Process) predecode(pc uint32) {
	var machInst = MachInst(process.memory.ReadWordAt(pc))

	process.pcToMachInsts[pc] = machInst

	if _, ok := process.machInstsToStaticInsts[machInst]; !ok {
		var staticInst = process.decode(machInst)
		process.machInstsToStaticInsts[machInst] = staticInst
	}
}

func (process *Process) GetStaticInst(pc uint32) *StaticInst {
	return process.machInstsToStaticInsts[process.pcToMachInsts[pc]]
}

func (process *Process) Disassemble(pc uint32) string {
	var staticInst = process.GetStaticInst(pc)

	return Disassemble(pc, string(staticInst.Mnemonic.Name), staticInst.MachInst)
}
