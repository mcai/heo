package cpu

import (
	"fmt"
	"github.com/mcai/heo/cpu/regs"
)

type OoOThread struct {
	*MemoryHierarchyThread

	BranchPredictor BranchPredictor

	IntPhysicalRegs  *PhysicalRegisterFile
	FpPhysicalRegs   *PhysicalRegisterFile
	MiscPhysicalRegs *PhysicalRegisterFile

	RenameTable map[uint32]*PhysicalRegister

	DecodeBuffer   *PipelineBuffer
	ReorderBuffer  *PipelineBuffer
	LoadStoreQueue *PipelineBuffer

	FetchNpc  uint32
	FetchNnpc uint32

	lastDecodedDynamicInst          *DynamicInst
	lastDecodedDynamicInstCommitted bool

	LastCommitCycle                        int64
	noDynamicInstCommittedCounterThreshold int64
}

func NewOoOThread(core Core, num int32) *OoOThread {
	var thread = &OoOThread{
		MemoryHierarchyThread: NewMemoryHierarchyThread(core, num),

		IntPhysicalRegs:  NewPhysicalRegisterFile(RegisterDependencyType_INT, core.Processor().Experiment.CPUConfig.PhysicalRegisterFileSize),
		FpPhysicalRegs:   NewPhysicalRegisterFile(RegisterDependencyType_FP, core.Processor().Experiment.CPUConfig.PhysicalRegisterFileSize),
		MiscPhysicalRegs: NewPhysicalRegisterFile(RegisterDependencyType_MISC, core.Processor().Experiment.CPUConfig.PhysicalRegisterFileSize),

		RenameTable: make(map[uint32]*PhysicalRegister),

		DecodeBuffer:   NewPipelineBuffer(core.Processor().Experiment.CPUConfig.DecodeBufferSize),
		ReorderBuffer:  NewPipelineBuffer(core.Processor().Experiment.CPUConfig.ReorderBufferSize),
		LoadStoreQueue: NewPipelineBuffer(core.Processor().Experiment.CPUConfig.LoadStoreQueueSize),
	}

	switch core.Processor().Experiment.CPUConfig.BranchPredictorType {
	case BranchPredictorType_PERFECT:
		thread.BranchPredictor = NewPerfectBranchPredictor(thread)
	case BranchPredictorType_TWO_BIT:
		thread.BranchPredictor = NewTwoBitBranchPredictor(
			thread,
			core.Processor().Experiment.CPUConfig.TwoBitBranchPredictorSize,
			core.Processor().Experiment.CPUConfig.BranchTargetBufferNumSets,
			core.Processor().Experiment.CPUConfig.BranchTargetBufferAssoc,
			core.Processor().Experiment.CPUConfig.ReturnAddressStackSize,
		)
	default:
		panic("Impossible")
	}

	for i := uint32(0); i < regs.NUM_INT_REGISTERS; i++ {
		var dependency = RegisterDependencyToInt(RegisterDependencyType_INT, i)
		var physicalReg = thread.IntPhysicalRegs.PhysicalRegisters[i]
		physicalReg.Reserve(dependency)
		thread.RenameTable[dependency] = physicalReg
	}

	for i := uint32(0); i < regs.NUM_FP_REGISTERS; i++ {
		var dependency = RegisterDependencyToInt(RegisterDependencyType_FP, i)
		var physicalReg = thread.FpPhysicalRegs.PhysicalRegisters[i]
		physicalReg.Reserve(dependency)
		thread.RenameTable[dependency] = physicalReg
	}

	for i := uint32(0); i < regs.NUM_MISC_REGISTERS; i++ {
		var dependency = RegisterDependencyToInt(RegisterDependencyType_MISC, i)
		var physicalReg = thread.MiscPhysicalRegs.PhysicalRegisters[i]
		physicalReg.Reserve(dependency)
		thread.RenameTable[dependency] = physicalReg
	}

	return thread
}

