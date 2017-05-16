package regs

import (
	"encoding/binary"
	"math"
	"bytes"
	"fmt"
)

var GPR_NAMES = []string{
	"zero", "at", "v0", "v1", "a0", "a1", "a2", "a3",
	"t0", "t1", "t2", "t3", "t4", "t5", "t6", "t6",
	"s0", "s1", "s2", "s3", "s4", "s5", "s6", "s7",
	"t8", "t9", "k0", "k1", "gp", "sp", "fp", "ra",
}

const (
	NUM_INT_REGISTERS = 32

	NUM_FP_REGISTERS = 32

	NUM_MISC_REGISTERS = 3

	REGISTER_ZERO = 0

	REGISTER_AT = 1

	REGISTER_V0 = 2

	REGISTER_V1 = 3

	REGISTER_A0 = 4

	REGISTER_A1 = 5

	REGISTER_A2 = 6

	REGISTER_A3 = 7

	REGISTER_T0 = 8

	REGISTER_T1 = 9

	REGISTER_T2 = 10

	REGISTER_T3 = 11

	REGISTER_T4 = 12

	REGISTER_T5 = 13

	REGISTER_T6 = 14

	REGISTER_T7 = 15

	REGISTER_S0 = 16

	REGISTER_S1 = 17

	REGISTER_S2 = 18

	REGISTER_S3 = 19

	REGISTER_S4 = 20

	REGISTER_S5 = 21

	REGISTER_S6 = 22

	REGISTER_S7 = 23

	REGISTER_T8 = 24

	REGISTER_T9 = 25

	REGISTER_K0 = 26

	REGISTER_K1 = 27

	REGISTER_GP = 28

	REGISTER_SP = 29

	REGISTER_FP = 30

	REGISTER_RA = 31

	REGISTER_LO = 0

	REGISTER_HI = 1

	REGISTER_FCSR = 2
)

type ArchitecturalRegisterFile struct {
	LittleEndian bool
	Pc           uint32
	Npc          uint32
	Nnpc         uint32
	Gpr          []uint32
	Fpr          *FloatingPointRegisters
	Hi           uint32
	Lo           uint32
	Fcsr         uint32
}

func NewArchitecturalRegisterFile(littleEndian bool) *ArchitecturalRegisterFile {
	var regs = &ArchitecturalRegisterFile{
		LittleEndian:littleEndian,
		Gpr:make([]uint32, 32),
		Fpr:NewFloatingPointRegisters(littleEndian),
	}

	return regs
}

func (regs *ArchitecturalRegisterFile) Clone() *ArchitecturalRegisterFile {
	var newArchitecturalRegisterFile = NewArchitecturalRegisterFile(regs.LittleEndian)

	newArchitecturalRegisterFile.Pc = regs.Pc
	newArchitecturalRegisterFile.Npc = regs.Npc
	newArchitecturalRegisterFile.Nnpc = regs.Nnpc

	copy(newArchitecturalRegisterFile.Gpr, regs.Gpr)

	newArchitecturalRegisterFile.Fpr = NewFloatingPointRegisters(regs.LittleEndian)
	copy(newArchitecturalRegisterFile.Fpr.data, regs.Fpr.data)

	newArchitecturalRegisterFile.Hi = regs.Hi
	newArchitecturalRegisterFile.Lo = regs.Lo
	newArchitecturalRegisterFile.Fcsr = regs.Fcsr

	return newArchitecturalRegisterFile
}

func (regs *ArchitecturalRegisterFile) Sgpr(i uint32) int32 {
	return int32(regs.Gpr[i])
}

func (regs *ArchitecturalRegisterFile) SetSgpr(i uint32, v int32) {
	regs.Gpr[i] = uint32(v)
}

func (regs *ArchitecturalRegisterFile) Dump() string {
	var buf bytes.Buffer

	for i := 0; i < 32; i++ {
		buf.WriteString(fmt.Sprintf("%s = 0x%08x, \n", GPR_NAMES[i], regs.Gpr[i]))
	}

	buf.WriteString(
		fmt.Sprintf("pc = 0x%08x, npc = 0x%08x, nnpc = 0x%08x, hi = 0x%08x, lo = 0x%08x, fcsr = 0x%08x",
			regs.Pc,
			regs.Npc,
			regs.Nnpc,
			regs.Hi,
			regs.Lo,
			regs.Fcsr))

	return buf.String()
}

type FloatingPointRegisters struct {
	LittleEndian bool
	ByteOrder    binary.ByteOrder
	data         []byte
}

func NewFloatingPointRegisters(littleEndian bool) *FloatingPointRegisters {
	var fprs = &FloatingPointRegisters{
		LittleEndian:littleEndian,
		data:make([]byte, 4 * 32),
	}

	if littleEndian {
		fprs.ByteOrder = binary.LittleEndian
	} else {
		fprs.ByteOrder = binary.BigEndian
	}

	return fprs
}

func (fprs *FloatingPointRegisters) Uint32(index uint32) uint32 {
	var size = uint32(4)

	var buffer = make([]byte, size)

	copy(buffer, fprs.data[index * size:index * size + size])

	return fprs.ByteOrder.Uint32(buffer)
}

func (fprs *FloatingPointRegisters) SetUint32(index uint32, value uint32) {
	var size = uint32(4)

	var buffer = make([]byte, size)

	fprs.ByteOrder.PutUint32(buffer, value)

	copy(fprs.data[index * size:index * size + size], buffer)
}

func (fprs *FloatingPointRegisters) Float32(index uint32) float32 {
	var size = uint32(4)

	var buffer = make([]byte, size)

	copy(buffer, fprs.data[index * size:index * size + size])

	return math.Float32frombits(fprs.ByteOrder.Uint32(buffer))
}

func (fprs *FloatingPointRegisters) SetFloat32(index uint32, value float32) {
	var size = uint32(4)

	var buffer = make([]byte, size)

	fprs.ByteOrder.PutUint32(buffer, math.Float32bits(value))

	copy(fprs.data[index * size:index * size + size], buffer)
}

func (fprs *FloatingPointRegisters) Uint64(index uint32) uint64 {
	var size = uint32(8)

	var buffer = make([]byte, size)

	copy(buffer, fprs.data[index * size:index * size + size])

	return fprs.ByteOrder.Uint64(buffer)
}

func (fprs *FloatingPointRegisters) SetUint64(index uint32, value uint64) {
	var size = uint32(8)

	var buffer = make([]byte, size)

	fprs.ByteOrder.PutUint64(buffer, value)

	copy(fprs.data[index * size:index * size + size], buffer)
}

func (fprs *FloatingPointRegisters) Float64(index uint32) float64 {
	var size = uint32(8)

	var buffer = make([]byte, size)

	copy(buffer, fprs.data[index * size:index * size + size])

	return math.Float64frombits(fprs.ByteOrder.Uint64(buffer))
}

func (fprs *FloatingPointRegisters) SetFloat64(index uint32, value float64) {
	var size = uint32(8)

	var buffer = make([]byte, size)

	fprs.ByteOrder.PutUint64(buffer, math.Float64bits(value))

	copy(fprs.data[index * size:index * size + size], buffer)
}
