package cpu

type DynamicInst struct {
	Thread           Thread
	Pc               uint32
	StaticInst       *StaticInst
	EffectiveAddress int32
}

func NewDynamicInst(thread Thread, pc uint32, staticInst *StaticInst) *DynamicInst {
	var dynamicInst = &DynamicInst{
		Thread:thread,
		Pc:pc,
		StaticInst:staticInst,
	}

	if staticInst.Mnemonic.StaticInstType == StaticInstType_LD ||
		staticInst.Mnemonic.StaticInstType == StaticInstType_ST {
		dynamicInst.EffectiveAddress = int32(GetEffectiveAddress(thread.Context(), staticInst.MachInst))
	} else {
		dynamicInst.EffectiveAddress = -1
	}

	return dynamicInst
}