func (thread *OoOThread) UpdateFetchNpcAndNnpcFromRegs() {
	thread.FetchNpc = thread.Context().Regs().Npc
	thread.FetchNnpc = thread.Context().Regs().Nnpc

	thread.LastCommitCycle = thread.Core().Processor().Experiment.CycleAccurateEventQueue().CurrentCycle
}

func (thread *OoOThread) TryFetch() bool {
	if thread.FetchStalled {
		return false
	}

	var cacheLineToFetch = thread.Core().L1IController().Cache.GetTag(thread.FetchNpc)
	if int32(cacheLineToFetch) != thread.LastFetchedCacheLine {
		if !thread.Core().CanIfetch(thread, thread.FetchNpc) {
			return false
		} else {
			thread.Core().Ifetch(thread, thread.FetchNpc, thread.FetchNpc, func() {
				thread.FetchStalled = false
			})

			thread.FetchStalled = true
			thread.LastFetchedCacheLine = int32(cacheLineToFetch)

			return false
		}
	}

	return true
}

func (thread *OoOThread) Fetch() {
	if !thread.TryFetch() {
		return
	}

	var hasDone = false

	for !hasDone {
		if thread.Context().State != ContextState_RUNNING {
			break
		}

		if thread.DecodeBuffer.Full() {
			break
		}

		if thread.Context().Regs().Npc != thread.FetchNpc {
			if thread.Context().Speculative {
				thread.Context().Regs().Npc = thread.FetchNpc
			} else {
				thread.Context().EnterSpeculativeState()
			}
		}

		var dynamicInst *DynamicInst

		for {
			var staticInst = thread.Context().DecodeNextStaticInst()

			dynamicInst = NewDynamicInst(thread, thread.Context().Regs().Pc, staticInst)

			staticInst.Execute(thread.Context())

			if dynamicInst.StaticInst.Mnemonic.StaticInstType == StaticInstType_NOP {
				thread.UpdateFetchNpcAndNnpcFromRegs()
			} else {
				break
			}
		}

		thread.FetchNpc = thread.FetchNnpc

		if !thread.Context().Speculative && thread.Context().State != ContextState_RUNNING {
			thread.lastDecodedDynamicInst = dynamicInst
			thread.lastDecodedDynamicInstCommitted = false
		}

		if (thread.FetchNpc+4)%thread.Core().L1IController().Cache.LineSize() == 0 {
			hasDone = true
		}

		var returnAddressStackRecoverTop uint32

		var branchPredictorUpdate interface{}

		if dynamicInst.StaticInst.Mnemonic.StaticInstType.IsControl() {
			thread.FetchNnpc, returnAddressStackRecoverTop, branchPredictorUpdate = thread.BranchPredictor.Predict(thread.FetchNpc, dynamicInst.StaticInst.Mnemonic)
		} else {
			thread.FetchNnpc, returnAddressStackRecoverTop, branchPredictorUpdate = thread.FetchNpc+4, 0, NewTwoBitBranchPredictorUpdate()
		}

		if thread.FetchNnpc != thread.FetchNpc+4 {
			hasDone = true
		}

		thread.DecodeBuffer.Entries = append(
			thread.DecodeBuffer.Entries,
			NewDecodeBufferEntry(
				dynamicInst,
				thread.Context().Regs().Npc,
				thread.Context().Regs().Nnpc,
				thread.FetchNnpc,
				returnAddressStackRecoverTop,
				branchPredictorUpdate,
				thread.Context().Speculative,
			),
		)
	}
}

