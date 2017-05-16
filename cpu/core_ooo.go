package cpu

type OoOCore struct {
	*MemoryHierarchyCore

	fuPool                  *FUPool
	waitingInstructionQueue []GeneralReorderBufferEntry
	readyInstructionQueue   []GeneralReorderBufferEntry
	readyLoadQueue          []GeneralReorderBufferEntry
	waitingStoreQueue       []GeneralReorderBufferEntry
	readyStoreQueue         []GeneralReorderBufferEntry
	oooEventQueue           []GeneralReorderBufferEntry

	RegisterRenameScheduler *RoundRobinScheduler
	DispatchScheduler       *RoundRobinScheduler
}

func NewOoOCore(processor *Processor, num int32) *OoOCore {
	var core = &OoOCore{
		MemoryHierarchyCore: NewMemoryHierarchyCore(processor, num),
	}

	core.fuPool = NewFUPool(core)

	var resources []interface{}

	for i := int32(0); i < core.Processor().Experiment.CPUConfig.NumThreadsPerCore; i++ {
		resources = append(resources, i)
	}

	core.RegisterRenameScheduler = NewRoundRobinScheduler(
		resources,
		func(resource interface{}) bool {
			var thread = core.Threads()[resource.(int32)].(*OoOThread)

			if thread.Context() == nil {
				return false
			} else if thread.DecodeBuffer.Empty() {
				return false
			} else if thread.ReorderBuffer.Full() {
				return false
			} else {
				return true
			}
		},
		func(resource interface{}) bool {
			var thread = core.Threads()[resource.(int32)].(*OoOThread)

			return thread.RegisterRenameOne()
		},
		core.Processor().Experiment.CPUConfig.DecodeWidth,
	)

	core.DispatchScheduler = NewRoundRobinScheduler(
		resources,
		func(resource interface{}) bool {
			var thread = core.Threads()[resource.(int32)].(*OoOThread)

			return thread.Context() != nil
		},
		func(resource interface{}) bool {
			var thread = core.Threads()[resource.(int32)].(*OoOThread)

			return thread.DispatchOne()
		},
		core.Processor().Experiment.CPUConfig.DecodeWidth,
	)

	return core
}

func (core *OoOCore) FUPool() *FUPool {
	return core.fuPool
}

func (core *OoOCore) WaitingInstructionQueue() []GeneralReorderBufferEntry {
	return core.waitingInstructionQueue
}

func (core *OoOCore) SetWaitingInstructionQueue(waitingInstructionQueue []GeneralReorderBufferEntry) {
	core.waitingInstructionQueue = waitingInstructionQueue
}

func (core *OoOCore) ReadyInstructionQueue() []GeneralReorderBufferEntry {
	return core.readyInstructionQueue
}

func (core *OoOCore) SetReadyInstructionQueue(readyInstructionQueue []GeneralReorderBufferEntry) {
	core.readyInstructionQueue = readyInstructionQueue
}

func (core *OoOCore) ReadyLoadQueue() []GeneralReorderBufferEntry {
	return core.readyLoadQueue
}

func (core *OoOCore) SetReadyLoadQueue(readyLoadQueue []GeneralReorderBufferEntry) {
	core.readyLoadQueue = readyLoadQueue
}

func (core *OoOCore) WaitingStoreQueue() []GeneralReorderBufferEntry {
	return core.waitingStoreQueue
}

func (core *OoOCore) SetWaitingStoreQueue(waitingStoreQueue []GeneralReorderBufferEntry) {
	core.waitingStoreQueue = waitingStoreQueue
}

func (core *OoOCore) ReadyStoreQueue() []GeneralReorderBufferEntry {
	return core.readyStoreQueue
}

func (core *OoOCore) SetReadyStoreQueue(readyStoreQueue []GeneralReorderBufferEntry) {
	core.readyStoreQueue = readyStoreQueue
}

func (core *OoOCore) OoOEventQueue() []GeneralReorderBufferEntry {
	return core.oooEventQueue
}

func (core *OoOCore) SetOoOEventQueue(oooEventQueue []GeneralReorderBufferEntry) {
	core.oooEventQueue = oooEventQueue
}

func (core *OoOCore) MeasurementOneCycle() {
	core.Commit()
	core.Writeback()
	core.RefreshLoadStoreQueue()
	core.Wakeup()
	core.Issue()
	core.Dispatch()
	core.RegisterRename()
	core.Fetch()
}

func (core *OoOCore) Fetch() {
	for _, thread := range core.Threads() {
		if thread.Context() != nil && thread.Context().State == ContextState_RUNNING {
			thread.(*OoOThread).Fetch()
		}
	}
}

func (core *OoOCore) RegisterRename() {
	core.RegisterRenameScheduler.ConsumeNext()
}

func (core *OoOCore) Dispatch() {
	core.DispatchScheduler.ConsumeNext()
}

