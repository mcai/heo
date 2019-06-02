package uncore

import (
	"github.com/mcai/heo/simutil"
)

type DirectoryEntry struct {
	Owner   *CacheController
	Sharers []*CacheController
}

func NewDirectoryEntry() *DirectoryEntry {
	var directoryEntry = &DirectoryEntry{
	}

	return directoryEntry
}

type DirectoryControllerFiniteStateMachine struct {
	*simutil.BaseFiniteStateMachine
	DirectoryController *DirectoryController
	DirectoryEntry      *DirectoryEntry
	Set                 uint32
	Way                 uint32
	NumRecallAcks       int32
	StalledEvents       []func()
	OnCompletedCallback func()
	EvicterTag          int32
	VictimTag           int32
}

func NewDirectoryControllerFiniteStateMachine(set uint32, way uint32, directoryController *DirectoryController) *DirectoryControllerFiniteStateMachine {
	var fsm = &DirectoryControllerFiniteStateMachine{
		BaseFiniteStateMachine: simutil.NewBaseFiniteStateMachine(DirectoryControllerState_I),
		DirectoryController:    directoryController,
		DirectoryEntry:         NewDirectoryEntry(),
		Set:                    set,
		Way:                    way,
		EvicterTag:             INVALID_TAG,
		VictimTag:              INVALID_TAG,
	}

	return fsm
}

func (fsm *DirectoryControllerFiniteStateMachine) Valid() bool {
	return fsm.State() != DirectoryControllerState_I
}

func (fsm *DirectoryControllerFiniteStateMachine) Line() *CacheLine {
	return fsm.DirectoryController.Cache.Sets[fsm.Set].Lines[fsm.Way]
}

func (fsm *DirectoryControllerFiniteStateMachine) fireTransition(event DirectoryControllerEvent) {
	event.Complete()
	fsm.DirectoryController.FsmFactory.FireTransition(fsm, event.EventType(), event)
}

func (fsm *DirectoryControllerFiniteStateMachine) OnEventGetS(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController, onStalledCallback func()) {
	var event = NewGetSEvent(
		fsm.DirectoryController,
		producerFlow,
		producerFlow.Access(),
		tag,
		requester,
		fsm.Set,
		fsm.Way,
		onStalledCallback,
	)

	fsm.fireTransition(event)
}

func (fsm *DirectoryControllerFiniteStateMachine) OnEventGetM(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController, onStalledCallback func()) {
	var event = NewGetMEvent(
		fsm.DirectoryController,
		producerFlow,
		producerFlow.Access(),
		tag,
		requester,
		fsm.Set,
		fsm.Way,
		onStalledCallback,
	)

	fsm.fireTransition(event)
}

func (fsm *DirectoryControllerFiniteStateMachine) OnEventReplacement(producerFlow CacheCoherenceFlow, tag uint32, cacheAccess *CacheAccess, requester *CacheController, onCompletedCallback func(), onStalledCallback func()) {
	var event = NewDirReplacementEvent(
		fsm.DirectoryController,
		producerFlow,
		producerFlow.Access(),
		tag,
		cacheAccess,
		fsm.Set,
		fsm.Way,
		onCompletedCallback,
		onStalledCallback,
	)

	fsm.fireTransition(event)
}

func (fsm *DirectoryControllerFiniteStateMachine) OnEventRecallAck(producerFlow CacheCoherenceFlow, tag uint32, sender *CacheController) {
	var event = NewRecallAckEvent(
		fsm.DirectoryController,
		producerFlow,
		producerFlow.Access(),
		tag,
		sender,
	)

	fsm.fireTransition(event)

	if fsm.NumRecallAcks == 0 {
		var event = NewLastRecallAckEvent(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
		)

		fsm.fireTransition(event)
	}
}

func (fsm *DirectoryControllerFiniteStateMachine) OnEventPutS(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	if len(fsm.DirectoryEntry.Sharers) > 1 {
		var event = NewPutSNotLastEvent(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
			requester,
		)

		fsm.fireTransition(event)
	} else {
		var event = NewPutSLastEvent(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
			requester,
		)

		fsm.fireTransition(event)
	}
}