func (thread *OoOThread) RegisterRenameOne() bool {
	var decodeBufferEntry = thread.DecodeBuffer.Entries[0].(*DecodeBufferEntry)

	var dynamicInst = decodeBufferEntry.DynamicInst

	for outputDependencyType, numPhysicalRegistersToAllocate := range dynamicInst.StaticInst.NumPhysicalRegistersToAllocate {
		if thread.GetPhysicalRegisterFile(outputDependencyType).NumFreePhysicalRegisters < numPhysicalRegistersToAllocate {
			return false
		}
	}

	if dynamicInst.StaticInst.Mnemonic.StaticInstType.IsLoadOrStore() && thread.LoadStoreQueue.Full() {
		return false
	}

	var reorderBufferEntry = NewReorderBufferEntry(
		thread,
		dynamicInst,
		decodeBufferEntry.Npc,
		decodeBufferEntry.Nnpc,
		decodeBufferEntry.PredictedNnpc,
		decodeBufferEntry.ReturnAddressStackRecoverTop,
		decodeBufferEntry.BranchPredictorUpdate,
		decodeBufferEntry.Speculative,
	)

	reorderBufferEntry.EffectiveAddressComputation = dynamicInst.StaticInst.Mnemonic.StaticInstType.IsLoadOrStore()

	for _, inputDependency := range dynamicInst.StaticInst.InputDependencies {
		reorderBufferEntry.SourcePhysicalRegisters()[inputDependency] = thread.RenameTable[inputDependency]
	}

	for _, outputDependency := range dynamicInst.StaticInst.OutputDependencies {
		if outputDependency != 0 {
			var outputDependencyType, _ = RegisterDependencyFromInt(outputDependency)

			reorderBufferEntry.OldPhysicalRegisters()[outputDependency] = thread.RenameTable[outputDependency]
			var physicalReg = thread.GetPhysicalRegisterFile(outputDependencyType).Allocate(reorderBufferEntry, outputDependency)
			thread.RenameTable[outputDependency] = physicalReg
			reorderBufferEntry.TargetPhysicalRegisters()[outputDependency] = physicalReg
		}
	}

	for _, sourcePhysicalReg := range reorderBufferEntry.SourcePhysicalRegisters() {
		if !sourcePhysicalReg.Ready() {
			reorderBufferEntry.AddNotReadyOperand(uint32(sourcePhysicalReg.Dependency))
			sourcePhysicalReg.Dependents = append(sourcePhysicalReg.Dependents, reorderBufferEntry)
		}
	}

	if reorderBufferEntry.EffectiveAddressComputation {
		var physicalReg = reorderBufferEntry.SourcePhysicalRegisters()[dynamicInst.StaticInst.InputDependencies[0]]

		if !physicalReg.Ready() {
			physicalReg.EffectiveAddressComputationOperandDependents = append(
				physicalReg.EffectiveAddressComputationOperandDependents,
				reorderBufferEntry,
			)
		} else {
			reorderBufferEntry.EffectiveAddressComputationOperandReady = true
		}
	}

	if dynamicInst.StaticInst.Mnemonic.StaticInstType.IsLoadOrStore() {
		var loadStoreQueueEntry = NewLoadStoreQueueEntry(
			thread,
			dynamicInst,
			decodeBufferEntry.Npc,
			decodeBufferEntry.Nnpc,
			decodeBufferEntry.PredictedNnpc,
			0,
			nil,
			false,
		)

		loadStoreQueueEntry.EffectiveAddress = dynamicInst.EffectiveAddress

		loadStoreQueueEntry.SetSourcePhysicalRegisters(reorderBufferEntry.SourcePhysicalRegisters())
		loadStoreQueueEntry.SetTargetPhysicalRegisters(reorderBufferEntry.TargetPhysicalRegisters())

		for _, sourcePhysicalReg := range loadStoreQueueEntry.SourcePhysicalRegisters() {
			if !sourcePhysicalReg.Ready() {
				sourcePhysicalReg.Dependents = append(sourcePhysicalReg.Dependents, loadStoreQueueEntry)
			}
		}

		loadStoreQueueEntry.SetNotReadyOperands(reorderBufferEntry.NotReadyOperands())

		var storeAddressPhysicalReg = loadStoreQueueEntry.SourcePhysicalRegisters()[dynamicInst.StaticInst.InputDependencies[0]]

		if !storeAddressPhysicalReg.Ready() {
			storeAddressPhysicalReg.StoreAddressDependents = append(
				storeAddressPhysicalReg.StoreAddressDependents,
				loadStoreQueueEntry,
			)
		} else {
			loadStoreQueueEntry.StoreAddressReady = true
		}

		thread.LoadStoreQueue.Entries = append(thread.LoadStoreQueue.Entries, loadStoreQueueEntry)

		reorderBufferEntry.LoadStoreQueueEntry = loadStoreQueueEntry
	}

	thread.ReorderBuffer.Entries = append(thread.ReorderBuffer.Entries, reorderBufferEntry)

	thread.DecodeBuffer.Entries = thread.DecodeBuffer.Entries[1:]

	return true
}