func (core *OoOCore) WakeupInstructionQueue() {
	var waitingInstructionQueueToReserve []GeneralReorderBufferEntry

	for _, entry := range core.WaitingInstructionQueue() {
		if entry.AllOperandReady() {
			core.SetReadyInstructionQueue(
				append(
					core.ReadyInstructionQueue(),
					entry,
				),
			)
		} else {
			waitingInstructionQueueToReserve = append(
				waitingInstructionQueueToReserve,
				entry,
			)
		}
	}

	core.SetWaitingInstructionQueue(waitingInstructionQueueToReserve)
}

func (core *OoOCore) WakeupStoreQueue() {
	var waitingStoreQueueToReserve []GeneralReorderBufferEntry

	for _, entry := range core.WaitingStoreQueue() {
		if entry.AllOperandReady() {
			core.SetReadyStoreQueue(
				append(
					core.ReadyStoreQueue(),
					entry,
				),
			)
		} else {
			waitingStoreQueueToReserve = append(
				waitingStoreQueueToReserve,
				entry,
			)
		}
	}

	core.SetWaitingStoreQueue(waitingStoreQueueToReserve)
}

func (core *OoOCore) Wakeup() {
	core.WakeupInstructionQueue()
	core.WakeupStoreQueue()
}

func (core *OoOCore) IssueInstructionQueue(quant uint32) uint32 {
	var readyInstructionQueueToRemove []GeneralReorderBufferEntry

	for _, entry := range core.ReadyInstructionQueue() {
		if quant <= 0 {
			break
		}

		var reorderBufferEntry = entry.(*ReorderBufferEntry)

		if reorderBufferEntry.DynamicInst().StaticInst.Mnemonic.FUOperationType != FUOperationType_NONE {
			if core.FUPool().Acquire(reorderBufferEntry, func() {
				SignalCompleted(reorderBufferEntry)
			}) {
				reorderBufferEntry.SetIssued(true)
			} else {
				continue
			}
		} else {
			reorderBufferEntry.SetIssued(true)
			reorderBufferEntry.SetCompleted(true)
			reorderBufferEntry.Writeback()
		}

		readyInstructionQueueToRemove = append(readyInstructionQueueToRemove, reorderBufferEntry)

		quant--
	}

	for _, entryToRemove := range readyInstructionQueueToRemove {
		core.removeFromReadyInstructionQueue(entryToRemove)
	}

	return quant
}

func (core *OoOCore) IssueLoadQueue(quant uint32) uint32 {
	var readyLoadQueueToRemove []GeneralReorderBufferEntry

	for _, entry := range core.ReadyLoadQueue() {
		if quant <= 0 {
			break
		}

		var loadStoreQueueEntry = entry.(*LoadStoreQueueEntry)

		var hitInLoadStoreQueue = false

		for _, entryFound := range loadStoreQueueEntry.Thread().(*OoOThread).LoadStoreQueue.Entries {
			var loadStoreQueueEntryFound = entryFound.(*LoadStoreQueueEntry)

			if loadStoreQueueEntryFound.DynamicInst().StaticInst.Mnemonic.StaticInstType == StaticInstType_ST &&
				loadStoreQueueEntryFound.EffectiveAddress == loadStoreQueueEntry.EffectiveAddress {
				hitInLoadStoreQueue = true
				break
			}
		}

		if hitInLoadStoreQueue {
			loadStoreQueueEntry.SetIssued(true)
			SignalCompleted(loadStoreQueueEntry)
		} else {
			if !core.CanLoad(loadStoreQueueEntry.Thread(), uint32(loadStoreQueueEntry.EffectiveAddress)) {
				break
			}

			core.Load(
				loadStoreQueueEntry.Thread(),
				uint32(loadStoreQueueEntry.EffectiveAddress),
				loadStoreQueueEntry.DynamicInst().Pc,
				func() {
					SignalCompleted(loadStoreQueueEntry)
				},
			)

			loadStoreQueueEntry.SetIssued(true)
		}

		readyLoadQueueToRemove = append(readyLoadQueueToRemove, loadStoreQueueEntry)

		quant--
	}

	for _, entryToRemove := range readyLoadQueueToRemove {
		core.removeFromReadyLoadQueue(entryToRemove)
	}

	return quant
}

func (core *OoOCore) IssueStoreQueue(quant uint32) uint32 {
	var readyStoreQueueToRemove []GeneralReorderBufferEntry

	for _, entry := range core.ReadyStoreQueue() {
		if quant <= 0 {
			break
		}

		var loadStoreQueueEntry = entry.(*LoadStoreQueueEntry)

		if !core.CanStore(loadStoreQueueEntry.Thread(), uint32(loadStoreQueueEntry.EffectiveAddress)) {
			break
		}

		core.Store(
			loadStoreQueueEntry.Thread(),
			uint32(loadStoreQueueEntry.EffectiveAddress),
			loadStoreQueueEntry.DynamicInst().Pc,
			func() {
			},
		)

		loadStoreQueueEntry.SetIssued(true)
		SignalCompleted(loadStoreQueueEntry)

		readyStoreQueueToRemove = append(readyStoreQueueToRemove, loadStoreQueueEntry)

		quant--
	}

	for _, entryToRemove := range readyStoreQueueToRemove {
		core.removeFromReadyStoreQueue(entryToRemove)
	}

	return quant
}