func (fsm *DirectoryControllerFiniteStateMachine) OnEventPutMAndData(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	if requester == fsm.DirectoryEntry.Owner {
		var event = NewPutMAndDataFromOwnerEvent(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
			requester,
		)

		fsm.fireTransition(event)
	} else {
		var event = NewPutMAndDataFromNonOwnerEvent(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
			requester,
		)

		fsm.fireTransition(event)
	}
}

func (directoryControllerFsm *DirectoryControllerFiniteStateMachine) OnEventData(producerFlow CacheCoherenceFlow, tag uint32, sender *CacheController) {
	var event = NewDataEvent(
		directoryControllerFsm.DirectoryController,
		producerFlow,
		producerFlow.Access(),
		tag,
		sender,
	)

	directoryControllerFsm.fireTransition(event)
}

func (fsm *DirectoryControllerFiniteStateMachine) SendDataToRequester(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController, numInvAcks int32) {
	fsm.DirectoryController.TransferMessage(
		requester,
		fsm.DirectoryController.Cache.LineSize()+8,
		NewDataMessage(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.DirectoryController,
			numInvAcks,
		),
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) SendPutAckToRequester(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	fsm.DirectoryController.SendPutAckToRequester(
		producerFlow,
		tag,
		requester,
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) CopyDataToMem(tag uint32) {
	fsm.DirectoryController.Next().(*MemoryController).ReceiveMemWriteRequest(
		fsm.DirectoryController,
		tag,
		func() {},
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) SendFwdGetSToOwner(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	fsm.DirectoryController.TransferMessage(
		fsm.DirectoryEntry.Owner,
		8,
		NewFwdGetSMessage(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
			requester,
		),
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) SendFwdGetMToOwner(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	fsm.DirectoryController.TransferMessage(
		fsm.DirectoryEntry.Owner,
		8,
		NewFwdGetMMessage(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
			requester,
		),
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) SendInvToSharers(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	for _, sharer := range fsm.DirectoryEntry.Sharers {
		if requester != sharer {
			fsm.DirectoryController.TransferMessage(
				sharer,
				8,
				NewInvMessage(
					fsm.DirectoryController,
					producerFlow,
					producerFlow.Access(),
					tag,
					requester,
				),
			)
		}
	}
}

func (fsm *DirectoryControllerFiniteStateMachine) SendRecallToOwner(producerFlow CacheCoherenceFlow, tag uint32) {
	var owner = fsm.DirectoryEntry.Owner

	if owner.Cache.FindWay(tag) == INVALID_WAY {
		panic("Impossible")
	}

	fsm.DirectoryController.TransferMessage(
		owner,
		8,
		NewRecallMessage(
			fsm.DirectoryController,
			producerFlow,
			producerFlow.Access(),
			tag,
		),
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) SendRecallToSharers(producerFlow CacheCoherenceFlow, tag uint32) {
	for _, sharer := range fsm.DirectoryEntry.Sharers {
		if sharer.Cache.FindWay(tag) == INVALID_WAY {
			panic("Impossible")
		}

		fsm.DirectoryController.TransferMessage(
			sharer,
			8,
			NewRecallMessage(
				fsm.DirectoryController,
				producerFlow,
				producerFlow.Access(),
				tag,
			),
		)
	}
}

func (fsm *DirectoryControllerFiniteStateMachine) AddRequesterAndOwnerToSharers(requester *CacheController) {
	fsm.DirectoryEntry.Sharers = append(
		fsm.DirectoryEntry.Sharers,
		requester,
		fsm.DirectoryEntry.Owner,
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) AddRequesterToSharers(requester *CacheController) {
	fsm.DirectoryEntry.Sharers = append(
		fsm.DirectoryEntry.Sharers,
		requester,
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) RemoveRequesterFromSharers(requester *CacheController) {
	var sharersToPreserve []*CacheController

	for _, sharer := range fsm.DirectoryEntry.Sharers {
		if requester != sharer {
			sharersToPreserve = append(sharersToPreserve, sharer)
		}
	}

	fsm.DirectoryEntry.Sharers = sharersToPreserve
}

func (fsm *DirectoryControllerFiniteStateMachine) SetOwnerToRequester(requester *CacheController) {
	fsm.DirectoryEntry.Owner = requester
}

func (fsm *DirectoryControllerFiniteStateMachine) ClearSharers() {
	fsm.DirectoryEntry.Sharers = []*CacheController{}
}

func (fsm *DirectoryControllerFiniteStateMachine) ClearOwner() {
	fsm.DirectoryEntry.Owner = nil
}

func (fsm *DirectoryControllerFiniteStateMachine) FireServiceNonblockingRequestEvent(access *MemoryHierarchyAccess, tag uint32, hitInCache bool) {
	fsm.DirectoryController.MemoryHierarchy().Driver().BlockingEventDispatcher().Dispatch(
		NewGeneralCacheControllerServiceNonblockingRequestEvent(
			fsm.DirectoryController,
			access,
			tag,
			fsm.Set,
			fsm.Way,
			hitInCache,
		),
	)

	fsm.DirectoryController.UpdateStats(access.AccessType.IsWrite(), hitInCache)
}

func (fsm *DirectoryControllerFiniteStateMachine) FireCacheLineInsertEvent(access *MemoryHierarchyAccess, tag uint32, victimTag int32) {
	fsm.DirectoryController.MemoryHierarchy().Driver().BlockingEventDispatcher().Dispatch(
		NewLastLevelCacheControllerLineInsertEvent(
			fsm.DirectoryController,
			access,
			tag,
			fsm.Set,
			fsm.Way,
			victimTag,
		),
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) FireReplacementEvent(access *MemoryHierarchyAccess, tag uint32) {
	fsm.DirectoryController.MemoryHierarchy().Driver().BlockingEventDispatcher().Dispatch(
		NewGeneralCacheControllerLineReplacementEvent(
			fsm.DirectoryController,
			access,
			tag,
			fsm.Set,
			fsm.Way,
		),
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) FirePutSOrPutMAndDataFromOwnerEvent(access *MemoryHierarchyAccess, tag uint32) {
	fsm.DirectoryController.MemoryHierarchy().Driver().BlockingEventDispatcher().Dispatch(
		NewGeneralCacheControllerLastPutSOrPutMAndDataFromOwnerEvent(
			fsm.DirectoryController,
			access,
			tag,
			fsm.Set,
			fsm.Way,
		),
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) FireNonblockingRequestHitToTransientTagEvent(access *MemoryHierarchyAccess, tag uint32) {
	fsm.DirectoryController.MemoryHierarchy().Driver().BlockingEventDispatcher().Dispatch(
		NewGeneralCacheControllerNonblockingRequestHitToTransientTagEvent(
			fsm.DirectoryController,
			access,
			tag,
			fsm.Set,
			fsm.Way,
		),
	)
}

func (fsm *DirectoryControllerFiniteStateMachine) Hit(access *MemoryHierarchyAccess, tag uint32, set uint32, way uint32) {
	fsm.FireServiceNonblockingRequestEvent(access, tag, true)
	fsm.DirectoryController.Cache.ReplacementPolicy.HandlePromotionOnHit(access, set, way)
	fsm.Line().Access = access
}

func (fsm *DirectoryControllerFiniteStateMachine) Stall(action func()) {
	fsm.StalledEvents = append(fsm.StalledEvents, action)
}

func (fsm *DirectoryControllerFiniteStateMachine) StallEvent(event DirectoryControllerEvent) {
	fsm.Stall(func() {
		fsm.fireTransition(event)
	})
}

type DirectoryControllerFiniteStateMachineFactory struct {
	*simutil.FiniteStateMachineFactory
}

func NewDirectoryControllerFiniteStateMachineFactory() *DirectoryControllerFiniteStateMachineFactory {
	var fsmFactory = &DirectoryControllerFiniteStateMachineFactory{
		FiniteStateMachineFactory: simutil.NewFiniteStateMachineFactory(),
	}

	var actionWhenStateChanged = func(fsm simutil.FiniteStateMachine) {
		var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)

		if directoryControllerFsm.PreviousState() != directoryControllerFsm.State() {
			if directoryControllerFsm.State().(DirectoryControllerState).Stable() {
				var onCompletedCallback = directoryControllerFsm.OnCompletedCallback
				if onCompletedCallback != nil {
					directoryControllerFsm.OnCompletedCallback = nil
					onCompletedCallback()
				}
			}

			var stalledEventsToProcess []func()

			for _, stalledEvent := range directoryControllerFsm.StalledEvents {
				stalledEventsToProcess = append(stalledEventsToProcess, stalledEvent)
			}

			directoryControllerFsm.StalledEvents = []func(){}

			for _, stalledEventToProcess := range stalledEventsToProcess {
				stalledEventToProcess()
			}
		}
	}

	fsmFactory.InState(DirectoryControllerState_I).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		DirectoryControllerEventType_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetSEvent)

			directoryControllerFsm.DirectoryController.NumPendingMemoryAccesses++

			directoryControllerFsm.DirectoryController.Transfer(
				directoryControllerFsm.DirectoryController.Next(),
				8,
				func() {
					directoryControllerFsm.DirectoryController.Next().(*MemoryController).ReceiveMemReadRequest(
						directoryControllerFsm.DirectoryController,
						event.Tag(),
						func() {
							directoryControllerFsm.DirectoryController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
								func() {
									directoryControllerFsm.DirectoryController.NumPendingMemoryAccesses--

									var dataFromMemEvent = NewDataFromMemEvent(
										directoryControllerFsm.DirectoryController,
										event,
										event.Access(),
										event.Tag(),
										event.Requester,
									)

									directoryControllerFsm.fireTransition(dataFromMemEvent)
								},
								int(directoryControllerFsm.DirectoryController.HitLatency()),
							)
						},
					)
				},
			)

			directoryControllerFsm.FireServiceNonblockingRequestEvent(event.Access(), event.Tag(), false)
			directoryControllerFsm.Line().Access = event.Access()
			directoryControllerFsm.Line().Tag = int32(event.Tag())
		},
		DirectoryControllerState_IS_D,
	).OnCondition(
		DirectoryControllerEventType_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetMEvent)

			directoryControllerFsm.DirectoryController.NumPendingMemoryAccesses++

			directoryControllerFsm.DirectoryController.Transfer(
				directoryControllerFsm.DirectoryController.Next(),
				8,
				func() {
					directoryControllerFsm.DirectoryController.Next().(*MemoryController).ReceiveMemReadRequest(
						directoryControllerFsm.DirectoryController,
						event.Tag(),
						func() {
							directoryControllerFsm.DirectoryController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
								func() {
									directoryControllerFsm.DirectoryController.NumPendingMemoryAccesses--

									var dataFromMemEvent = NewDataFromMemEvent(
										directoryControllerFsm.DirectoryController,
										event,
										event.Access(),
										event.Tag(),
										event.Requester,
									)

									directoryControllerFsm.fireTransition(dataFromMemEvent)
								},
								int(directoryControllerFsm.DirectoryController.HitLatency()),
							)
						},
					)
				},
			)

			directoryControllerFsm.FireServiceNonblockingRequestEvent(event.Access(), event.Tag(), false)
			directoryControllerFsm.Line().Access = event.Access()
			directoryControllerFsm.Line().Tag = int32(event.Tag())
		},
		DirectoryControllerState_IM_D,
	)

	fsmFactory.InState(DirectoryControllerState_IS_D).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		DirectoryControllerEventType_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetSEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
			directoryControllerFsm.FireNonblockingRequestHitToTransientTagEvent(
				event.Access(),
				event.Tag(),
			)
		},
		DirectoryControllerState_IS_D,
	).OnCondition(
		DirectoryControllerEventType_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetMEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
			directoryControllerFsm.FireNonblockingRequestHitToTransientTagEvent(
				event.Access(),
				event.Tag(),
			)
		},
		DirectoryControllerState_IS_D,
	).OnCondition(
		DirectoryControllerEventType_DIR_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DirReplacementEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_IS_D,
	).OnCondition(
		DirectoryControllerEventType_PUTS_NOT_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSNotLastEvent)

			directoryControllerFsm.StallEvent(event)
		},
		DirectoryControllerState_IS_D,
	).OnCondition(
		DirectoryControllerEventType_PUTS_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSLastEvent)

			directoryControllerFsm.StallEvent(event)
		},
		DirectoryControllerState_IS_D,
	).OnCondition(
		DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutMAndDataFromNonOwnerEvent)

			directoryControllerFsm.StallEvent(event)
		},
		DirectoryControllerState_IS_D,
	).OnCondition(
		DirectoryControllerEventType_DATA_FROM_MEM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DataFromMemEvent)

			directoryControllerFsm.SendDataToRequester(event, event.Tag(), event.Requester, 0)
			directoryControllerFsm.AddRequesterToSharers(event.Requester)
			directoryControllerFsm.FireCacheLineInsertEvent(event.Access(), event.Tag(), directoryControllerFsm.VictimTag)
			directoryControllerFsm.EvicterTag = INVALID_TAG
			directoryControllerFsm.VictimTag = INVALID_TAG
			directoryControllerFsm.DirectoryController.Cache.ReplacementPolicy.HandleInsertionOnMiss(
				event.Access(),
				directoryControllerFsm.Set,
				directoryControllerFsm.Way,
			)
		},
		DirectoryControllerState_S,
	)

	fsmFactory.InState(DirectoryControllerState_IM_D).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		DirectoryControllerEventType_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetSEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
			directoryControllerFsm.FireNonblockingRequestHitToTransientTagEvent(event.Access(), event.Tag())
		},
		DirectoryControllerState_IM_D,
	).OnCondition(
		DirectoryControllerEventType_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetMEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
			directoryControllerFsm.FireNonblockingRequestHitToTransientTagEvent(event.Access(), event.Tag())
		},
		DirectoryControllerState_IM_D,
	).OnCondition(
		DirectoryControllerEventType_DIR_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DirReplacementEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_IM_D,
	).OnCondition(
		DirectoryControllerEventType_PUTS_NOT_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSNotLastEvent)

			directoryControllerFsm.StallEvent(event)
		},
		DirectoryControllerState_IM_D,
	).OnCondition(
		DirectoryControllerEventType_PUTS_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSLastEvent)

			directoryControllerFsm.StallEvent(event)
		},
		DirectoryControllerState_IM_D,
	).OnCondition(
		DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutMAndDataFromNonOwnerEvent)

			directoryControllerFsm.StallEvent(event)
		},
		DirectoryControllerState_IM_D,
	).OnCondition(
		DirectoryControllerEventType_DATA_FROM_MEM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DataFromMemEvent)

			directoryControllerFsm.SendDataToRequester(event, event.Tag(), event.Requester, 0)
			directoryControllerFsm.SetOwnerToRequester(event.Requester)
			directoryControllerFsm.FireCacheLineInsertEvent(event.Access(), event.Tag(), directoryControllerFsm.VictimTag)
			directoryControllerFsm.EvicterTag = INVALID_TAG
			directoryControllerFsm.VictimTag = INVALID_TAG
			directoryControllerFsm.DirectoryController.Cache.ReplacementPolicy.HandleInsertionOnMiss(
				event.Access(),
				directoryControllerFsm.Set,
				directoryControllerFsm.Way,
			)
		},
		DirectoryControllerState_M,
	)

	fsmFactory.InState(DirectoryControllerState_S).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		DirectoryControllerEventType_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetSEvent)

			directoryControllerFsm.SendDataToRequester(event, event.Tag(), event.Requester, 0)
			directoryControllerFsm.AddRequesterToSharers(event.Requester)
			directoryControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
		},
		DirectoryControllerState_S,
	).OnCondition(
		DirectoryControllerEventType_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetMEvent)

			var numInvAcks = int32(0)

			for _, sharer := range directoryControllerFsm.DirectoryEntry.Sharers {
				if sharer != event.Requester {
					numInvAcks++
				}
			}

			directoryControllerFsm.SendDataToRequester(event, event.Tag(), event.Requester, numInvAcks)

			directoryControllerFsm.SendInvToSharers(event, event.Tag(), event.Requester)
			directoryControllerFsm.ClearSharers()
			directoryControllerFsm.SetOwnerToRequester(event.Requester)
			directoryControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
		},
		DirectoryControllerState_M,
	).OnCondition(
		DirectoryControllerEventType_DIR_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DirReplacementEvent)

			directoryControllerFsm.NumRecallAcks = int32(len(directoryControllerFsm.DirectoryEntry.Sharers))
			directoryControllerFsm.SendRecallToSharers(event, uint32(directoryControllerFsm.Line().Tag))
			directoryControllerFsm.ClearSharers()
			directoryControllerFsm.OnCompletedCallback = event.OnCompletedCallback
			directoryControllerFsm.FireReplacementEvent(event.Access(), event.Tag())
			directoryControllerFsm.EvicterTag = int32(event.Tag())
			directoryControllerFsm.VictimTag = directoryControllerFsm.Line().Tag
			directoryControllerFsm.DirectoryController.NumEvictions++
		},
		DirectoryControllerState_SI_A,
	).OnCondition(
		DirectoryControllerEventType_PUTS_NOT_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSNotLastEvent)

			directoryControllerFsm.RemoveRequesterFromSharers(event.Requester)
			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_S,
	).OnCondition(
		DirectoryControllerEventType_PUTS_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSLastEvent)

			directoryControllerFsm.RemoveRequesterFromSharers(event.Requester)
			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
			directoryControllerFsm.FirePutSOrPutMAndDataFromOwnerEvent(event.Access(), event.Tag())
			directoryControllerFsm.Line().Access = nil
			directoryControllerFsm.Line().Tag = INVALID_TAG
		},
		DirectoryControllerState_I,
	).OnCondition(
		DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutMAndDataFromNonOwnerEvent)

			directoryControllerFsm.RemoveRequesterFromSharers(event.Requester)
			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_S,
	)

	fsmFactory.InState(DirectoryControllerState_M).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		DirectoryControllerEventType_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetSEvent)

			directoryControllerFsm.SendFwdGetSToOwner(event, event.Tag(), event.Requester)
			directoryControllerFsm.AddRequesterAndOwnerToSharers(event.Requester)
			directoryControllerFsm.ClearOwner()
			directoryControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
		},
		DirectoryControllerState_S_D,
	).OnCondition(
		DirectoryControllerEventType_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetMEvent)

			directoryControllerFsm.SendFwdGetMToOwner(event, event.Tag(), event.Requester)
			directoryControllerFsm.SetOwnerToRequester(event.Requester)
			directoryControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
		},
		DirectoryControllerState_M,
	).OnCondition(
		DirectoryControllerEventType_DIR_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DirReplacementEvent)

			directoryControllerFsm.NumRecallAcks = 1
			directoryControllerFsm.SendRecallToOwner(event, uint32(directoryControllerFsm.Line().Tag))
			directoryControllerFsm.ClearOwner()
			directoryControllerFsm.OnCompletedCallback = event.OnCompletedCallback
			directoryControllerFsm.FireReplacementEvent(event.Access(), event.Tag())
			directoryControllerFsm.EvicterTag = int32(event.Tag())
			directoryControllerFsm.VictimTag = directoryControllerFsm.Line().Tag
			directoryControllerFsm.DirectoryController.NumEvictions++
		},
		DirectoryControllerState_MI_A,
	).OnCondition(
		DirectoryControllerEventType_PUTS_NOT_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSNotLastEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_M,
	).OnCondition(
		DirectoryControllerEventType_PUTS_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSLastEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_M,
	).OnCondition(
		DirectoryControllerEventType_PUTM_AND_DATA_FROM_OWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutMAndDataFromOwnerEvent)

			directoryControllerFsm.CopyDataToMem(event.Tag())
			directoryControllerFsm.ClearOwner()
			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
			directoryControllerFsm.FirePutSOrPutMAndDataFromOwnerEvent(event.Access(), event.Tag())
			directoryControllerFsm.Line().Access = nil
			directoryControllerFsm.Line().Tag = INVALID_TAG
		},
		DirectoryControllerState_I,
	).OnCondition(
		DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutMAndDataFromNonOwnerEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_M,
	)

	fsmFactory.InState(DirectoryControllerState_S_D).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		DirectoryControllerEventType_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetSEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_S_D,
	).OnCondition(
		DirectoryControllerEventType_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetMEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_S_D,
	).OnCondition(
		DirectoryControllerEventType_DIR_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DirReplacementEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_S_D,
	).OnCondition(
		DirectoryControllerEventType_PUTS_NOT_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSNotLastEvent)

			directoryControllerFsm.RemoveRequesterFromSharers(event.Requester)
			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_S_D,
	).OnCondition(
		DirectoryControllerEventType_PUTS_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSLastEvent)

			directoryControllerFsm.RemoveRequesterFromSharers(event.Requester)
			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_S_D,
	).OnCondition(
		DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutMAndDataFromNonOwnerEvent)

			directoryControllerFsm.RemoveRequesterFromSharers(event.Requester)
			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_S_D,
	).OnCondition(
		DirectoryControllerEventType_DATA,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DataEvent)

			directoryControllerFsm.CopyDataToMem(event.Tag())
		},
		DirectoryControllerState_S,
	)

	fsmFactory.InState(DirectoryControllerState_MI_A).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		DirectoryControllerEventType_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetSEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_MI_A,
	).OnCondition(
		DirectoryControllerEventType_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetMEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_MI_A,
	).OnCondition(
		DirectoryControllerEventType_DIR_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DirReplacementEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_MI_A,
	).OnCondition(
		DirectoryControllerEventType_RECALL_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)

			directoryControllerFsm.NumRecallAcks--
		},
		DirectoryControllerState_MI_A,
	).OnCondition(
		DirectoryControllerEventType_LAST_RECALL_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*LastRecallAckEvent)

			directoryControllerFsm.CopyDataToMem(event.Tag())
			directoryControllerFsm.Line().Access = nil
			directoryControllerFsm.Line().Tag = INVALID_TAG
		},
		DirectoryControllerState_I,
	).OnCondition(
		DirectoryControllerEventType_PUTS_NOT_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSNotLastEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_MI_A,
	).OnCondition(
		DirectoryControllerEventType_PUTS_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSLastEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_MI_A,
	).OnCondition(
		DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutMAndDataFromNonOwnerEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_MI_A,
	)

	fsmFactory.InState(DirectoryControllerState_SI_A).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		DirectoryControllerEventType_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetSEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_SI_A,
	).OnCondition(
		DirectoryControllerEventType_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*GetMEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_SI_A,
	).OnCondition(
		DirectoryControllerEventType_DIR_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*DirReplacementEvent)

			directoryControllerFsm.Stall(event.OnStalledCallback)
		},
		DirectoryControllerState_SI_A,
	).OnCondition(
		DirectoryControllerEventType_RECALL_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)

			directoryControllerFsm.NumRecallAcks--
		},
		DirectoryControllerState_SI_A,
	).OnCondition(
		DirectoryControllerEventType_LAST_RECALL_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)

			directoryControllerFsm.Line().Access = nil
			directoryControllerFsm.Line().Tag = INVALID_TAG
		},
		DirectoryControllerState_I,
	).OnCondition(
		DirectoryControllerEventType_PUTS_NOT_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSNotLastEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_SI_A,
	).OnCondition(
		DirectoryControllerEventType_PUTS_LAST,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutSLastEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_SI_A,
	).OnCondition(
		DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var directoryControllerFsm = fsm.(*DirectoryControllerFiniteStateMachine)
			var event = params.(*PutMAndDataFromNonOwnerEvent)

			directoryControllerFsm.SendPutAckToRequester(event, event.Tag(), event.Requester)
		},
		DirectoryControllerState_SI_A,
	)

	return fsmFactory
}
