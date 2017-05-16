package cpu

type GeneralReorderBufferEntry interface {
	Id() int32
	Thread() Thread
	DynamicInst() *DynamicInst

	Npc() uint32
	Nnpc() uint32
	PredictedNnpc() uint32

	ReturnAddressStackRecoverTop() uint32
	BranchPredictorUpdate() interface{}
	Speculative() bool

	OldPhysicalRegisters() map[uint32]*PhysicalRegister

	TargetPhysicalRegisters() map[uint32]*PhysicalRegister
	SetTargetPhysicalRegisters(targetPhysicalRegisters map[uint32]*PhysicalRegister)

	SourcePhysicalRegisters() map[uint32]*PhysicalRegister
	SetSourcePhysicalRegisters(sourcePhysicalRegisters map[uint32]*PhysicalRegister)

	Dispatched() bool
	SetDispatched(dispatched bool)
	Issued() bool
	SetIssued(issued bool)
	Completed() bool
	SetCompleted(completed bool)
	Squashed() bool
	SetSquashed(squashed bool)

	NotReadyOperands() []uint32
	SetNotReadyOperands(notReadyOperands []uint32)
	AddNotReadyOperand(notReadyOperand uint32)
	RemoveNotReadyOperand(notReadyOperand uint32)

	Writeback()

	AllOperandReady() bool
}

func SignalCompleted(reorderBufferEntry GeneralReorderBufferEntry) {
	if !reorderBufferEntry.Squashed() {
		reorderBufferEntry.Thread().Core().SetOoOEventQueue(
			append(
				reorderBufferEntry.Thread().Core().OoOEventQueue(),
				reorderBufferEntry,
			),
		)
	}
}

type BaseReorderBufferEntry struct {
	id                           int32
	thread                       Thread
	dynamicInst                  *DynamicInst

	npc                          uint32
	nnpc                         uint32
	predictedNnpc                uint32

	returnAddressStackRecoverTop uint32
	branchPredictorUpdate        interface{}
	speculative                  bool

	oldPhysicalRegisters         map[uint32]*PhysicalRegister
	targetPhysicalRegisters      map[uint32]*PhysicalRegister
	sourcePhysicalRegisters      map[uint32]*PhysicalRegister

	dispatched                   bool
	issued                       bool
	completed                    bool
	squashed                     bool

	notReadyOperands             []uint32
}

func NewBaseReorderBufferEntry(thread Thread, dynamicInst *DynamicInst, npc uint32, nnpc uint32, predictedNnpc uint32, returnAddressStackRecoverTop uint32, branchPredictorUpdate interface{}, speculative bool) *BaseReorderBufferEntry {
	var reorderBufferEntry = &BaseReorderBufferEntry{
		id:thread.Core().Processor().Experiment.OoO.CurrentReorderBufferEntryId,
		thread:thread,
		dynamicInst:dynamicInst,

		npc:npc,
		nnpc:nnpc,
		predictedNnpc:predictedNnpc,

		returnAddressStackRecoverTop:returnAddressStackRecoverTop,
		branchPredictorUpdate:branchPredictorUpdate,
		speculative:speculative,

		oldPhysicalRegisters:make(map[uint32]*PhysicalRegister),
		targetPhysicalRegisters:make(map[uint32]*PhysicalRegister),
		sourcePhysicalRegisters:make(map[uint32]*PhysicalRegister),
	}

	thread.Core().Processor().Experiment.OoO.CurrentReorderBufferEntryId++

	return reorderBufferEntry
}

func (reorderBufferEntry *BaseReorderBufferEntry) Id() int32 {
	return reorderBufferEntry.id
}

func (reorderBufferEntry *BaseReorderBufferEntry) Thread() Thread {
	return reorderBufferEntry.thread
}

func (reorderBufferEntry *BaseReorderBufferEntry) DynamicInst() *DynamicInst {
	return reorderBufferEntry.dynamicInst
}

func (reorderBufferEntry *BaseReorderBufferEntry) Npc() uint32 {
	return reorderBufferEntry.npc
}

func (reorderBufferEntry *BaseReorderBufferEntry) Nnpc() uint32 {
	return reorderBufferEntry.nnpc
}

func (reorderBufferEntry *BaseReorderBufferEntry) PredictedNnpc() uint32 {
	return reorderBufferEntry.predictedNnpc
}

func (reorderBufferEntry *BaseReorderBufferEntry) ReturnAddressStackRecoverTop() uint32 {
	return reorderBufferEntry.returnAddressStackRecoverTop
}

func (reorderBufferEntry *BaseReorderBufferEntry) BranchPredictorUpdate() interface{} {
	return reorderBufferEntry.branchPredictorUpdate
}

func (reorderBufferEntry *BaseReorderBufferEntry) Speculative() bool {
	return reorderBufferEntry.speculative
}

func (reorderBufferEntry *BaseReorderBufferEntry) OldPhysicalRegisters() map[uint32]*PhysicalRegister {
	return reorderBufferEntry.oldPhysicalRegisters
}

func (reorderBufferEntry *BaseReorderBufferEntry) TargetPhysicalRegisters() map[uint32]*PhysicalRegister {
	return reorderBufferEntry.targetPhysicalRegisters
}

func (reorderBufferEntry *BaseReorderBufferEntry) SetTargetPhysicalRegisters(targetPhysicalRegisters map[uint32]*PhysicalRegister) {
	reorderBufferEntry.targetPhysicalRegisters = targetPhysicalRegisters
}