func (thread *OoOThread) DispatchOne() bool {
	for _, entry := range thread.ReorderBuffer.Entries {
		var reorderBufferEntry = entry.(*ReorderBufferEntry)

		if !reorderBufferEntry.Dispatched() {
			if reorderBufferEntry.AllOperandReady() {
				thread.Core().SetReadyInstructionQueue(
					append(
						thread.Core().ReadyInstructionQueue(),
						reorderBufferEntry,
					),
				)
			} else {
				thread.Core().SetWaitingInstructionQueue(
					append(
						thread.Core().WaitingInstructionQueue(),
						reorderBufferEntry,
					),
				)
			}

			reorderBufferEntry.SetDispatched(true)

			if reorderBufferEntry.LoadStoreQueueEntry != nil {
				var loadStoreQueueEntry = reorderBufferEntry.LoadStoreQueueEntry

				if loadStoreQueueEntry.DynamicInst().StaticInst.Mnemonic.StaticInstType == StaticInstType_ST {
					if loadStoreQueueEntry.AllOperandReady() {
						thread.Core().SetReadyStoreQueue(
							append(
								thread.Core().ReadyStoreQueue(),
								loadStoreQueueEntry,
							),
						)
					} else {
						thread.Core().SetWaitingStoreQueue(
							append(
								thread.Core().WaitingStoreQueue(),
								loadStoreQueueEntry,
							),
						)
					}
				}

				loadStoreQueueEntry.SetDispatched(true)
			}

			return true
		}
	}

	return false
}

func (thread *OoOThread) RefreshLoadStoreQueue() {
	var stdUnknowns []int32

	for _, entry := range thread.LoadStoreQueue.Entries {
		var loadStoreQueueEntry = entry.(*LoadStoreQueueEntry)

		if loadStoreQueueEntry.DynamicInst().StaticInst.Mnemonic.StaticInstType == StaticInstType_ST {
			if loadStoreQueueEntry.StoreAddressReady {
				break
			} else if !loadStoreQueueEntry.AllOperandReady() {
				stdUnknowns = append(stdUnknowns, loadStoreQueueEntry.EffectiveAddress)
			} else {
				for i, stdUnknown := range stdUnknowns {
					if stdUnknown == loadStoreQueueEntry.EffectiveAddress {
						stdUnknowns[i] = -1
					}
				}
			}
		}

		if loadStoreQueueEntry.DynamicInst().StaticInst.Mnemonic.StaticInstType == StaticInstType_LD &&
			loadStoreQueueEntry.Dispatched() &&
			!loadStoreQueueEntry.Issued() &&
			!loadStoreQueueEntry.Completed() &&
			loadStoreQueueEntry.AllOperandReady() {
			var foundInReadyLoadQueue bool

			for _, readyLoad := range thread.Core().ReadyLoadQueue() {
				if readyLoad == loadStoreQueueEntry {
					foundInReadyLoadQueue = true
					break
				}
			}

			var foundInStdUnknowns bool

			for _, stdUnknown := range stdUnknowns {
				if stdUnknown == loadStoreQueueEntry.EffectiveAddress {
					foundInStdUnknowns = true
					break
				}
			}

			if !foundInReadyLoadQueue && !foundInStdUnknowns {
				thread.Core().SetReadyLoadQueue(
					append(
						thread.Core().ReadyLoadQueue(),
						loadStoreQueueEntry,
					),
				)
			}
		}
	}
}

