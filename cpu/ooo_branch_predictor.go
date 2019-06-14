package cpu

type BranchPredictorType string

const (
	BranchPredictorType_PERFECT = BranchPredictorType("PERFECT")

	BranchPredictorType_TWO_BIT = BranchPredictorType("TWO_BIT")
)

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
	return branchAddress >> BRANCH_SHIFT & (branchTargetBuffer.NumSets - 1)
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

type BranchPredictor interface {
	Thread() Thread

	Predict(branchAddress uint32, mnemonic *Mnemonic) (uint32, uint32, interface{})
	Update(branchAddress uint32, branchTarget uint32, taken bool, correct bool, mnemonic *Mnemonic, branchPredictorUpdate interface{})
	Recover(returnAddressStackRecoverTop uint32)

	NumHits() int64
	NumMisses() int64
	NumAccesses() int64
	HitRatio() float64
}

type BaseBranchPredictor struct {
	thread    Thread
	numHits   int64
	numMisses int64
}

func NewBaseBranchPredictor(thread Thread) *BaseBranchPredictor {
	var branchPredictor = &BaseBranchPredictor{
		thread: thread,
	}

	return branchPredictor
}

func (branchPredictor *BaseBranchPredictor) Thread() Thread {
	return branchPredictor.thread
}

func (branchPredictor *BaseBranchPredictor) NumHits() int64 {
	return branchPredictor.numHits
}

func (branchPredictor *BaseBranchPredictor) NumMisses() int64 {
	return branchPredictor.numMisses
}

func (branchPredictor *BaseBranchPredictor) NumAccesses() int64 {
	return branchPredictor.numHits + branchPredictor.numMisses
}

func (branchPredictor *BaseBranchPredictor) HitRatio() float64 {
	if branchPredictor.NumAccesses() == 0 {
		return float64(0)
	}

	return float64(branchPredictor.numHits) / float64(branchPredictor.NumAccesses())
}

func (branchPredictor *BaseBranchPredictor) Update(branchAddress uint32, branchTarget uint32, taken bool, correct bool, mnemonic *Mnemonic, branchPredictorUpdate interface{}) {
	if correct {
		branchPredictor.numHits++
	} else {
		branchPredictor.numMisses++
	}
}
