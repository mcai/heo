package cpu

const (
	Mnemonic_NOP = "nop"
	Mnemonic_BREAK = "break"
	Mnemonic_SYSCALL = "syscall"
	Mnemonic_ADD = "add"
	Mnemonic_ADDI = "addi"
	Mnemonic_ADDIU = "addiu"
	Mnemonic_ADDU = "addu"
	Mnemonic_SUBU = "subu"
	Mnemonic_AND = "and"
	Mnemonic_ANDI = "andi"
	Mnemonic_NOR = "nor"
	Mnemonic_OR = "or"
	Mnemonic_ORI = "ori"
	Mnemonic_XOR = "xor"
	Mnemonic_XORI = "xori"
	Mnemonic_MULT = "mult"
	Mnemonic_MULTU = "multu"
	Mnemonic_DIV = "div"
	Mnemonic_DIVU = "divu"
	Mnemonic_SLL = "sll"
	Mnemonic_SLLV = "sllv"
	Mnemonic_SLT = "slt"
	Mnemonic_SLTI = "slti"
	Mnemonic_SLTIU = "sltiu"
	Mnemonic_SLTU = "sltu"
	Mnemonic_SRA = "sra"
	Mnemonic_SRAV = "srav"
	Mnemonic_SRL = "srl"
	Mnemonic_SRLV = "srlv"
	Mnemonic_MADD = "madd"
	Mnemonic_MSUB = "msub"
	Mnemonic_B = "b"
	Mnemonic_BAL = "bal"
	Mnemonic_BEQ = "beq"
	Mnemonic_BGEZ = "bgez"
	Mnemonic_BGEZAL = "bgezal"
	Mnemonic_BGTZ = "bgtz"
	Mnemonic_BLEZ = "blez"
	Mnemonic_BLTZ = "bltz"
	Mnemonic_BNE = "bne"
	Mnemonic_J = "j"
	Mnemonic_JAL = "jal"
	Mnemonic_JALR = "jalr"
	Mnemonic_JR = "jr"
	Mnemonic_LB = "lb"
	Mnemonic_LBU = "lbu"
	Mnemonic_LH = "lh"
	Mnemonic_LHU = "lhu"
	Mnemonic_LUI = "lui"
	Mnemonic_LW = "lw"
	Mnemonic_LWL = "lwl"
	Mnemonic_LWR = "lwr"
	Mnemonic_SB = "sb"
	Mnemonic_SH = "sh"
	Mnemonic_SW = "sw"
	Mnemonic_SWL = "swl"
	Mnemonic_SWR = "swr"
	Mnemonic_LDC1 = "ldc1"
	Mnemonic_LWC1 = "lwc1"
	Mnemonic_SDC1 = "sdc1"
	Mnemonic_SWC1 = "swc1"
	Mnemonic_MFHI = "mfhi"
	Mnemonic_MFLO = "mflo"
	Mnemonic_MTLO = "mtlo"
	Mnemonic_CFC1 = "cfc1"
	Mnemonic_CTC1 = "ctc1"
	Mnemonic_MFC1 = "mfc1"
	Mnemonic_MTC1 = "mtc1"
	Mnemonic_LL = "ll"
	Mnemonic_SC = "sc"
	Mnemonic_NEG_D = "neg_d"
	Mnemonic_MOV_D = "mov_d"
	Mnemonic_SQRT_D = "sqrt_d"
	Mnemonic_ABS_D = "abs_d"
	Mnemonic_MUL_D = "mul_d"
	Mnemonic_DIV_D = "div_d"
	Mnemonic_ADD_D = "add_d"
	Mnemonic_SUB_D = "sub_d"
	Mnemonic_MUL_S = "mul_s"
	Mnemonic_DIV_S = "div_s"
	Mnemonic_ADD_S = "add_s"
	Mnemonic_SUB_S = "sub_s"
	Mnemonic_MOV_S = "mov_s"
	Mnemonic_NEG_S = "neg_s"
	Mnemonic_ABS_S = "abs_s"
	Mnemonic_SQRT_S = "sqrt_s"
	Mnemonic_C_COND_D = "c_cond_d"
	Mnemonic_C_COND_S = "c_cond_s"
	Mnemonic_CVT_D_L = "cvt_d_l"
	Mnemonic_CVT_S_L = "cvt_s_l"
	Mnemonic_CVT_D_W = "cvt_d_w"
	Mnemonic_CVT_S_W = "cvt_s_w"
	Mnemonic_CVT_W_D = "cvt_w_d"
	Mnemonic_CVT_S_D = "cvt_s_d"
	Mnemonic_CVT_W_S = "cvt_w_s"
	Mnemonic_CVT_D_S = "cvt_d_s"
	Mnemonic_BC1F = "bc1f"
	Mnemonic_BC1T = "bc1t"
)

