package cpu

import "github.com/mcai/heo/simutil"

type TwoBitBranchPredictorUpdate struct {
	SaturatingCounter *simutil.SaturatingCounter
	Ras               bool
}

func NewTwoBitBranchPredictorUpdate() *TwoBitBranchPredictorUpdate {
	var branchPredictorUpdate = &TwoBitBranchPredictorUpdate{
	}

	return branchPredictorUpdate
}

type TwoBitBranchPredictor struct {
	*BaseBranchPredictor

	branchTargetBuffer *BranchTargetBuffer
	returnAddressStack *ReturnAddressStack

	size               uint32
	saturatingCounters []*simutil.SaturatingCounter
}

func NewTwoBitBranchPredictor(thread Thread, branchTargetBufferNumSets uint32, branchTargetBufferAssoc uint32, returnAddressStackSize uint32, size uint32) *TwoBitBranchPredictor {
	var branchPredictor = &TwoBitBranchPredictor{
		BaseBranchPredictor:NewBaseBranchPredictor(thread),

		branchTargetBuffer:NewBranchTargetBuffer(branchTargetBufferNumSets, branchTargetBufferAssoc),
		returnAddressStack:NewReturnAddressStack(returnAddressStackSize),

		size:size,
	}

	var flipFlop = uint32(1)

	for i := uint32(0); i < size; i++ {
		branchPredictor.saturatingCounters = append(
			branchPredictor.saturatingCounters,
			simutil.NewSaturatingCounter(0, 2, 3, flipFlop),
		)

		flipFlop = 3 - flipFlop
	}

	return branchPredictor
}

func (branchPredictor *TwoBitBranchPredictor) getSaturatingCounter(branchAddress uint32) *simutil.SaturatingCounter {
	var index = (branchAddress >> BRANCH_SHIFT) & (branchPredictor.size - 1)

	return branchPredictor.saturatingCounters[index]
}

func (branchPredictor *TwoBitBranchPredictor) Predict(branchAddress uint32, mnemonic *Mnemonic) (uint32, uint32, interface{}) {
	var branchPredictorUpdate = NewTwoBitBranchPredictorUpdate()

	if mnemonic.StaticInstType == StaticInstType_COND {
		branchPredictorUpdate.SaturatingCounter = branchPredictor.getSaturatingCounter(branchAddress)
	}

	var returnAddressStackRecoverTop = branchPredictor.returnAddressStack.Top()

	if mnemonic.StaticInstType == StaticInstType_FUNC_RET && branchPredictor.returnAddressStack.Size() > 0 {
		branchPredictorUpdate.Ras = true
		return branchPredictor.returnAddressStack.Pop(), returnAddressStackRecoverTop, branchPredictorUpdate
	}

	if mnemonic.StaticInstType == StaticInstType_FUNC_CALL && branchPredictor.returnAddressStack.Size() > 0 {
		branchPredictor.returnAddressStack.Push(branchAddress)
	}

	if mnemonic.StaticInstType != StaticInstType_COND || branchPredictorUpdate.SaturatingCounter.Taken() {
		var branchTargetBufferEntry = branchPredictor.branchTargetBuffer.Lookup(branchAddress)

		if branchTargetBufferEntry != nil {
			return branchTargetBufferEntry.Target, returnAddressStackRecoverTop, branchPredictorUpdate
		} else {
			return 0, returnAddressStackRecoverTop, branchPredictorUpdate
		}
	} else {
		return 0, returnAddressStackRecoverTop, branchPredictorUpdate
	}
}

func (branchPredictor *TwoBitBranchPredictor) Update(branchAddress uint32, branchTarget uint32, taken bool, correct bool, mnemonic *Mnemonic, branchPredictorUpdate interface{}) {
	branchPredictor.BaseBranchPredictor.Update(branchAddress, branchTarget, taken, correct, mnemonic, branchPredictorUpdate)

	var twoBitBranchPredictorUpdate = branchPredictorUpdate.(*TwoBitBranchPredictorUpdate)

	if mnemonic.StaticInstType == StaticInstType_FUNC_RET {
		if !twoBitBranchPredictorUpdate.Ras {
			return
		}
	}

	if mnemonic.StaticInstType == StaticInstType_COND {
		twoBitBranchPredictorUpdate.SaturatingCounter.Update(taken)
	}

	branchPredictor.branchTargetBuffer.Update(branchAddress, branchTarget, taken)
}

func (branchPredictor *TwoBitBranchPredictor) Recover(returnAddressStackRecoverTop uint32) {
	branchPredictor.returnAddressStack.Recover(returnAddressStackRecoverTop)
}
