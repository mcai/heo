package cpu

import (
	"bytes"
	"fmt"
	"github.com/mcai/heo/cpu/cpuutil"
	"github.com/mcai/heo/cpu/regs"
	"strings"
)

const (
	MachInstType_R = "r"
	MachInstType_I = "i"
	MachInstType_J = "j"
	MachInstType_F = "f"
)

type MachInstType string

type MachInst uint32

func (machInst MachInst) ValueOf(field *BitField) uint32 {
	return cpuutil.Bits32(uint32(machInst), field.Hi, field.Lo)
}

func (machInst MachInst) GetType() MachInstType {
	var opcode = machInst.ValueOf(OPCODE)

	switch opcode {
	case 0:
		return MachInstType_R
	case 0x02, 0x03:
		return MachInstType_J
	case 0x11:
		return MachInstType_F
	default:
		return MachInstType_I
	}
}

func (machInst MachInst) IsRMt() bool {
	var func_ = machInst.Func()
	return func_ == 0x10 || func_ == 0x11
}

func (machInst MachInst) IsRMf() bool {
	var func_ = machInst.Func()
	return func_ == 0x12 || func_ == 0x13
}

func (machInst MachInst) IsROneOp() bool {
	var func_ = machInst.Func()
	return func_ == 0x08 || func_ == 0x09
}

func (machInst MachInst) IsRTwoOp() bool {
	var func_ = machInst.Func()
	return func_ >= 0x18 && func_ <= 0x1b
}

func (machInst MachInst) IsLoadStore() bool {
	var opcode = machInst.OpCode()
	return opcode >= 0x20 && opcode <= 0x2e || opcode == 0x30 || opcode == 0x38
}

func (machInst MachInst) IsFPLoadStore() bool {
	var opcode = machInst.OpCode()
	return opcode == 0x31 || opcode == 0x39
}

func (machInst MachInst) IsOneOpBranch() bool {
	var opcode = machInst.OpCode()
	return opcode == 0x00 || opcode == 0x01 || opcode == 0x06 || opcode == 0x07
}

func (machInst MachInst) IsShift() bool {
	var func_ = machInst.Func()
	return func_ == 0x00 || func_ == 0x01 || func_ == 0x03
}

func (machInst MachInst) IsCVT() bool {
	var func_ = machInst.Func()
	return func_ == 32 || func_ == 33 || func_ == 36
}

func (machInst MachInst) IsCompare() bool {
	var func_ = machInst.Func()
	return func_ >= 48
}

func (machInst MachInst) IsGprFpMove() bool {
	var rs = machInst.Rs()
	return rs == 0 || rs == 4
}

func (machInst MachInst) IsGprFcrMove() bool {
	var rs = machInst.Rs()
	return rs == 2 || rs == 6
}

func (machInst MachInst) IsFpBranch() bool {
	var rs = machInst.Rs()
	return rs == 8
}

func (machInst MachInst) IsSyscall() bool {
	var opcodeLo = machInst.OpCodeLo()
	var funcHi = machInst.FuncHi()
	var funcLo = machInst.FuncLo()
	return opcodeLo == 0x0 && funcHi == 0x1 && funcLo == 0x4
}

func (machInst MachInst) OpCode() uint32 {
	return machInst.ValueOf(OPCODE)
}

func (machInst MachInst) OpCodeHi() uint32 {
	return machInst.ValueOf(OPCODE_HI)
}

func (machInst MachInst) OpCodeLo() uint32 {
	return machInst.ValueOf(OPCODE_LO)
}

func (machInst MachInst) Rs() uint32 {
	return machInst.ValueOf(RS)
}

func (machInst MachInst) Rt() uint32 {
	return machInst.ValueOf(RT)
}

func (machInst MachInst) Rd() uint32 {
	return machInst.ValueOf(RD)
}

func (machInst MachInst) Shift() uint32 {
	return machInst.ValueOf(SHIFT)
}

func (machInst MachInst) Func() uint32 {
	return machInst.ValueOf(FUNC)
}

func (machInst MachInst) FuncHi() uint32 {
	return machInst.ValueOf(FUNC_HI)
}

func (machInst MachInst) FuncLo() uint32 {
	return machInst.ValueOf(FUNC_LO)
}

func (machInst MachInst) Cond() uint32 {
	return machInst.ValueOf(COND)
}

func (machInst MachInst) Uimm() uint32 {
	return machInst.ValueOf(UIMM)
}

func (machInst MachInst) Imm() int32 {
	return int32(cpuutil.Sext32(machInst.Uimm(), 16))
}

