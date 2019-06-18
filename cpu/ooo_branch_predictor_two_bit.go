package cpu

import "github.com/mcai/heo/simutil"

const (
	BRANCH_SHIFT = 2
)

type BranchTargetBufferEntry struct {
	Source uint32
	Target uint32
}

func NewBranchTargetBufferEntry() *BranchTargetBufferEntry {
	var branchTargetBufferEntry = &BranchTargetBufferEntry{
	}

	return branchTargetBufferEntry
}

type BranchTargetBuffer struct {
	NumSets uint32
	Assoc   uint32
	Entries [][]*BranchTargetBufferEntry
}

func NewBranchTargetBuffer(numSets uint32, assoc uint32) *BranchTargetBuffer {
	var branchTargetBuffer = &BranchTargetBuffer{
		NumSets: numSets,
		Assoc:   assoc,
	}

	for i := uint32(0); i < numSets; i++ {
		var entriesPerSet []*BranchTargetBufferEntry

		for j := uint32(0); j < assoc; j++ {
			entriesPerSet = append(entriesPerSet, NewBranchTargetBufferEntry())
		}

		branchTargetBuffer.Entries = append(branchTargetBuffer.Entries, entriesPerSet)
	}

	return branchTargetBuffer
}

func (branchTargetBuffer *BranchTargetBuffer) GetSet(branchAddress uint32) uint32 {
	return (branchAddress >> BRANCH_SHIFT) & (branchTargetBuffer.NumSets - 1)
}

func (branchTargetBuffer *BranchTargetBuffer) Lookup(branchAddress uint32) *BranchTargetBufferEntry {
	var set = branchTargetBuffer.GetSet(branchAddress)

	for _, entry := range branchTargetBuffer.Entries[set] {
		if entry.Source == branchAddress {
			return entry
		}

	}

	return nil
}

func (branchTargetBuffer *BranchTargetBuffer) Update(branchAddress uint32, branchTarget uint32, taken bool) {
	if !taken {
		return
	}

	var set = branchTargetBuffer.GetSet(branchAddress)

	var entryFound *BranchTargetBufferEntry

	for _, entry := range branchTargetBuffer.Entries[set] {
		if entry.Source == branchAddress {
			entryFound = entry
			break
		}
	}

	if entryFound == nil {
		entryFound = branchTargetBuffer.Entries[set][branchTargetBuffer.Assoc-1]
		entryFound.Source = branchAddress
	}

	entryFound.Target = branchTarget

	branchTargetBuffer.removeFromEntries(set, entryFound)

	branchTargetBuffer.Entries[set] = append(
		[]*BranchTargetBufferEntry{entryFound},
		branchTargetBuffer.Entries[set]...,
	)
}

func (branchTargetBuffer *BranchTargetBuffer) removeFromEntries(set uint32, entryToRemove *BranchTargetBufferEntry) {
	var entriesToReserve []*BranchTargetBufferEntry

	for _, entry := range branchTargetBuffer.Entries[set] {
		if entry != entryToRemove {
			entriesToReserve = append(entriesToReserve, entry)
		}
	}

	branchTargetBuffer.Entries[set] = entriesToReserve
}

type ReturnAddressStack struct {
	size    uint32
	top     uint32
	entries []*BranchTargetBufferEntry
}

func NewReturnAddressStack(size uint32) *ReturnAddressStack {
	var returnAddressStack = &ReturnAddressStack{
		size: size,
		top:  size - 1,
	}

	for i := uint32(0); i < size; i++ {
		returnAddressStack.entries = append(
			returnAddressStack.entries,
			NewBranchTargetBufferEntry(),
		)
	}

	return returnAddressStack
}

func (returnAddressStack *ReturnAddressStack) Size() uint32 {
	return returnAddressStack.size
}

func (returnAddressStack *ReturnAddressStack) Top() uint32 {
	if returnAddressStack.size > 0 {
		return returnAddressStack.top
	}

	return 0
}

func (returnAddressStack *ReturnAddressStack) Recover(top uint32) {
	returnAddressStack.top = top
}

func (returnAddressStack *ReturnAddressStack) Push(branchAddress uint32) {
	returnAddressStack.top = (returnAddressStack.top + 1) % returnAddressStack.size
	returnAddressStack.entries[returnAddressStack.top].Target = branchAddress + 8
}

func (returnAddressStack *ReturnAddressStack) Pop() uint32 {
	var target = returnAddressStack.entries[returnAddressStack.top].Target
	returnAddressStack.top = (returnAddressStack.top + returnAddressStack.size - 1) % returnAddressStack.size
	return target
}

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
		BaseBranchPredictor: NewBaseBranchPredictor(thread),

		branchTargetBuffer: NewBranchTargetBuffer(branchTargetBufferNumSets, branchTargetBufferAssoc),
		returnAddressStack: NewReturnAddressStack(returnAddressStackSize),

		size: size,
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
		}
	}

	return 0, returnAddressStackRecoverTop, branchPredictorUpdate
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
