package cpu

import (
	"github.com/mcai/heo/cpu/cpuutil"
	"github.com/mcai/heo/cpu/regs"
	"math"
)

const (
	FMT_SINGLE = 16
	FMT_DOUBLE = 17
	FMT_WORD   = 20
	FMT_LONG   = 21
)

type ISA struct {
	Mnemonics []*Mnemonic
}

func NewISA() *ISA {
	var isa = &ISA{
	}

	isa.addMnemonics()

	return isa
}

func (isa *ISA) ResetStats() {
}

func (isa *ISA) addMnemonic(name MnemonicName, decodeMethod *DecodeMethod, decodeCondition *DecodeCondition, fuOperationType FUOperationType, staticInstType StaticInstType, staticInstFlags []StaticInstFlag, inputDependencies []StaticInstDependency, outputDependencies []StaticInstDependency, execute func(context *Context, machInst MachInst)) {
	var mnemonic = NewMnemonic(name, decodeMethod, decodeCondition, fuOperationType, staticInstType, staticInstFlags, inputDependencies, outputDependencies, execute)

	isa.Mnemonics = append(isa.Mnemonics, mnemonic)
}

func (isa *ISA) addMnemonics() {
	isa.addMnemonic(
		Mnemonic_NOP,
		NewDecodeMethod(0x00000000, 0xffffffff),
		nil,
		FUOperationType_NONE,
		StaticInstType_NOP,
		[]StaticInstFlag{
			StaticInstFlag_NOP,
		},
		[]StaticInstDependency{},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
		},
	)

	isa.addMnemonic(
		Mnemonic_BC1F,
		NewDecodeMethod(0x45000000, 0xffe30000),
		nil,
		FUOperationType_NONE,
		StaticInstType_COND,
		[]StaticInstFlag{
			StaticInstFlag_COND,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_FCSR,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var branchCc = machInst.BranchCc()
			var imm = machInst.Imm()

			if fPCC(context, branchCc) == 0 {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_BC1T,
		NewDecodeMethod(0x45010000, 0xffe30000),
		nil,
		FUOperationType_NONE,
		StaticInstType_COND,
		[]StaticInstFlag{
			StaticInstFlag_COND,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_FCSR,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var branchCc = machInst.BranchCc()
			var imm = machInst.Imm()

			if fPCC(context, branchCc) != 0 {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_MFC1,
		NewDecodeMethod(0x44000000, 0xffe007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var rt = machInst.Rt()

			var temp = context.Regs().Fpr.Uint32(fs)
			context.Regs().Gpr[rt] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_MTC1,
		NewDecodeMethod(0x44800000, 0xffe007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()
			var fs = machInst.Fs()

			var temp = context.Regs().Gpr[rt]
			context.Regs().Fpr.SetUint32(fs, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_CFC1,
		NewDecodeMethod(0x44400000, 0xffe007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_FCSR,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var rt = machInst.Rt()

			if fs == 31 {
				var temp = context.Regs().Fcsr
				context.Regs().Gpr[rt] = temp
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_CTC1,
		NewDecodeMethod(0x44c00000, 0xffe007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_FCSR,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var rt = machInst.Rt()

			if fs != 0 {
				var temp = context.Regs().Gpr[rt]
				context.Regs().Fcsr = temp
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_ABS_S,
		NewDecodeMethod(0x44000005, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_CMP,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp1 = context.Regs().Fpr.Float32(fs)

			var temp2 float32

			if temp1 < 0.0 {
				temp2 = -temp1
			} else {
				temp2 = temp1
			}

			context.Regs().Fpr.SetFloat32(fd, temp2)
		},
	)

	isa.addMnemonic(
		Mnemonic_ABS_D,
		NewDecodeMethod(0x44000005, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_CMP,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp1 = context.Regs().Fpr.Float64(fs)

			var temp2 float64

			if temp1 < 0.0 {
				temp2 = -temp1
			} else {
				temp2 = temp1
			}

			context.Regs().Fpr.SetFloat64(fd, temp2)
		},
	)

	isa.addMnemonic(
		Mnemonic_ADD,
		NewDecodeMethod(0x00000020, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()

			var temp = context.Regs().Sgpr(rs) + context.Regs().Sgpr(rt)
			context.Regs().Gpr[machInst.Rd()] = uint32(temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_ADD_S,
		NewDecodeMethod(0x44000000, 0xfc00003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_ADD,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float32(fs) + context.Regs().Fpr.Float32(ft)
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_ADD_D,
		NewDecodeMethod(0x44000000, 0xfc00003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_ADD,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float64(fs) + context.Regs().Fpr.Float64(ft)
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_ADDI,
		NewDecodeMethod(0x20000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
			StaticInstFlag_IMM,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()
			var rt = machInst.Rt()

			var temp = context.Regs().Sgpr(rs) + imm
			context.Regs().Gpr[rt] = uint32(temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_ADDIU,
		NewDecodeMethod(0x24000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
			StaticInstFlag_IMM,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()
			var rt = machInst.Rt()

			var temp = context.Regs().Sgpr(rs) + imm
			context.Regs().Gpr[rt] = uint32(temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_ADDU,
		NewDecodeMethod(0x00000021, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			var temp = context.Regs().Gpr[rs] + context.Regs().Gpr[rt]
			context.Regs().Gpr[rd] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_AND,
		NewDecodeMethod(0x00000024, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			var temp = context.Regs().Gpr[rs] & context.Regs().Gpr[rt]
			context.Regs().Gpr[rd] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_ANDI,
		NewDecodeMethod(0x30000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
			StaticInstFlag_IMM,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var uimm = machInst.Uimm()
			var rt = machInst.Rt()

			context.Regs().Gpr[rt] = context.Regs().Gpr[rs] & uimm
		},
	)

	isa.addMnemonic(
		Mnemonic_B,
		NewDecodeMethod(0x10000000, 0xffff0000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_UNCOND,
		[]StaticInstFlag{
			StaticInstFlag_UNCOND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var imm = machInst.Imm()

			relBranch(context, imm<<2)
		},
	)

	isa.addMnemonic(
		Mnemonic_BAL,
		NewDecodeMethod(0x04110000, 0xffff0000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_UNCOND,
		[]StaticInstFlag{
			StaticInstFlag_UNCOND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_RA,
		},
		func(context *Context, machInst MachInst) {
			var imm = machInst.Imm()

			context.Regs().Gpr[regs.REGISTER_RA] = context.Regs().Pc + 8
			relBranch(context, imm<<2)
		},
	)

	isa.addMnemonic(
		Mnemonic_BEQ,
		NewDecodeMethod(0x10000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_COND,
		[]StaticInstFlag{
			StaticInstFlag_COND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var imm = machInst.Imm()

			if context.Regs().Gpr[rs] == context.Regs().Gpr[rt] {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_BGEZ,
		NewDecodeMethod(0x04010000, 0xfc1f0000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_COND,
		[]StaticInstFlag{
			StaticInstFlag_COND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()

			if context.Regs().Sgpr(rs) >= 0 {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_BGEZAL,
		NewDecodeMethod(0x04110000, 0xfc1f0000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_FUNC_CALL,
		[]StaticInstFlag{
			StaticInstFlag_COND,
			StaticInstFlag_DIRECT_JMP,
			StaticInstFlag_FUNC_CALL,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_RA,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()

			context.Regs().Gpr[regs.REGISTER_RA] = context.Regs().Pc + 8
			if context.Regs().Sgpr(rs) >= 0 {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_BGTZ,
		NewDecodeMethod(0x1c000000, 0xfc1f0000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_COND,
		[]StaticInstFlag{
			StaticInstFlag_COND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()

			if context.Regs().Sgpr(rs) > 0 {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_BLEZ,
		NewDecodeMethod(0x18000000, 0xfc1f0000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_COND,
		[]StaticInstFlag{
			StaticInstFlag_COND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()

			if context.Regs().Sgpr(rs) <= 0 {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_BLTZ,
		NewDecodeMethod(0x04000000, 0xfc1f0000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_COND,
		[]StaticInstFlag{
			StaticInstFlag_COND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()

			if context.Regs().Sgpr(rs) < 0 {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_BNE,
		NewDecodeMethod(0x14000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_COND,
		[]StaticInstFlag{
			StaticInstFlag_COND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var imm = machInst.Imm()

			if context.Regs().Gpr[rs] != context.Regs().Gpr[rt] {
				relBranch(context, imm<<2)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_BREAK,
		NewDecodeMethod(0x0000000d, 0xfc00003f),
		nil,
		FUOperationType_NONE,
		StaticInstType_TRAP,
		[]StaticInstFlag{
			StaticInstFlag_TRAP,
		},
		[]StaticInstDependency{},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			if !context.Speculative {
				context.Finish()
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_C_COND_D,
		NewDecodeMethod(0x44000030, 0xfc0000f0),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_CMP,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
			StaticInstDependency_REGISTER_FCSR,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_FCSR,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()

			var op1 = context.Regs().Fpr.Float64(fs)
			var op2 = context.Regs().Fpr.Float64(ft)

			var less = op1 < op2
			var equal = op1 == op2

			var unordered = false

			cCond(context, machInst, less, equal, unordered)
		},
	)

	isa.addMnemonic(
		Mnemonic_C_COND_S,
		NewDecodeMethod(0x44000030, 0xfc0000f0),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_CMP,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
			StaticInstDependency_REGISTER_FCSR,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_FCSR,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()

			var op1 = context.Regs().Fpr.Float32(fs)
			var op2 = context.Regs().Fpr.Float32(ft)

			var less = op1 < op2
			var equal = op1 == op2

			var unordered = false

			cCond(context, machInst, less, equal, unordered)
		},
	)

	isa.addMnemonic(
		Mnemonic_CVT_D_S,
		NewDecodeMethod(0x44000021, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_CVT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = float64(context.Regs().Fpr.Float32(fs))
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_CVT_D_W,
		NewDecodeMethod(0x44000021, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_WORD),
		FUOperationType_FP_CVT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = float64(uint64(context.Regs().Fpr.Uint32(fs)))
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_CVT_D_L,
		NewDecodeMethod(0x44000021, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_LONG),
		FUOperationType_FP_CVT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = float64(context.Regs().Fpr.Uint64(fs))
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_CVT_S_D,
		NewDecodeMethod(0x44000020, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_CVT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()
			var temp = float32(context.Regs().Fpr.Float64(fs))
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_CVT_S_W,
		NewDecodeMethod(0x44000020, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_WORD),
		FUOperationType_FP_CVT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = float32(context.Regs().Fpr.Uint32(fs))
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_CVT_S_L,
		NewDecodeMethod(0x44000020, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_LONG),
		FUOperationType_FP_CVT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = float32(context.Regs().Fpr.Uint64(fs))
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_CVT_W_S,
		NewDecodeMethod(0x44000024, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_CVT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = uint32(context.Regs().Fpr.Float32(fs))
			context.Regs().Fpr.SetUint32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_CVT_W_D,
		NewDecodeMethod(0x44000024, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_CVT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = uint32(context.Regs().Fpr.Float64(fs))
			context.Regs().Fpr.SetUint32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_DIV,
		NewDecodeMethod(0x0000001a, 0xfc00ffff),
		nil,
		FUOperationType_INT_DIV,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_HI,
			StaticInstDependency_REGISTER_LO,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()

			if context.Regs().Gpr[rt] != 0 {
				context.Regs().Lo = uint32(context.Regs().Sgpr(rs) /
					context.Regs().Sgpr(rt))
				context.Regs().Hi = uint32(context.Regs().Sgpr(rs) %
					context.Regs().Sgpr(rt))
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_DIV_S,
		NewDecodeMethod(0x44000003, 0xfc00003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_DIV,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float32(fs) /
				context.Regs().Fpr.Float32(ft)
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_DIV_D,
		NewDecodeMethod(0x44000003, 0xfc00003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_DIV,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float64(fs) /
				context.Regs().Fpr.Float64(ft)
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_DIVU,
		NewDecodeMethod(0x0000001b, 0xfc00003f),
		nil,
		FUOperationType_INT_DIV,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_HI,
			StaticInstDependency_REGISTER_LO,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()

			if context.Regs().Gpr[rt] != 0 {
				context.Regs().Lo = context.Regs().Gpr[rs] /
					context.Regs().Gpr[rt]
				context.Regs().Hi = context.Regs().Gpr[rs] %
					context.Regs().Gpr[rt]
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_J,
		NewDecodeMethod(0x08000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_UNCOND,
		[]StaticInstFlag{
			StaticInstFlag_UNCOND,
			StaticInstFlag_DIRECT_JMP,
		},
		[]StaticInstDependency{},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var target = machInst.Target()

			var dest = (cpuutil.GetBits32(context.Regs().Pc+4, 32, 28) << 28) | (target << 2)
			branch(context, dest)
		},
	)

	isa.addMnemonic(
		Mnemonic_JAL,
		NewDecodeMethod(0x0c000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_FUNC_CALL,
		[]StaticInstFlag{
			StaticInstFlag_UNCOND,
			StaticInstFlag_DIRECT_JMP,
			StaticInstFlag_FUNC_CALL,
		},
		[]StaticInstDependency{},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_RA,
		},
		func(context *Context, machInst MachInst) {
			var target = machInst.Target()

			var dest = (cpuutil.GetBits32(context.Regs().Pc+4, 32, 28) << 28) | (target << 2)
			context.Regs().Gpr[regs.REGISTER_RA] = context.Regs().Pc + 8
			branch(context, dest)
		},
	)

	isa.addMnemonic(
		Mnemonic_JALR,
		NewDecodeMethod(0x00000009, 0xfc00003f),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_FUNC_CALL,
		[]StaticInstFlag{
			StaticInstFlag_UNCOND,
			StaticInstFlag_INDIRECT_JMP,
			StaticInstFlag_FUNC_CALL,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rd = machInst.Rd()

			branch(context, context.Regs().Gpr[rs])
			context.Regs().Gpr[rd] = context.Regs().Pc + 8
		},
	)

	isa.addMnemonic(
		Mnemonic_JR,
		NewDecodeMethod(0x00000008, 0xfc00003f),
		nil,
		FUOperationType_NONE,
		StaticInstType_FUNC_RET,
		[]StaticInstFlag{
			StaticInstFlag_UNCOND,
			StaticInstFlag_INDIRECT_JMP,
			StaticInstFlag_FUNC_RET,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()

			branch(context, context.Regs().Gpr[rs])
		},
	)

	isa.addMnemonic(
		Mnemonic_LB,
		NewDecodeMethod(0x80000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)
			var temp = context.Process.Memory().ReadUInt8At(addr)
			context.Regs().Gpr[rt] = cpuutil.SignExtend32(uint32(temp), 8)
		},
	)

	isa.addMnemonic(
		Mnemonic_LBU,
		NewDecodeMethod(0x90000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)
			var temp = context.Process.Memory().ReadUInt8At(addr)
			context.Regs().Gpr[rt] = uint32(temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_LDC1,
		NewDecodeMethod(0xd4000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FT,
		},
		func(context *Context, machInst MachInst) {
			var ft = machInst.Ft()

			var addr = GetEffectiveAddress(context, machInst)
			var temp = context.Process.Memory().ReadUInt64At(addr)
			context.Regs().Fpr.SetFloat64(ft, math.Float64frombits(temp))
		},
	)

	isa.addMnemonic(
		Mnemonic_LH,
		NewDecodeMethod(0x84000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)
			var temp = context.Process.Memory().ReadUInt16At(addr)
			context.Regs().Gpr[rt] = cpuutil.SignExtend32(uint32(temp), 16)
		},
	)

	isa.addMnemonic(
		Mnemonic_LHU,
		NewDecodeMethod(0x94000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)
			var temp = context.Process.Memory().ReadUInt16At(addr)
			context.Regs().Gpr[rt] = uint32(temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_LL,
		NewDecodeMethod(0xc0000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)
			var temp = context.Process.Memory().ReadUInt32At(addr)
			context.Regs().Gpr[rt] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_LUI,
		NewDecodeMethod(0x3c000000, 0xffe00000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()
			var uimm = machInst.Uimm()

			context.Regs().Gpr[rt] = uimm << 16
		},
	)

	isa.addMnemonic(
		Mnemonic_LW,
		NewDecodeMethod(0x8c000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)
			var temp = context.Process.Memory().ReadUInt32At(addr)
			context.Regs().Gpr[rt] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_LWC1,
		NewDecodeMethod(0xc4000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FT,
		},
		func(context *Context, machInst MachInst) {
			var ft = machInst.Ft()

			var addr = GetEffectiveAddress(context, machInst)
			var temp = context.Process.Memory().ReadUInt32At(addr)
			context.Regs().Fpr.SetFloat32(ft, math.Float32frombits(temp))
		},
	)

	isa.addMnemonic(
		Mnemonic_LWL,
		NewDecodeMethod(0x88000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()
			var addr = GetEffectiveAddress(context, machInst)

			var size = 4 - (addr & 3)

			var src = context.Process.Memory().ReadBlockAt(addr, size)

			var dst = make([]byte, 4)

			context.Process.Memory().ByteOrder.PutUint32(dst, context.Regs().Gpr[rt])

			for i := uint32(0); i < size; i++ {
				dst[3-i] = src[i]
			}

			var temp = context.Process.Memory().ByteOrder.Uint32(dst)

			context.Regs().Gpr[rt] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_LWR,
		NewDecodeMethod(0x98000000, 0xfc000000),
		nil,
		FUOperationType_READ_PORT,
		StaticInstType_LD,
		[]StaticInstFlag{
			StaticInstFlag_LD,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)

			var size = 1 + (addr & 3)

			var src = context.Process.Memory().ReadBlockAt(addr-size+1, size)

			var dst = make([]byte, 4)

			context.Process.Memory().ByteOrder.PutUint32(dst, context.Regs().Gpr[rt])

			for i := uint32(0); i < size; i++ {
				dst[size-i-1] = src[i]
			}

			var temp = context.Process.Memory().ByteOrder.Uint32(dst)

			context.Regs().Gpr[rt] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_MADD,
		NewDecodeMethod(0x70000000, 0xfc00ffff),
		nil,
		FUOperationType_INT_MULT,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
			StaticInstDependency_REGISTER_HI,
			StaticInstDependency_REGISTER_LO,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_HI,
			StaticInstDependency_REGISTER_LO,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()

			var temp1 = int64(context.Regs().Sgpr(rs))
			var temp2 = int64(context.Regs().Sgpr(rt))
			var temp3 = int64(context.Regs().Hi<<32) | int64(context.Regs().Lo)
			var temp = temp1*temp2 + temp3
			context.Regs().Hi = uint32(cpuutil.GetBits64(uint64(temp), 63, 32))
			context.Regs().Lo = uint32(cpuutil.GetBits64(uint64(temp), 31, 0))
		},
	)

	isa.addMnemonic(
		Mnemonic_MFHI,
		NewDecodeMethod(0x00000010, 0xffff07ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_HI,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rd = machInst.Rd()

			context.Regs().Gpr[rd] = context.Regs().Hi
		},
	)

	isa.addMnemonic(
		Mnemonic_MFLO,
		NewDecodeMethod(0x00000012, 0xffff07ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_LO,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rd = machInst.Rd()

			context.Regs().Gpr[rd] = context.Regs().Lo
		},
	)

	isa.addMnemonic(
		Mnemonic_MOV_S,
		NewDecodeMethod(0x44000006, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_NONE,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float32(fs)
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_MOV_D,
		NewDecodeMethod(0x44000006, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_NONE,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float64(fs)
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_MSUB,
		NewDecodeMethod(0x70000004, 0xfc00ffff),
		nil,
		FUOperationType_INT_MULT,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
			StaticInstDependency_REGISTER_HI,
			StaticInstDependency_REGISTER_LO,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_HI,
			StaticInstDependency_REGISTER_LO,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()

			var temp1 = int64(context.Regs().Sgpr(rs))
			var temp2 = int64(context.Regs().Sgpr(rt))
			var temp3 = int64(context.Regs().Hi<<32) | int64(context.Regs().Lo)
			var temp = temp3 - temp1*temp2 + temp3
			context.Regs().Hi = uint32(cpuutil.GetBits64(uint64(temp), 63, 32))
			context.Regs().Lo = uint32(cpuutil.GetBits64(uint64(temp), 31, 0))
		},
	)

	isa.addMnemonic(
		Mnemonic_MTLO,
		NewDecodeMethod(0x00000013, 0xfc1fffff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_LO,
		},
		func(context *Context, machInst MachInst) {
			var rd = machInst.Rd()

			context.Regs().Lo = context.Regs().Gpr[rd]
		},
	)

	isa.addMnemonic(
		Mnemonic_MUL_S,
		NewDecodeMethod(0x44000002, 0xfc00003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_MULT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float32(fs) *
				context.Regs().Fpr.Float32(ft)
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_MUL_D,
		NewDecodeMethod(0x44000002, 0xfc00003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_MULT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float64(fs) *
				context.Regs().Fpr.Float64(ft)
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_MULT,
		NewDecodeMethod(0x00000018, 0xfc00003f),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_HI,
			StaticInstDependency_REGISTER_LO,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()

			var temp = uint64(int64(context.Regs().Sgpr(rs)) *
				int64(context.Regs().Sgpr(rt)))
			context.Regs().Lo = uint32(cpuutil.GetBits64(temp, 31, 0))
			context.Regs().Hi = uint32(cpuutil.GetBits64(temp, 63, 32))
		},
	)

	isa.addMnemonic(
		Mnemonic_MULTU,
		NewDecodeMethod(0x00000019, 0xfc00003f),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_HI,
			StaticInstDependency_REGISTER_LO,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()

			var temp = uint64(context.Regs().Gpr[rs]) *
				uint64(context.Regs().Gpr[rt])
			context.Regs().Lo = uint32(cpuutil.GetBits64(temp, 31, 0))
			context.Regs().Hi = uint32(cpuutil.GetBits64(temp, 63, 32))
		},
	)

	isa.addMnemonic(
		Mnemonic_NEG_S,
		NewDecodeMethod(0x44000007, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_CMP,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = -context.Regs().Fpr.Float32(fs)
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_NEG_D,
		NewDecodeMethod(0x44000007, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_CMP,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = -context.Regs().Fpr.Float64(fs)
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_NOR,
		NewDecodeMethod(0x00000027, 0xfc00003f),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			var temp = context.Regs().Gpr[rs] | context.Regs().Gpr[rt]
			context.Regs().Gpr[rd] = ^temp
		},
	)

	isa.addMnemonic(
		Mnemonic_OR,
		NewDecodeMethod(0x00000025, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			var temp = context.Regs().Gpr[rs] | context.Regs().Gpr[rt]
			context.Regs().Gpr[rd] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_ORI,
		NewDecodeMethod(0x34000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
			StaticInstFlag_IMM,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var uimm = machInst.Uimm()
			var rt = machInst.Rt()

			var temp = context.Regs().Gpr[rs] | uimm
			context.Regs().Gpr[rt] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_SB,
		NewDecodeMethod(0xa0000000, 0xfc000000),
		nil,
		FUOperationType_WRITE_PORT,
		StaticInstType_ST,
		[]StaticInstFlag{
			StaticInstFlag_ST,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var temp = byte(context.Regs().Gpr[rt])
			var addr = GetEffectiveAddress(context, machInst)
			context.Process.Memory().WriteUInt8At(addr, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_SC,
		NewDecodeMethod(0xe0000000, 0xfc000000),
		nil,
		FUOperationType_WRITE_PORT,
		StaticInstType_ST,
		[]StaticInstFlag{
			StaticInstFlag_ST,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var temp = context.Regs().Gpr[rt]
			var addr = GetEffectiveAddress(context, machInst)
			context.Process.Memory().WriteUInt32At(addr, temp)
			context.Regs().Gpr[rt] = 1
		},
	)

	isa.addMnemonic(
		Mnemonic_SDC1,
		NewDecodeMethod(0xf4000000, 0xfc000000),
		nil,
		FUOperationType_WRITE_PORT,
		StaticInstType_ST,
		[]StaticInstFlag{
			StaticInstFlag_ST,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var ft = machInst.Ft()

			var dbl = context.Regs().Fpr.Float64(ft)
			var temp = math.Float64bits(dbl)
			var addr = GetEffectiveAddress(context, machInst)
			context.Process.Memory().WriteUInt64At(addr, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_SH,
		NewDecodeMethod(0xa4000000, 0xfc000000),
		nil,
		FUOperationType_WRITE_PORT,
		StaticInstType_ST,
		[]StaticInstFlag{
			StaticInstFlag_ST,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var temp = uint16(context.Regs().Gpr[rt])
			var addr = GetEffectiveAddress(context, machInst)
			context.Process.Memory().WriteUInt16At(addr, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_SLL,
		NewDecodeMethod(0x00000000, 0xffe0003f),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()
			var shift = machInst.Shift()
			var rd = machInst.Rd()

			var temp = context.Regs().Gpr[rt] << shift
			context.Regs().Gpr[rd] = temp
		},
	)

	isa.addMnemonic(
		Mnemonic_SLLV,
		NewDecodeMethod(0x00000004, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			var s = cpuutil.GetBits32(context.Regs().Gpr[rs], 4, 0)
			context.Regs().Gpr[rd] = context.Regs().Gpr[rt] << s
		},
	)

	isa.addMnemonic(
		Mnemonic_SLT,
		NewDecodeMethod(0x0000002a, 0xfc00003f),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			if context.Regs().Sgpr(rs) < context.Regs().Sgpr(rt) {
				context.Regs().Gpr[rd] = 1
			} else {
				context.Regs().Gpr[rd] = 0
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_SLTI,
		NewDecodeMethod(0x28000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
			StaticInstFlag_IMM,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()
			var rt = machInst.Rt()

			if context.Regs().Sgpr(rs) < imm {
				context.Regs().Gpr[rt] = 1
			} else {
				context.Regs().Gpr[rt] = 0
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_SLTIU,
		NewDecodeMethod(0x2c000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
			StaticInstFlag_IMM,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var imm = machInst.Imm()
			var rt = machInst.Rt()

			if context.Regs().Gpr[rs] < uint32(imm) {
				context.Regs().Gpr[rt] = 1
			} else {
				context.Regs().Gpr[rt] = 0
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_SLTU,
		NewDecodeMethod(0x0000002b, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			if context.Regs().Gpr[rs] < context.Regs().Gpr[rt] {
				context.Regs().Gpr[rd] = 1
			} else {
				context.Regs().Gpr[rd] = 0
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_SQRT_S,
		NewDecodeMethod(0x44000004, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_SQRT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = float32(math.Sqrt(float64(context.Regs().Fpr.Float32(fs))))
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_SQRT_D,
		NewDecodeMethod(0x44000004, 0xfc1f003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_SQRT,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var fd = machInst.Fd()

			var temp = math.Sqrt(context.Regs().Fpr.Float64(fs))
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_SRA,
		NewDecodeMethod(0x00000003, 0xffe0003f),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()
			var shift = machInst.Shift()
			var rd = machInst.Rd()

			context.Regs().Gpr[rd] = uint32(context.Regs().Sgpr(rt) >> shift)
		},
	)

	isa.addMnemonic(
		Mnemonic_SRAV,
		NewDecodeMethod(0x00000007, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			var s = int32(cpuutil.GetBits32(context.Regs().Gpr[rs], 4, 0))
			context.Regs().Gpr[rd] = uint32(context.Regs().Sgpr(rt) >> uint32(s))
		},
	)

	isa.addMnemonic(
		Mnemonic_SRL,
		NewDecodeMethod(0x00000002, 0xffe0003f),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()
			var shift = machInst.Shift()
			var rd = machInst.Rd()

			context.Regs().Gpr[rd] = context.Regs().Gpr[rt] >> shift
		},
	)

	isa.addMnemonic(
		Mnemonic_SRLV,
		NewDecodeMethod(0x00000006, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			var s = cpuutil.GetBits32(context.Regs().Gpr[rs], 4, 0)
			context.Regs().Gpr[rd] = context.Regs().Gpr[rt] >> s
		},
	)

	isa.addMnemonic(
		Mnemonic_SUB_S,
		NewDecodeMethod(0x44000001, 0xfc00003f),
		NewDecodeCondition(FMT, FMT_SINGLE),
		FUOperationType_FP_ADD,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float32(fs) -
				context.Regs().Fpr.Float32(ft)
			context.Regs().Fpr.SetFloat32(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_SUB_D,
		NewDecodeMethod(0x44000001, 0xfc00003f),
		NewDecodeCondition(FMT, FMT_DOUBLE),
		FUOperationType_FP_ADD,
		StaticInstType_FP_COMP,
		[]StaticInstFlag{
			StaticInstFlag_FP_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_FS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{
			StaticInstDependency_FD,
		},
		func(context *Context, machInst MachInst) {
			var fs = machInst.Fs()
			var ft = machInst.Ft()
			var fd = machInst.Fd()

			var temp = context.Regs().Fpr.Float64(fs) -
				context.Regs().Fpr.Float64(ft)
			context.Regs().Fpr.SetFloat64(fd, temp)
		},
	)

	isa.addMnemonic(
		Mnemonic_SUBU,
		NewDecodeMethod(0x00000023, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			context.Regs().Gpr[rd] = context.Regs().Gpr[rs] -
				context.Regs().Gpr[rt]
		},
	)

	isa.addMnemonic(
		Mnemonic_SW,
		NewDecodeMethod(0xac000000, 0xfc000000),
		nil,
		FUOperationType_WRITE_PORT,
		StaticInstType_ST,
		[]StaticInstFlag{
			StaticInstFlag_ST,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var data = context.Regs().Gpr[rt]
			var addr = GetEffectiveAddress(context, machInst)
			context.Process.Memory().WriteUInt32At(addr, data)
		},
	)

	isa.addMnemonic(
		Mnemonic_SWC1,
		NewDecodeMethod(0xe4000000, 0xfc000000),
		nil,
		FUOperationType_WRITE_PORT,
		StaticInstType_ST,
		[]StaticInstFlag{
			StaticInstFlag_ST,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_FT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var ft = machInst.Ft()

			var f = context.Regs().Fpr.Float32(ft)
			var data = math.Float32bits(f)
			var addr = GetEffectiveAddress(context, machInst)
			context.Process.Memory().WriteUInt32At(addr, data)
		},
	)

	isa.addMnemonic(
		Mnemonic_SWL,
		NewDecodeMethod(0xa8000000, 0xfc000000),
		nil,
		FUOperationType_WRITE_PORT,
		StaticInstType_ST,
		[]StaticInstFlag{
			StaticInstFlag_ST,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)

			var size = 4 - (addr & 3)

			var src = make([]byte, 4)

			context.Process.Memory().ByteOrder.PutUint32(src, context.Regs().Gpr[rt])

			var dst = make([]byte, 4)

			for i := uint32(0); i < size; i++ {
				dst[i] = src[3-i]
			}

			context.Process.Memory().WriteBlockAt(addr, size, dst)
		},
	)

	isa.addMnemonic(
		Mnemonic_SWR,
		NewDecodeMethod(0xb8000000, 0xfc000000),
		nil,
		FUOperationType_WRITE_PORT,
		StaticInstType_ST,
		[]StaticInstFlag{
			StaticInstFlag_ST,
			StaticInstFlag_DISPLACED_ADDRESSING,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			var rt = machInst.Rt()

			var addr = GetEffectiveAddress(context, machInst)

			var size = 1 + (addr & 3)

			var src = make([]byte, 4)

			context.Process.Memory().ByteOrder.PutUint32(src, context.Regs().Gpr[rt])

			var dst = make([]byte, 4)

			for i := uint32(0); i < size; i++ {
				dst[i] = src[size-i-1]
			}

			context.Process.Memory().WriteBlockAt(addr-size+1, size, dst)
		},
	)

	isa.addMnemonic(
		Mnemonic_SYSCALL,
		NewDecodeMethod(0x0000000c, 0xfc00003f),
		nil,
		FUOperationType_NONE,
		StaticInstType_TRAP,
		[]StaticInstFlag{
			StaticInstFlag_TRAP,
		},
		[]StaticInstDependency{
			StaticInstDependency_REGISTER_V0,
		},
		[]StaticInstDependency{},
		func(context *Context, machInst MachInst) {
			if !context.Speculative {
				context.Kernel.SyscallEmulation.DoSyscall(context.Regs().Gpr[regs.REGISTER_V0], context)
			}
		},
	)

	isa.addMnemonic(
		Mnemonic_XOR,
		NewDecodeMethod(0x00000026, 0xfc0007ff),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
			StaticInstDependency_RT,
		},
		[]StaticInstDependency{
			StaticInstDependency_RD,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var rt = machInst.Rt()
			var rd = machInst.Rd()

			context.Regs().Gpr[rd] = context.Regs().Gpr[rs] ^
				context.Regs().Gpr[rt]
		},
	)

	isa.addMnemonic(
		Mnemonic_XORI,
		NewDecodeMethod(0x38000000, 0xfc000000),
		nil,
		FUOperationType_INT_ALU,
		StaticInstType_INT_COMP,
		[]StaticInstFlag{
			StaticInstFlag_INT_COMP,
			StaticInstFlag_IMM,
		},
		[]StaticInstDependency{
			StaticInstDependency_RS,
		},
		[]StaticInstDependency{
			StaticInstDependency_RT,
		},
		func(context *Context, machInst MachInst) {
			var rs = machInst.Rs()
			var uimm = machInst.Uimm()
			var rt = machInst.Rt()

			context.Regs().Gpr[rt] = context.Regs().Gpr[rs] ^ uimm
		},
	)
}

func cCond(context *Context, machInst MachInst, less bool, equal bool, unordered bool) {
	var cc = machInst.Cc()

	var condition = (cpuutil.GetBit32(machInst.Cond(), 2) != 0 && less) ||
		(cpuutil.GetBit32(machInst.Cond(), 1) != 0 && equal) ||
		(cpuutil.GetBit32(machInst.Cond(), 0) != 0 && unordered)

	setFPCC(context, cc, condition)
}

func branch(context *Context, v uint32) {
	context.Regs().Nnpc = v
}

func relBranch(context *Context, v int32) {
	context.Regs().Nnpc = uint32(int32(context.Regs().Pc) + v + 4)
}

func fPCC(context *Context, c uint32) uint32 {
	if c != 0 {
		return cpuutil.GetBit32(context.Regs().Fcsr, 24+c)
	} else {
		return cpuutil.GetBit32(context.Regs().Fcsr, 23)
	}
}

func setFPCC(context *Context, c uint32, v bool) {
	if c != 0 {
		context.Regs().Fcsr = cpuutil.SetBitValue32(context.Regs().Fcsr, 24+c, v)
	} else {
		context.Regs().Fcsr = cpuutil.SetBitValue32(context.Regs().Fcsr, 23, v)
	}
}

func GetEffectiveAddress(context *Context, machInst MachInst) uint32 {
	return uint32(int32(context.Regs().Gpr[machInst.Rs()]) + machInst.Imm())
}