func (thread *OoOThread) DumpQueues() {
	for i, entry := range thread.DecodeBuffer.Entries {
		var decodeBufferEntry = entry.(*DecodeBufferEntry)

		fmt.Printf("thread.decodeBuffer[%d]={id=%d, speculative=%t}\n", i, decodeBufferEntry.Id, decodeBufferEntry.Speculative)
	}

	for i, entry := range thread.ReorderBuffer.Entries {
		var reorderBufferEntry = entry.(*ReorderBufferEntry)

		var loadStoreQueueEntryId = int32(-1)

		if reorderBufferEntry.LoadStoreQueueEntry != nil {
			loadStoreQueueEntryId = reorderBufferEntry.LoadStoreQueueEntry.Id()
		}

		fmt.Printf(
			"thread.reorderBuffer[%d]={id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, loadStoreQueueEntry.id=%d, notReadyOperands=%+v, allOperandReady=%t}\n",
			i,
			reorderBufferEntry.Id(),
			reorderBufferEntry.Dispatched(),
			reorderBufferEntry.Issued(),
			reorderBufferEntry.Completed(),
			reorderBufferEntry.Squashed(),
			loadStoreQueueEntryId,
			reorderBufferEntry.NotReadyOperands(),
			reorderBufferEntry.AllOperandReady(),
		)
	}

	for i, entry := range thread.LoadStoreQueue.Entries {
		var loadStoreQueueEntry = entry.(*LoadStoreQueueEntry)

		fmt.Printf(
			"thread.loadStoreQueue[%d]={id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, notReadyOperands=%+v, allOperandReady=%t}\n",
			i,
			loadStoreQueueEntry.Id(),
			loadStoreQueueEntry.Dispatched(),
			loadStoreQueueEntry.Issued(),
			loadStoreQueueEntry.Completed(),
			loadStoreQueueEntry.Squashed(),
			loadStoreQueueEntry.NotReadyOperands(),
			loadStoreQueueEntry.AllOperandReady(),
		)
	}

	for dependency := uint32(0); dependency < regs.NUM_INT_REGISTERS+regs.NUM_FP_REGISTERS+regs.NUM_MISC_REGISTERS; dependency++ {
		var physicalReg = thread.RenameTable[dependency]

		fmt.Printf("thread.renameTable[%d]=PhysicalRegister{type=%s, num=%d, dependency=%d, state=%s}\n",
			dependency, physicalReg.PhysicalRegisterFile.RegisterDependencyType, physicalReg.Num, physicalReg.Dependency, physicalReg.State)
	}

	for fuType, fuDescriptor := range thread.Core().FUPool().FUDescriptors {
		fmt.Printf("thread.core.fuPool.descriptors[%s]={numFree=%d, quantity=%d}\n", fuType, fuDescriptor.NumFree, fuDescriptor.Quantity)
	}

	for i, entry := range thread.Core().WaitingInstructionQueue() {
		fmt.Printf(
			"thread.core.waitingInstructionQueue[%d]={id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, notReadyOperands=%+v, allOperandReady=%t}\n",
			i,
			entry.Id(),
			entry.Dispatched(),
			entry.Issued(),
			entry.Completed(),
			entry.Squashed(),
			entry.NotReadyOperands(),
			entry.AllOperandReady(),
		)
	}

	for i, entry := range thread.Core().ReadyInstructionQueue() {
		fmt.Printf(
			"thread.core.readyInstructionQueue[%d]={id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, notReadyOperands=%+v, allOperandReady=%t}\n",
			i,
			entry.Id(),
			entry.Dispatched(),
			entry.Issued(),
			entry.Completed(),
			entry.Squashed(),
			entry.NotReadyOperands(),
			entry.AllOperandReady(),
		)
	}

	for i, entry := range thread.Core().ReadyLoadQueue() {
		fmt.Printf(
			"thread.core.readyLoadQueue[%d]={id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, notReadyOperands=%+v, allOperandReady=%t}\n",
			i,
			entry.Id(),
			entry.Dispatched(),
			entry.Issued(),
			entry.Completed(),
			entry.Squashed(),
			entry.NotReadyOperands(),
			entry.AllOperandReady(),
		)
	}

	for i, entry := range thread.Core().WaitingStoreQueue() {
		fmt.Printf(
			"thread.core.waitingStoreQueue[%d]={id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, notReadyOperands=%+v, allOperandReady=%t}\n",
			i,
			entry.Id(),
			entry.Dispatched(),
			entry.Issued(),
			entry.Completed(),
			entry.Squashed(),
			entry.NotReadyOperands(),
			entry.AllOperandReady(),
		)
	}

	for i, entry := range thread.Core().ReadyStoreQueue() {
		fmt.Printf(
			"thread.core.readyStoreQueue[%d]={id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, notReadyOperands=%+v, allOperandReady=%t}\n",
			i,
			entry.Id(),
			entry.Dispatched(),
			entry.Issued(),
			entry.Completed(),
			entry.Squashed(),
			entry.NotReadyOperands(),
			entry.AllOperandReady(),
		)
	}

	for i, entry := range thread.Core().OoOEventQueue() {
		fmt.Printf(
			"thread.core.oooEventQueue[%d]={id=%d, dispatched=%t, issued=%t, completed=%t, squashed=%t, notReadyOperands=%+v, allOperandReady=%t}\n",
			i,
			entry.Id(),
			entry.Dispatched(),
			entry.Issued(),
			entry.Completed(),
			entry.Squashed(),
			entry.NotReadyOperands(),
			entry.AllOperandReady(),
		)
	}

	fmt.Println("thread.intPhysicalRegs:")

	thread.IntPhysicalRegs.Dump()

	fmt.Println("thread.fpPhysicalRegs:")

	thread.FpPhysicalRegs.Dump()

	fmt.Println("thread.miscPhysicalRegs:")

	thread.MiscPhysicalRegs.Dump()
}