func (reorderBufferEntry *BaseReorderBufferEntry) SourcePhysicalRegisters() map[uint32]*PhysicalRegister {
	return reorderBufferEntry.sourcePhysicalRegisters
}

func (reorderBufferEntry *BaseReorderBufferEntry) SetSourcePhysicalRegisters(sourcePhysicalRegisters map[uint32]*PhysicalRegister) {
	reorderBufferEntry.sourcePhysicalRegisters = sourcePhysicalRegisters
}

func (reorderBufferEntry *BaseReorderBufferEntry) Dispatched() bool {
	return reorderBufferEntry.dispatched
}

func (reorderBufferEntry *BaseReorderBufferEntry) SetDispatched(dispatched bool) {
	reorderBufferEntry.dispatched = dispatched
}

func (reorderBufferEntry *BaseReorderBufferEntry) Issued() bool {
	return reorderBufferEntry.issued
}

func (reorderBufferEntry *BaseReorderBufferEntry) SetIssued(issued bool) {
	reorderBufferEntry.issued = issued
}

func (reorderBufferEntry *BaseReorderBufferEntry) Completed() bool {
	return reorderBufferEntry.completed
}

func (reorderBufferEntry *BaseReorderBufferEntry) SetCompleted(completed bool) {
	reorderBufferEntry.completed = completed
}

func (reorderBufferEntry *BaseReorderBufferEntry) Squashed() bool {
	return reorderBufferEntry.squashed
}

func (reorderBufferEntry *BaseReorderBufferEntry) SetSquashed(squashed bool) {
	reorderBufferEntry.squashed = squashed
}

func (reorderBufferEntry *BaseReorderBufferEntry) NotReadyOperands() []uint32 {
	return reorderBufferEntry.notReadyOperands
}

func (reorderBufferEntry *BaseReorderBufferEntry) SetNotReadyOperands(notReadyOperands []uint32) {
	reorderBufferEntry.notReadyOperands = notReadyOperands
}

func (reorderBufferEntry *BaseReorderBufferEntry) AddNotReadyOperand(notReadyOperand uint32) {
	reorderBufferEntry.notReadyOperands = append(reorderBufferEntry.notReadyOperands, notReadyOperand)
}

func (reorderBufferEntry *BaseReorderBufferEntry) RemoveNotReadyOperand(notReadyOperand uint32) {
	var notReadyOperandsToReserve []uint32

	for _, o := range reorderBufferEntry.notReadyOperands {
		if o != notReadyOperand {
			notReadyOperandsToReserve = append(notReadyOperandsToReserve, o)
		}
	}

	reorderBufferEntry.notReadyOperands = notReadyOperandsToReserve
}

func (reorderBufferEntry *BaseReorderBufferEntry) doWriteback() {
	for dependency, targetPhysicalRegister := range reorderBufferEntry.targetPhysicalRegisters {
		if dependency != 0 {
			targetPhysicalRegister.Writeback()
		}
	}
}

type ReorderBufferEntry struct {
	*BaseReorderBufferEntry

	EffectiveAddressComputation             bool
	LoadStoreBufferEntry                    *LoadStoreQueueEntry
	EffectiveAddressComputationOperandReady bool
}

func NewReorderBufferEntry(thread Thread, dynamicInst *DynamicInst, npc uint32, nnpc uint32, predictedNnpc uint32, returnAddressStackRecoverIndex uint32, branchPredictorUpdate interface{}, speculative bool) *ReorderBufferEntry {
	var reorderBufferEntry = &ReorderBufferEntry{
		BaseReorderBufferEntry:NewBaseReorderBufferEntry(
			thread,
			dynamicInst,
			npc,
			nnpc,
			predictedNnpc,
			returnAddressStackRecoverIndex,
			branchPredictorUpdate,
			speculative,
		),
	}

	return reorderBufferEntry
}

func (reorderBufferEntry *ReorderBufferEntry) Writeback() {
	if !reorderBufferEntry.EffectiveAddressComputation {
		reorderBufferEntry.doWriteback()
	}
}

func (reorderBufferEntry *ReorderBufferEntry) AllOperandReady() bool {
	if reorderBufferEntry.EffectiveAddressComputation {
		return reorderBufferEntry.EffectiveAddressComputationOperandReady
	}

	return len(reorderBufferEntry.notReadyOperands) == 0
}

type LoadStoreQueueEntry struct {
	*BaseReorderBufferEntry

	EffectiveAddress  int32
	StoreAddressReady bool
}

func NewLoadStoreQueueEntry(thread Thread, dynamicInst *DynamicInst, npc uint32, nnpc uint32, predictedNnpc uint32, returnAddressStackRecoverIndex uint32, branchPredictorUpdate interface{}, speculative bool) *LoadStoreQueueEntry {
	var loadStoreQueueEntry = &LoadStoreQueueEntry{
		BaseReorderBufferEntry:NewBaseReorderBufferEntry(
			thread,
			dynamicInst,
			npc,
			nnpc,
			predictedNnpc,
			returnAddressStackRecoverIndex,
			branchPredictorUpdate,
			speculative,
		),

		EffectiveAddress:-1,
	}

	return loadStoreQueueEntry
}

func (loadStoreQueueEntry *LoadStoreQueueEntry) Writeback() {
	loadStoreQueueEntry.doWriteback()
}

func (loadStoreQueueEntry *LoadStoreQueueEntry) AllOperandReady() bool {
	return len(loadStoreQueueEntry.notReadyOperands) == 0
}
