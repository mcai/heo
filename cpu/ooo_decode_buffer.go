package cpu

type DecodeBufferEntry struct {
	Id                           int32
	DynamicInst                  *DynamicInst
	Npc                          uint32
	Nnpc                         uint32
	PredictedNnpc                uint32
	ReturnAddressStackRecoverTop uint32
	BranchPredictorUpdate        interface{}
	Speculative                  bool
}

func NewDecodeBufferEntry(dynamicInst *DynamicInst, npc uint32, nnpc uint32, predictedNnpc uint32, returnAddressStackRecoverTop uint32, branchPredictorUpdate interface{}, speculative bool) *DecodeBufferEntry {
	var decodeBufferEntry = &DecodeBufferEntry{
		Id:                           dynamicInst.Thread.Core().Processor().Experiment.OoO.CurrentDecodeBufferEntryId,
		DynamicInst:                  dynamicInst,
		Npc:                          npc,
		Nnpc:                         nnpc,
		PredictedNnpc:                predictedNnpc,
		ReturnAddressStackRecoverTop: returnAddressStackRecoverTop,
		BranchPredictorUpdate:        branchPredictorUpdate,
		Speculative:                  speculative,
	}

	dynamicInst.Thread.Core().Processor().Experiment.OoO.CurrentDecodeBufferEntryId++

	return decodeBufferEntry
}