type MnemonicName string

type DecodeMethod struct {
	Bits uint32
	Mask uint32
}

func NewDecodeMethod(bits uint32, mask uint32) *DecodeMethod {
	var decodeMethod = &DecodeMethod{
		Bits:bits,
		Mask:mask,
	}

	return decodeMethod
}

type DecodeCondition struct {
	BitField *BitField
	Value    uint32
}

func NewDecodeCondition(bitField *BitField, value uint32) *DecodeCondition {
	var decodeCondition = &DecodeCondition{
		BitField:bitField,
		Value:value,
	}

	return decodeCondition
}

type Mnemonic struct {
	Name               MnemonicName

	DecodeMethod       *DecodeMethod
	DecodeCondition    *DecodeCondition

	Bits               uint32
	Mask               uint32

	ExtraBitField      *BitField
	ExtraBitFieldValue uint32

	FUOperationType    FUOperationType

	StaticInstType     StaticInstType
	StaticInstFlags    []StaticInstFlag

	InputDependencies  []StaticInstDependency
	OutputDependencies []StaticInstDependency

	Execute            func(context *Context, machInst MachInst)
}

func NewMnemonic(name MnemonicName, decodeMethod *DecodeMethod, decodeCondition *DecodeCondition, fuOperationType FUOperationType, staticInstType StaticInstType, staticInstFlags []StaticInstFlag, inputDependencies []StaticInstDependency, outputDependencies []StaticInstDependency, execute func(context *Context, machInst MachInst)) *Mnemonic {
	var mnemonic = &Mnemonic{
		Name:name,

		DecodeMethod:decodeMethod,
		DecodeCondition:decodeCondition,

		Bits:decodeMethod.Bits,
		Mask:decodeMethod.Mask,

		FUOperationType:fuOperationType,

		StaticInstType:staticInstType,
		StaticInstFlags:staticInstFlags,

		InputDependencies:inputDependencies,
		OutputDependencies:outputDependencies,

		Execute:execute,
	}

	if decodeCondition != nil {
		mnemonic.ExtraBitField = decodeCondition.BitField
		mnemonic.ExtraBitFieldValue = decodeCondition.Value
	}

	return mnemonic
}

func (mnemonic *Mnemonic) GetInputDependencies(machInst MachInst) []uint32 {
	var inputDependencies []uint32

	for _, staticInstDependency := range mnemonic.InputDependencies {
		inputDependencies = append(inputDependencies, RegisterDependencyToInt(staticInstDependency.ToRegisterDependency(machInst)))
	}

	return inputDependencies
}

func (mnemonic *Mnemonic) GetOutputDependencies(machInst MachInst) []uint32 {
	var outputDependencies []uint32

	for _, staticInstDependency := range mnemonic.OutputDependencies {
		outputDependencies = append(outputDependencies, RegisterDependencyToInt(staticInstDependency.ToRegisterDependency(machInst)))
	}

	return outputDependencies
}