func (thread *OoOThread) Commit() {
	var commitTimeout = int64(1000000)

	if thread.Core().Processor().Experiment.CycleAccurateEventQueue().CurrentCycle-thread.LastCommitCycle > commitTimeout {
		if thread.noDynamicInstCommittedCounterThreshold > 5 {
			thread.DumpQueues()
			thread.Core().Processor().Experiment.MemoryHierarchy.DumpPendingFlowTree()
			panic(fmt.Sprintf(
				"[%d] No dynamic insts committed for a long time (thread.NumDynamicInsts=%d)",
				thread.Core().Processor().Experiment.CycleAccurateEventQueue().CurrentCycle,
				thread.NumDynamicInsts(),
			))
		} else {
			thread.LastCommitCycle = thread.Core().Processor().Experiment.CycleAccurateEventQueue().CurrentCycle
			thread.noDynamicInstCommittedCounterThreshold++
		}
	}

	var numCommitted = uint32(0)

	for !thread.ReorderBuffer.Empty() && numCommitted < thread.Core().Processor().Experiment.CPUConfig.CommitWidth {
		var reorderBufferEntry = thread.ReorderBuffer.Entries[0].(*ReorderBufferEntry)

		if !reorderBufferEntry.Completed() {
			break
		}

		if reorderBufferEntry.Speculative() {
			thread.BranchPredictor.Recover(reorderBufferEntry.ReturnAddressStackRecoverTop())

			thread.Context().ExitSpeculativeState()

			thread.FetchNpc = thread.Context().Regs().Npc
			thread.FetchNnpc = thread.Context().Regs().Nnpc

			thread.Squash()
			break
		}

		if reorderBufferEntry.EffectiveAddressComputation {
			var loadStoreQueueEntry = reorderBufferEntry.LoadStoreQueueEntry

			if !loadStoreQueueEntry.Completed() {
				break
			}

			thread.Core().RemoveFromQueues(loadStoreQueueEntry)

			thread.removeFromLoadStoreQueue(loadStoreQueueEntry)
		}

		for _, outputDependency := range reorderBufferEntry.DynamicInst().StaticInst.OutputDependencies {
			if outputDependency != 0 {
				reorderBufferEntry.OldPhysicalRegisters()[outputDependency].Reclaim()
				reorderBufferEntry.TargetPhysicalRegisters()[outputDependency].Commit()
			}
		}

		if reorderBufferEntry.DynamicInst().StaticInst.Mnemonic.StaticInstType.IsControl() {
			thread.BranchPredictor.Update(
				reorderBufferEntry.DynamicInst().Pc,
				reorderBufferEntry.Nnpc(),
				reorderBufferEntry.Nnpc() != reorderBufferEntry.Npc()+4,
				reorderBufferEntry.PredictedNnpc() == reorderBufferEntry.Nnpc(),
				reorderBufferEntry.DynamicInst().StaticInst.Mnemonic,
				reorderBufferEntry.BranchPredictorUpdate(),
			)
		}

		thread.Core().RemoveFromQueues(reorderBufferEntry)

		if thread.Context().State == ContextState_FINISHED && reorderBufferEntry.DynamicInst() == thread.lastDecodedDynamicInst {
			thread.lastDecodedDynamicInstCommitted = true
		}

		thread.ReorderBuffer.Entries = thread.ReorderBuffer.Entries[1:]

		thread.numDynamicInsts++

		thread.LastCommitCycle = thread.Core().Processor().Experiment.CycleAccurateEventQueue().CurrentCycle

		numCommitted++
	}
}

