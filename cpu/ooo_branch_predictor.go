package cpu

type BranchPredictorType string

const (
	BranchPredictorType_PERFECT = BranchPredictorType("PERFECT")

	BranchPredictorType_TWO_BIT = BranchPredictorType("TWO_BIT")
)

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