func (core *OoOCore) Issue() {
	var quant = core.Processor().Experiment.CPUConfig.IssueWidth

	quant = core.IssueInstructionQueue(quant)
	quant = core.IssueLoadQueue(quant)
	quant = core.IssueStoreQueue(quant)
}

func (core *OoOCore) removeFromReadyInstructionQueue(entryToRemove GeneralReorderBufferEntry) {
	var readyInstructionQueueToReserve []GeneralReorderBufferEntry

	for _, entry := range core.ReadyInstructionQueue() {
		if entry != entryToRemove {
			readyInstructionQueueToReserve = append(readyInstructionQueueToReserve, entry)
		}
	}

	core.SetReadyInstructionQueue(readyInstructionQueueToReserve)
}

func (core *OoOCore) removeFromReadyLoadQueue(entryToRemove GeneralReorderBufferEntry) {
	var readyLoadQueueToReserve []GeneralReorderBufferEntry

	for _, entry := range core.ReadyLoadQueue() {
		if entry != entryToRemove {
			readyLoadQueueToReserve = append(readyLoadQueueToReserve, entry)
		}
	}

	core.SetReadyLoadQueue(readyLoadQueueToReserve)
}

func (core *OoOCore) removeFromReadyStoreQueue(entryToRemove GeneralReorderBufferEntry) {
	var readyStoreQueueToReserve []GeneralReorderBufferEntry

	for _, entry := range core.ReadyStoreQueue() {
		if entry != entryToRemove {
			readyStoreQueueToReserve = append(readyStoreQueueToReserve, entry)
		}
	}

	core.SetReadyStoreQueue(readyStoreQueueToReserve)
}

func (core *OoOCore) Writeback() {
	for _, entry := range core.OoOEventQueue() {
		entry.SetCompleted(true)
		entry.Writeback()
	}

	core.SetOoOEventQueue([]GeneralReorderBufferEntry{})
}

func (core *OoOCore) RefreshLoadStoreQueue() {
	for _, thread := range core.Threads() {
		if thread.Context() != nil {
			thread.(*OoOThread).RefreshLoadStoreQueue()
		}
	}
}

func (core *OoOCore) Commit() {
	for _, thread := range core.Threads() {
		if thread.Context() != nil {
			thread.(*OoOThread).Commit()
		}
	}
}

func (core *OoOCore) RemoveFromQueues(entryToRemove GeneralReorderBufferEntry) {
	var waitingInstructionQueueToReserve []GeneralReorderBufferEntry
	var readyInstructionQueueToReserve   []GeneralReorderBufferEntry

	var readyLoadQueueToReserve          []GeneralReorderBufferEntry

	var waitingStoreQueueToReserve       []GeneralReorderBufferEntry
	var readyStoreQueueToReserve         []GeneralReorderBufferEntry

	var oooEventQueueToReserve           []GeneralReorderBufferEntry

	for _, entry := range core.waitingInstructionQueue {
		if entry != entryToRemove {
			waitingInstructionQueueToReserve = append(waitingInstructionQueueToReserve, entry)
		}
	}

	for _, entry := range core.readyInstructionQueue {
		if entry != entryToRemove {
			readyInstructionQueueToReserve = append(readyInstructionQueueToReserve, entry)
		}
	}

	for _, entry := range core.readyLoadQueue {
		if entry != entryToRemove {
			readyLoadQueueToReserve = append(readyLoadQueueToReserve, entry)
		}
	}

	for _, entry := range core.waitingStoreQueue {
		if entry != entryToRemove {
			waitingStoreQueueToReserve = append(waitingStoreQueueToReserve, entry)
		}
	}

	for _, entry := range core.readyStoreQueue {
		if entry != entryToRemove {
			readyStoreQueueToReserve = append(readyStoreQueueToReserve, entry)
		}
	}

	for _, entry := range core.oooEventQueue {
		if entry != entryToRemove {
			oooEventQueueToReserve = append(oooEventQueueToReserve, entry)
		}
	}

	core.waitingInstructionQueue = waitingInstructionQueueToReserve
	core.readyInstructionQueue = readyInstructionQueueToReserve

	core.readyLoadQueue = readyLoadQueueToReserve

	core.waitingStoreQueue = waitingStoreQueueToReserve
	core.readyStoreQueue = readyStoreQueueToReserve

	core.oooEventQueue = oooEventQueueToReserve

	entryToRemove.SetSquashed(true)
}