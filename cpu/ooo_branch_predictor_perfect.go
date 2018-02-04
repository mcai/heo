package cpu

type PerfectBranchPredictor struct {
	*BaseBranchPredictor
}

func NewPerfectBranchPredictor(thread Thread) *PerfectBranchPredictor {
	var branchPredictor = &PerfectBranchPredictor{
		BaseBranchPredictor: NewBaseBranchPredictor(thread),
	}

	return branchPredictor
}

func (branchPredictor *PerfectBranchPredictor) Predict(branchAddress uint32, mnemonic *Mnemonic) (uint32, uint32, interface{}) {
	return branchPredictor.Thread().Context().Regs().Nnpc, 0, nil
}

func (branchPredictor *PerfectBranchPredictor) Update(branchAddress uint32, branchTarget uint32, taken bool, correct bool, mnemonic *Mnemonic, branchPredictorUpdate interface{}) {
	branchPredictor.BaseBranchPredictor.Update(branchAddress, branchTarget, taken, correct, mnemonic, branchPredictorUpdate)
}

func (branchPredictor *PerfectBranchPredictor) Recover(returnAddressStackRecoverTop uint32) {
}
