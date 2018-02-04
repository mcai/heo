package cpu

import "github.com/mcai/heo/cpu/regs"

type StaticInstFlag string

const (
	StaticInstFlag_INT_COMP             = StaticInstFlag("INT_COMP")
	StaticInstFlag_FP_COMP              = StaticInstFlag("FP_COMP")
	StaticInstFlag_UNCOND               = StaticInstFlag("UNCOND")
	StaticInstFlag_COND                 = StaticInstFlag("COND")
	StaticInstFlag_LD                   = StaticInstFlag("LD")
	StaticInstFlag_ST                   = StaticInstFlag("ST")
	StaticInstFlag_DIRECT_JMP           = StaticInstFlag("DIRECT_JUMP")
	StaticInstFlag_INDIRECT_JMP         = StaticInstFlag("INDIRECT_JUMP")
	StaticInstFlag_FUNC_CALL            = StaticInstFlag("FUNC_CALL")
	StaticInstFlag_FUNC_RET             = StaticInstFlag("FUNC_RET")
	StaticInstFlag_IMM                  = StaticInstFlag("IMM")
	StaticInstFlag_DISPLACED_ADDRESSING = StaticInstFlag("DISPLACED_ADDRESSING")
	StaticInstFlag_TRAP                 = StaticInstFlag("TRAP")
	StaticInstFlag_NOP                  = StaticInstFlag("NOP")
)

type StaticInstType string

const (
	StaticInstType_INT_COMP  = StaticInstType("INT_COMP")
	StaticInstType_FP_COMP   = StaticInstType("FP_COMP")
	StaticInstType_COND      = StaticInstType("COND")
	StaticInstType_UNCOND    = StaticInstType("UNCOND")
	StaticInstType_LD        = StaticInstType("LD")
	StaticInstType_ST        = StaticInstType("ST")
	StaticInstType_FUNC_CALL = StaticInstType("FUNC_CALL")
	StaticInstType_FUNC_RET  = StaticInstType("FUNC_RET")
	StaticInstType_TRAP      = StaticInstType("TRAP")
	StaticInstType_NOP       = StaticInstType("NOP")
)

func (staticInstType StaticInstType) IsControl() bool {
	return staticInstType == StaticInstType_COND ||
		staticInstType == StaticInstType_UNCOND ||
		staticInstType == StaticInstType_FUNC_CALL ||
		staticInstType == StaticInstType_FUNC_RET
}

func (staticInstType StaticInstType) IsLoadOrStore() bool {
	return staticInstType == StaticInstType_LD ||
		staticInstType == StaticInstType_ST
}

type RegisterDependencyType string

const (
	RegisterDependencyType_INT  = RegisterDependencyType("INT")
	RegisterDependencyType_FP   = RegisterDependencyType("FP")
	RegisterDependencyType_MISC = RegisterDependencyType("MISC")
)

func RegisterDependencyFromInt(i uint32) (RegisterDependencyType, uint32) {
	if i < regs.NUM_INT_REGISTERS {
		return RegisterDependencyType_INT, i
	} else if i < regs.NUM_INT_REGISTERS+regs.NUM_FP_REGISTERS {
		return RegisterDependencyType_FP, i - regs.NUM_INT_REGISTERS
	} else {
		return RegisterDependencyType_MISC, i - regs.NUM_INT_REGISTERS - regs.NUM_FP_REGISTERS
	}
}

func RegisterDependencyToInt(dependencyType RegisterDependencyType, num uint32) uint32 {
	switch dependencyType {
	case RegisterDependencyType_INT:
		return num
	case RegisterDependencyType_FP:
		return regs.NUM_INT_REGISTERS + num
	case RegisterDependencyType_MISC:
		return regs.NUM_INT_REGISTERS + regs.NUM_FP_REGISTERS + num
	default:
		panic("Impossible")
	}
}

type StaticInstDependency string

const (
	StaticInstDependency_RS            = StaticInstDependency("RS")
	StaticInstDependency_RT            = StaticInstDependency("RT")
	StaticInstDependency_RD            = StaticInstDependency("RD")
	StaticInstDependency_FS            = StaticInstDependency("FS")
	StaticInstDependency_FT            = StaticInstDependency("FT")
	StaticInstDependency_FD            = StaticInstDependency("FD")
	StaticInstDependency_REGISTER_RA   = StaticInstDependency("REGISTER_RA")
	StaticInstDependency_REGISTER_V0   = StaticInstDependency("REGISTER_V0")
	StaticInstDependency_REGISTER_HI   = StaticInstDependency("REGISTER_HI")
	StaticInstDependency_REGISTER_LO   = StaticInstDependency("REGISTER_LO")
	StaticInstDependency_REGISTER_FCSR = StaticInstDependency("REGISTER_FCSR")
)

func (staticInstDependency StaticInstDependency) ToRegisterDependency(machInst MachInst) (RegisterDependencyType, uint32) {
	switch staticInstDependency {
	case StaticInstDependency_RS:
		return RegisterDependencyType_INT, machInst.Rs()
	case StaticInstDependency_RT:
		return RegisterDependencyType_INT, machInst.Rt()
	case StaticInstDependency_RD:
		return RegisterDependencyType_INT, machInst.Rd()
	case StaticInstDependency_FS:
		return RegisterDependencyType_FP, machInst.Fs()
	case StaticInstDependency_FT:
		return RegisterDependencyType_FP, machInst.Ft()
	case StaticInstDependency_FD:
		return RegisterDependencyType_FP, machInst.Fd()
	case StaticInstDependency_REGISTER_RA:
		return RegisterDependencyType_INT, regs.REGISTER_RA
	case StaticInstDependency_REGISTER_V0:
		return RegisterDependencyType_INT, regs.REGISTER_V0
	case StaticInstDependency_REGISTER_HI:
		return RegisterDependencyType_MISC, regs.REGISTER_HI
	case StaticInstDependency_REGISTER_LO:
		return RegisterDependencyType_MISC, regs.REGISTER_LO
	case StaticInstDependency_REGISTER_FCSR:
		return RegisterDependencyType_FP, regs.REGISTER_FCSR
	default:
		panic("Impossible")
	}
}