func (thread *OoOThread) removeFromLoadStoreQueue(entryToRemove *LoadStoreQueueEntry) {
	var loadStoreQueueEntriesToReserve []interface{}

	for _, entry := range thread.LoadStoreQueue.Entries {
		if entry != entryToRemove {
			loadStoreQueueEntriesToReserve = append(loadStoreQueueEntriesToReserve, entry)
		}
	}

	thread.LoadStoreQueue.Entries = loadStoreQueueEntriesToReserve
}

func (thread *OoOThread) Squash() {
	for !thread.ReorderBuffer.Empty() {
		var reorderBufferEntry = thread.ReorderBuffer.Entries[len(thread.ReorderBuffer.Entries)-1].(*ReorderBufferEntry)

		if reorderBufferEntry.EffectiveAddressComputation {
			var loadStoreQueueEntry = reorderBufferEntry.LoadStoreQueueEntry

			thread.Core().RemoveFromQueues(loadStoreQueueEntry)

			thread.removeFromLoadStoreQueue(loadStoreQueueEntry)
		}

		thread.Core().RemoveFromQueues(reorderBufferEntry)

		for _, outputDependency := range reorderBufferEntry.DynamicInst().StaticInst.OutputDependencies {
			if outputDependency != 0 {
				reorderBufferEntry.TargetPhysicalRegisters()[outputDependency].Recover()
				thread.RenameTable[outputDependency] = reorderBufferEntry.OldPhysicalRegisters()[outputDependency]
			}
		}

		reorderBufferEntry.SetTargetPhysicalRegisters(make(map[uint32]*PhysicalRegister))

		thread.ReorderBuffer.Entries = thread.ReorderBuffer.Entries[:len(thread.ReorderBuffer.Entries)-1]
	}

	if !thread.ReorderBuffer.Empty() || !thread.LoadStoreQueue.Empty() {
		panic("Impossible")
	}

	thread.Core().FUPool().ReleaseAll()

	thread.DecodeBuffer.Entries = []interface{}{}
}

func (thread *OoOThread) IsLastDecodedDynamicInstCommitted() bool {
	return thread.lastDecodedDynamicInst == nil || thread.lastDecodedDynamicInstCommitted
}

func (thread *OoOThread) GetPhysicalRegisterFile(dependencyType RegisterDependencyType) *PhysicalRegisterFile {
	switch dependencyType {
	case RegisterDependencyType_INT:
		return thread.IntPhysicalRegs
	case RegisterDependencyType_FP:
		return thread.FpPhysicalRegs
	case RegisterDependencyType_MISC:
		return thread.MiscPhysicalRegs
	default:
		panic("Impossible")
	}
}
