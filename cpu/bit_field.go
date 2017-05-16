package cpu

type BitField struct {
	Hi uint32
	Lo uint32
}

func NewBitField(hi uint32, lo uint32) *BitField {
	var bitField = &BitField{
		Hi:hi,
		Lo:lo,
	}

	return bitField
}

var (
	OPCODE = NewBitField(31, 26)
	OPCODE_HI = NewBitField(31, 29)
	OPCODE_LO = NewBitField(28, 26)
	RS = NewBitField(25, 21)
	RT = NewBitField(20, 16)
	RD = NewBitField(15, 11)
	SHIFT = NewBitField(10, 6)
	FUNC = NewBitField(5, 0)
	FUNC_HI = NewBitField(5, 3)
	FUNC_LO = NewBitField(2, 0)
	COND = NewBitField(3, 0)
	UIMM = NewBitField(15, 0)
	TARGET = NewBitField(25, 0)
	FMT = NewBitField(25, 21)
	FMT3 = NewBitField(2, 0)
	FR = NewBitField(25, 21)
	FS = NewBitField(15, 11)
	FT = NewBitField(20, 16)
	FD = NewBitField(10, 6)
	BRANCH_CC = NewBitField(20, 18)
	CC = NewBitField(10, 8)
)