func (machInst MachInst) Target() uint32 {
	return machInst.ValueOf(TARGET)
}

func (machInst MachInst) Fmt() uint32 {
	return machInst.ValueOf(FMT)
}

func (machInst MachInst) Fmt3() uint32 {
	return machInst.ValueOf(FMT3)
}

func (machInst MachInst) Fr() uint32 {
	return machInst.ValueOf(FR)
}

func (machInst MachInst) Fs() uint32 {
	return machInst.ValueOf(FS)
}

func (machInst MachInst) Ft() uint32 {
	return machInst.ValueOf(FT)
}

func (machInst MachInst) Fd() uint32 {
	return machInst.ValueOf(FD)
}

func (machInst MachInst) BranchCc() uint32 {
	return machInst.ValueOf(BRANCH_CC)
}

func (machInst MachInst) Cc() uint32 {
	return machInst.ValueOf(CC)
}

func Disassemble(pc uint32, mnemonicName string, machInst MachInst) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("0x%08x: 0x%08x %s ",
		pc, machInst, strings.ToLower(mnemonicName)))

	if machInst == 0x00000000 {
		return buf.String()
	}

	var machInstType = machInst.GetType()

	var imm = machInst.Imm()

	var rs = machInst.Rs()
	var rt = machInst.Rt()
	var rd = machInst.Rd()

	var fs = machInst.Fs()
	var ft = machInst.Ft()
	var fd = machInst.Fd()

	var shift = machInst.Shift()

	var target = machInst.Target()

	switch machInstType {
	case MachInstType_J:
		buf.WriteString(fmt.Sprintf("%x", target))
	case MachInstType_I:
		if machInst.IsOneOpBranch() {
			buf.WriteString(fmt.Sprintf("$%s, %d", regs.GPR_NAMES[rs], imm))
		} else if machInst.IsLoadStore() {
			buf.WriteString(fmt.Sprintf("$%s, %d($%s)", regs.GPR_NAMES[rt], imm, regs.GPR_NAMES[rs]))
		} else if machInst.IsFPLoadStore() {
			buf.WriteString(fmt.Sprintf("$f%d, %d($%s)", ft, imm, regs.GPR_NAMES[rs]))
		} else {
			buf.WriteString(fmt.Sprintf("$%s, $%s, %d", regs.GPR_NAMES[rt], regs.GPR_NAMES[rs], imm))
		}
	case MachInstType_F:
		if machInst.IsCVT() {
			buf.WriteString(fmt.Sprintf("$f%d, $f%d", fd, fs))
		} else if machInst.IsCompare() {
			buf.WriteString(fmt.Sprintf("%d, $f%d, $f%d", fd>>2, fs, ft))
		} else if machInst.IsFpBranch() {
			buf.WriteString(fmt.Sprintf("%d, %d", fd>>2, imm))
		} else if machInst.IsGprFpMove() {
			buf.WriteString(fmt.Sprintf("$%s, $f%d", regs.GPR_NAMES[rt], fs))
		} else if machInst.IsGprFcrMove() {
			buf.WriteString(fmt.Sprintf("$%s, $%d", regs.GPR_NAMES[rt], fs))
		} else {
			buf.WriteString(fmt.Sprintf("$f%d, $f%d, $f%d", fd, fs, ft))
		}
	case MachInstType_R:
		if !machInst.IsSyscall() {
			if machInst.IsShift() {
				buf.WriteString(fmt.Sprintf("$%s, $%s, %d", regs.GPR_NAMES[rd], regs.GPR_NAMES[rt], shift))
			} else if machInst.IsROneOp() {
				buf.WriteString(fmt.Sprintf("$%s", regs.GPR_NAMES[rs]))
			} else if machInst.IsRTwoOp() {
				buf.WriteString(fmt.Sprintf("$%s, $%s", regs.GPR_NAMES[rs], regs.GPR_NAMES[rt]))
			} else if machInst.IsRMt() {
				buf.WriteString(fmt.Sprintf("$%s", regs.GPR_NAMES[rs]))
			} else if machInst.IsRMf() {
				buf.WriteString(fmt.Sprintf("$%s", regs.GPR_NAMES[rd]))
			} else {
				buf.WriteString(fmt.Sprintf("$%s, $%s, $%s", regs.GPR_NAMES[rd], regs.GPR_NAMES[rs], regs.GPR_NAMES[rt]))
			}
		}
	default:
		panic("Impossible!")
	}

	return buf.String()
}
