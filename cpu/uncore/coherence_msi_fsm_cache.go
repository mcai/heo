package uncore

import (
	"github.com/mcai/heo/simutil"
)

type CacheControllerFiniteStateMachine struct {
	*simutil.BaseFiniteStateMachine
	CacheController     *CacheController
	Set                 uint32
	Way                 uint32
	NumInvAcks          int32
	StalledEvents       []func()
	OnCompletedCallback func()
}

func NewCacheControllerFiniteStateMachine(set uint32, way uint32, cacheController *CacheController) *CacheControllerFiniteStateMachine {
	var fsm = &CacheControllerFiniteStateMachine{
		BaseFiniteStateMachine: simutil.NewBaseFiniteStateMachine(CacheControllerState_I),
		Set:                    set,
		Way:                    way,
		CacheController:        cacheController,
	}

	return fsm
}

func (fsm *CacheControllerFiniteStateMachine) Valid() bool {
	return fsm.State() != CacheControllerState_I
}

func (fsm *CacheControllerFiniteStateMachine) Line() *CacheLine {
	return fsm.CacheController.Cache.Sets[fsm.Set].Lines[fsm.Way]
}

func (fsm *CacheControllerFiniteStateMachine) fireTransition(event CacheControllerEvent) {
	event.Complete()
	fsm.CacheController.FsmFactory.FireTransition(fsm, event.EventType(), event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventLoad(producerFlow *LoadFlow, tag uint32, onCompletedCallback func(), onStalledCallback func()) {
	var event = NewLoadEvent(fsm.CacheController, producerFlow, producerFlow.Access(), tag, fsm.Set, fsm.Way, onCompletedCallback, onStalledCallback)
	fsm.fireTransition(event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventStore(producerFlow *StoreFlow, tag uint32, onCompletedCallback func(), onStalledCallback func()) {
	var event = NewStoreEvent(fsm.CacheController, producerFlow, producerFlow.Access(), tag, fsm.Set, fsm.Way, onCompletedCallback, onStalledCallback)
	fsm.fireTransition(event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventReplacement(producerFlow CacheCoherenceFlow, tag uint32, cacheAccess *CacheAccess, onCompletedCallback func(), onStalledCallback func()) {
	var event = NewReplacementEvent(fsm.CacheController, producerFlow, producerFlow.Access(), tag, cacheAccess, fsm.Set, fsm.Way, onCompletedCallback, onStalledCallback)
	fsm.fireTransition(event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventFwdGetS(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	var event = NewFwdGetSEvent(fsm.CacheController, producerFlow, producerFlow.Access(), tag, requester)
	fsm.fireTransition(event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventFwdGetM(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	var event = NewFwdGetMEvent(fsm.CacheController, producerFlow, producerFlow.Access(), tag, requester)
	fsm.fireTransition(event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventInv(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	var event = NewInvEvent(fsm.CacheController, producerFlow, producerFlow.Access(), tag, requester)
	fsm.fireTransition(event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventRecall(producerFlow CacheCoherenceFlow, tag uint32) {
	var event = NewRecallEvent(fsm.CacheController, producerFlow, producerFlow.Access(), tag)
	fsm.fireTransition(event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventPutAck(producerFlow CacheCoherenceFlow, tag uint32) {
	var event = NewPutAckEvent(fsm.CacheController, producerFlow, producerFlow.Access(), tag)
	fsm.fireTransition(event)
}

func (fsm *CacheControllerFiniteStateMachine) OnEventData(producerFlow CacheCoherenceFlow, tag uint32, sender Controller, numInvAcks int32) {
	fsm.NumInvAcks += numInvAcks

	switch sender.(type) {
	case *DirectoryController:
		if numInvAcks == 0 {
			var event = NewDataFromDirAcksEq0Event(
				fsm.CacheController,
				producerFlow,
				producerFlow.Access(),
				tag,
				sender,
			)

			fsm.fireTransition(event)

		} else {
			var event = NewDataFromDirAcksGt0Event(
				fsm.CacheController,
				producerFlow,
				producerFlow.Access(),
				tag,
				sender,
			)

			fsm.fireTransition(event)

			if fsm.NumInvAcks == 0 {
				fsm.OnEventLastInvAck(producerFlow, tag)
			}
		}
	default:
		var event = NewDataFromOwnerEvent(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			sender,
		)

		fsm.fireTransition(event)
	}
}

func (fsm *CacheControllerFiniteStateMachine) OnEventInvAck(producerFlow CacheCoherenceFlow, tag uint32, sender *CacheController) {
	var event = NewInvAckEvent(
		fsm.CacheController,
		producerFlow,
		producerFlow.Access(),
		tag,
		sender,
	)

	fsm.fireTransition(event)

	if fsm.NumInvAcks == 0 {
		fsm.OnEventLastInvAck(producerFlow, tag)
	}
}

func (fsm *CacheControllerFiniteStateMachine) OnEventLastInvAck(producerFlow CacheCoherenceFlow, tag uint32) {
	var event = NewLastInvAckEvent(
		fsm.CacheController,
		producerFlow,
		producerFlow.Access(),
		tag,
	)

	fsm.fireTransition(event)

	fsm.NumInvAcks = 0
}

func (fsm *CacheControllerFiniteStateMachine) SendGetSToDir(producerFlow CacheCoherenceFlow, tag uint32) {
	fsm.CacheController.TransferMessage(
		fsm.CacheController.Next().(*DirectoryController),
		8,
		NewGetSMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) SendGetMToDir(producerFlow CacheCoherenceFlow, tag uint32) {
	fsm.CacheController.TransferMessage(
		fsm.CacheController.Next().(*DirectoryController),
		8,
		NewGetMMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) SendPutSToDir(producerFlow CacheCoherenceFlow, tag uint32) {
	fsm.CacheController.TransferMessage(
		fsm.CacheController.Next().(*DirectoryController),
		8,
		NewPutSMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) SendPutMAndDataToDir(producerFlow CacheCoherenceFlow, tag uint32) {
	fsm.CacheController.TransferMessage(
		fsm.CacheController.Next().(*DirectoryController),
		fsm.CacheController.Cache.LineSize()+8,
		NewPutMAndDataMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) SendDataToRequesterAndDir(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	fsm.CacheController.TransferMessage(
		requester,
		10,
		NewDataMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
			0,
		),
	)

	fsm.CacheController.TransferMessage(
		fsm.CacheController.Next().(*DirectoryController),
		fsm.CacheController.Cache.LineSize()+8,
		NewDataMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
			0,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) SendDataToRequester(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	fsm.CacheController.TransferMessage(
		requester,
		fsm.CacheController.Cache.LineSize()+8,
		NewDataMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
			0,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) SendInvAckToRequester(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	fsm.CacheController.TransferMessage(
		requester,
		8,
		NewInvAckMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) SendRecallAckToDir(producerFlow CacheCoherenceFlow, tag uint32, size uint32) {
	fsm.CacheController.TransferMessage(
		fsm.CacheController.Next().(*DirectoryController),
		size,
		NewRecallAckMessage(
			fsm.CacheController,
			producerFlow,
			producerFlow.Access(),
			tag,
			fsm.CacheController,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) FireServiceNonblockingRequestEvent(access *MemoryHierarchyAccess, tag uint32, hitInCache bool) {
	fsm.CacheController.MemoryHierarchy().Driver().BlockingEventDispatcher().Dispatch(
		NewGeneralCacheControllerServiceNonblockingRequestEvent(
			fsm.CacheController,
			access,
			tag,
			fsm.Set,
			fsm.Way,
			hitInCache,
		),
	)

	fsm.CacheController.UpdateStats(access.AccessType.IsWrite(), hitInCache)
}

func (fsm *CacheControllerFiniteStateMachine) FireReplacementEvent(access *MemoryHierarchyAccess, tag uint32) {
	fsm.CacheController.MemoryHierarchy().Driver().BlockingEventDispatcher().Dispatch(
		NewGeneralCacheControllerLineReplacementEvent(
			fsm.CacheController,
			access,
			tag,
			fsm.Set,
			fsm.Way,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) FireNonblockingRequestHitToTransientTagEvent(access *MemoryHierarchyAccess, tag uint32) {
	fsm.CacheController.MemoryHierarchy().Driver().BlockingEventDispatcher().Dispatch(
		NewGeneralCacheControllerNonblockingRequestHitToTransientTagEvent(
			fsm.CacheController,
			access,
			tag,
			fsm.Set,
			fsm.Way,
		),
	)
}

func (fsm *CacheControllerFiniteStateMachine) Hit(access *MemoryHierarchyAccess, tag uint32, set uint32, way uint32) {
	fsm.FireServiceNonblockingRequestEvent(access, tag, true)
	fsm.CacheController.Cache.ReplacementPolicy.HandlePromotionOnHit(access, set, way)
	fsm.Line().Access = access
}

func (fsm *CacheControllerFiniteStateMachine) Stall(action func()) {
	fsm.StalledEvents = append(fsm.StalledEvents, action)
}

func (fsm *CacheControllerFiniteStateMachine) StallEvent(event CacheControllerEvent) {
	fsm.Stall(func() {
		fsm.fireTransition(event)
	})
}

type CacheControllerFiniteStateMachineFactory struct {
	*simutil.FiniteStateMachineFactory
}

func NewCacheControllerFiniteStateMachineFactory() *CacheControllerFiniteStateMachineFactory {
	var fsmFactory = &CacheControllerFiniteStateMachineFactory{
		FiniteStateMachineFactory: simutil.NewFiniteStateMachineFactory(),
	}

	var actionWhenStateChanged = func(fsm simutil.FiniteStateMachine) {
		var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)

		if cacheControllerFsm.PreviousState() != cacheControllerFsm.State() {
			if cacheControllerFsm.State().(CacheControllerState).Stable() {
				var onCompletedCallback = cacheControllerFsm.OnCompletedCallback
				if onCompletedCallback != nil {
					cacheControllerFsm.OnCompletedCallback = nil
					onCompletedCallback()
				}
			}

			var stalledEventsToProcess []func()

			for _, stalledEvent := range cacheControllerFsm.StalledEvents {
				stalledEventsToProcess = append(stalledEventsToProcess, stalledEvent)
			}

			cacheControllerFsm.StalledEvents = []func(){}

			for _, stalledEventToProcess := range stalledEventsToProcess {
				stalledEventToProcess()
			}
		}
	}

	fsmFactory.InState(CacheControllerState_I).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.SendGetSToDir(event, event.Tag())
			cacheControllerFsm.FireServiceNonblockingRequestEvent(event.Access(), event.Tag(), false)
			cacheControllerFsm.Line().Access = event.Access()
			cacheControllerFsm.Line().Tag = int32(event.Tag())
			cacheControllerFsm.OnCompletedCallback = func() {
				cacheControllerFsm.CacheController.Cache.ReplacementPolicy.HandleInsertionOnMiss(
					event.Access(),
					cacheControllerFsm.Set,
					cacheControllerFsm.Way,
				)
				event.OnCompletedCallback()
			}
		},
		CacheControllerState_IS_D,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.SendGetMToDir(event, event.Tag())
			cacheControllerFsm.FireServiceNonblockingRequestEvent(event.Access(), event.Tag(), false)
			cacheControllerFsm.Line().Access = event.Access()
			cacheControllerFsm.Line().Tag = int32(event.Tag())
			cacheControllerFsm.OnCompletedCallback = func() {
				cacheControllerFsm.CacheController.Cache.ReplacementPolicy.HandleInsertionOnMiss(
					event.Access(),
					cacheControllerFsm.Set,
					cacheControllerFsm.Way,
				)
				event.OnCompletedCallback()
			}
		},
		CacheControllerState_IM_AD,
	)

	fsmFactory.InState(CacheControllerState_IS_D).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
			cacheControllerFsm.FireNonblockingRequestHitToTransientTagEvent(event.Access(), event.Tag())
		},
		CacheControllerState_IS_D,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
			cacheControllerFsm.FireNonblockingRequestHitToTransientTagEvent(event.Access(), event.Tag())
		},
		CacheControllerState_IS_D,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var event = params.(*ReplacementEvent)

			event.OnStalledCallback()
		},
		CacheControllerState_IS_D,
	).OnCondition(
		CacheControllerEventType_INV,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*InvEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_IS_D,
	).OnCondition(
		CacheControllerEventType_DATA_FROM_DIR_ACKS_EQ_0,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_S,
	).OnCondition(
		CacheControllerEventType_DATA_FROM_OWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_S,
	)

	fsmFactory.InState(CacheControllerState_IM_AD).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
			cacheControllerFsm.FireNonblockingRequestHitToTransientTagEvent(event.Access(), event.Tag())
		},
		CacheControllerState_IM_AD,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
			cacheControllerFsm.FireNonblockingRequestHitToTransientTagEvent(event.Access(), event.Tag())
		},
		CacheControllerState_IM_AD,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var event = params.(*ReplacementEvent)
			event.OnStalledCallback()
		},
		CacheControllerState_IM_AD,
	).OnCondition(
		CacheControllerEventType_FWD_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetSEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_IM_AD,
	).OnCondition(
		CacheControllerEventType_FWD_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetMEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_IM_AD,
	).OnCondition(
		CacheControllerEventType_DATA_FROM_DIR_ACKS_EQ_0,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_M,
	).OnCondition(
		CacheControllerEventType_DATA_FROM_DIR_ACKS_GT_0,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_IM_A,
	).OnCondition(
		CacheControllerEventType_DATA_FROM_OWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_M,
	).OnCondition(
		CacheControllerEventType_INV_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)

			cacheControllerFsm.NumInvAcks--
		},
		CacheControllerState_IM_AD,
	)

	fsmFactory.InState(CacheControllerState_IM_A).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.FireNonblockingRequestHitToTransientTagEvent(event.Access(), event.Tag())
		},
		CacheControllerState_IM_A,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.FireNonblockingRequestHitToTransientTagEvent(event.Access(), event.Tag())
		},
		CacheControllerState_IM_A,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var event = params.(*ReplacementEvent)

			event.OnStalledCallback()
		},
		CacheControllerState_IM_A,
	).OnCondition(
		CacheControllerEventType_FWD_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetSEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_IM_A,
	).OnCondition(
		CacheControllerEventType_FWD_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetMEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_IM_A,
	).OnCondition(
		CacheControllerEventType_RECALL,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*RecallEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_IM_A,
	).OnCondition(
		CacheControllerEventType_INV_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			cacheControllerFsm.NumInvAcks--
		},
		CacheControllerState_IM_A,
	).OnCondition(
		CacheControllerEventType_LAST_INV_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_M,
	)

	fsmFactory.InState(CacheControllerState_S).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
			cacheControllerFsm.CacheController.memoryHierarchy.Driver().CycleAccurateEventQueue().Schedule(
				event.OnCompletedCallback,
				0,
			)
		},
		CacheControllerState_S,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.SendGetMToDir(event, event.Tag())
			cacheControllerFsm.OnCompletedCallback = event.OnCompletedCallback
			cacheControllerFsm.FireServiceNonblockingRequestEvent(event.Access(), event.Tag(), true)
		},
		CacheControllerState_SM_AD,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*ReplacementEvent)

			cacheControllerFsm.SendPutSToDir(event, uint32(cacheControllerFsm.Line().Tag))
			cacheControllerFsm.OnCompletedCallback = event.OnCompletedCallback
			cacheControllerFsm.FireReplacementEvent(event.Access(), event.Tag())
			cacheControllerFsm.CacheController.NumReplacements++
		},
		CacheControllerState_SI_A,
	).OnCondition(
		CacheControllerEventType_INV,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*InvEvent)

			cacheControllerFsm.SendInvAckToRequester(event, event.Tag(), event.Requester)
			cacheControllerFsm.Line().Access = nil
			cacheControllerFsm.Line().Tag = INVALID_TAG
		},
		CacheControllerState_I,
	).OnCondition(
		CacheControllerEventType_RECALL,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*RecallEvent)

			cacheControllerFsm.SendRecallAckToDir(event, event.Tag(), 8)
			cacheControllerFsm.Line().Access = nil
			cacheControllerFsm.Line().Tag = INVALID_TAG
		},
		CacheControllerState_I,
	)

	fsmFactory.InState(CacheControllerState_SM_AD).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
			cacheControllerFsm.CacheController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
				event.OnCompletedCallback,
				0,
			)
		},
		CacheControllerState_SM_AD,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
		},
		CacheControllerState_SM_AD,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var event = params.(*ReplacementEvent)

			event.OnStalledCallback()
		},
		CacheControllerState_SM_AD,
	).OnCondition(
		CacheControllerEventType_FWD_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetSEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_SM_AD,
	).OnCondition(
		CacheControllerEventType_FWD_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetMEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_SM_AD,
	).OnCondition(
		CacheControllerEventType_INV,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*InvEvent)

			cacheControllerFsm.SendInvAckToRequester(event, event.Tag(), event.Requester)
		},
		CacheControllerState_IM_AD,
	).OnCondition(
		CacheControllerEventType_RECALL,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*RecallEvent)

			cacheControllerFsm.SendRecallAckToDir(event, event.Tag(), 8)
		},
		CacheControllerState_IM_AD,
	).OnCondition(
		CacheControllerEventType_DATA_FROM_DIR_ACKS_EQ_0,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_M,
	).OnCondition(
		CacheControllerEventType_DATA_FROM_DIR_ACKS_GT_0,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_SM_A,
	).OnCondition(
		CacheControllerEventType_DATA_FROM_OWNER,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_M,
	).OnCondition(
		CacheControllerEventType_INV_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)

			cacheControllerFsm.NumInvAcks--
		},
		CacheControllerState_SM_AD,
	)

	fsmFactory.InState(CacheControllerState_SM_A).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
			cacheControllerFsm.CacheController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
				event.OnCompletedCallback,
				0,
			)
		},
		CacheControllerState_SM_A,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
		},
		CacheControllerState_SM_A,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var event = params.(*ReplacementEvent)

			event.OnStalledCallback()
		},
		CacheControllerState_SM_A,
	).OnCondition(
		CacheControllerEventType_FWD_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetSEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_SM_A,
	).OnCondition(
		CacheControllerEventType_FWD_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetMEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_SM_A,
	).OnCondition(
		CacheControllerEventType_INV_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)

			cacheControllerFsm.NumInvAcks--
		},
		CacheControllerState_SM_A,
	).OnCondition(
		CacheControllerEventType_LAST_INV_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
		},
		CacheControllerState_M,
	).OnCondition(
		CacheControllerEventType_RECALL,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*RecallEvent)

			cacheControllerFsm.StallEvent(event)
		},
		CacheControllerState_SM_A,
	)

	fsmFactory.InState(CacheControllerState_M).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
			cacheControllerFsm.CacheController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
				event.OnCompletedCallback,
				0,
			)
		},
		CacheControllerState_M,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.Hit(event.Access(), event.Tag(), event.Set, event.Way)
			cacheControllerFsm.CacheController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
				event.OnCompletedCallback,
				0,
			)
		},
		CacheControllerState_M,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*ReplacementEvent)

			cacheControllerFsm.SendPutMAndDataToDir(event, uint32(cacheControllerFsm.Line().Tag))
			cacheControllerFsm.OnCompletedCallback = event.OnCompletedCallback
			cacheControllerFsm.FireReplacementEvent(event.Access(), event.Tag())
			cacheControllerFsm.CacheController.NumReplacements++
		},
		CacheControllerState_MI_A,
	).OnCondition(
		CacheControllerEventType_FWD_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetSEvent)

			cacheControllerFsm.SendDataToRequesterAndDir(event, event.Tag(), event.Requester)
		},
		CacheControllerState_S,
	).OnCondition(
		CacheControllerEventType_FWD_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetMEvent)

			cacheControllerFsm.SendDataToRequester(event, event.Tag(), event.Requester)
			cacheControllerFsm.Line().Access = nil
			cacheControllerFsm.Line().Tag = INVALID_TAG
		},
		CacheControllerState_I,
	).OnCondition(
		CacheControllerEventType_RECALL,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*RecallEvent)

			cacheControllerFsm.SendRecallAckToDir(
				event,
				event.Tag(),
				cacheControllerFsm.CacheController.Cache.LineSize()+8,
			)
			cacheControllerFsm.Line().Access = nil
			cacheControllerFsm.Line().Tag = INVALID_TAG
		},
		CacheControllerState_I,
	)

	fsmFactory.InState(CacheControllerState_MI_A).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
		},
		CacheControllerState_MI_A,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
		},
		CacheControllerState_MI_A,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var event = params.(*ReplacementEvent)

			event.OnStalledCallback()
		},
		CacheControllerState_MI_A,
	).OnCondition(
		CacheControllerEventType_RECALL,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*RecallEvent)

			cacheControllerFsm.SendRecallAckToDir(event, event.Tag(), 8)
		},
		CacheControllerState_II_A,
	).OnCondition(
		CacheControllerEventType_FWD_GETS,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetSEvent)

			cacheControllerFsm.SendDataToRequesterAndDir(event, event.Tag(), event.Requester)
		},
		CacheControllerState_SI_A,
	).OnCondition(
		CacheControllerEventType_FWD_GETM,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*FwdGetMEvent)

			cacheControllerFsm.SendDataToRequester(event, event.Tag(), event.Requester)
		},
		CacheControllerState_II_A,
	).OnCondition(
		CacheControllerEventType_PUT_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)

			cacheControllerFsm.Line().Access = nil
			cacheControllerFsm.Line().Tag = INVALID_TAG
		},
		CacheControllerState_I,
	)

	fsmFactory.InState(CacheControllerState_SI_A).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
		},
		CacheControllerState_SI_A,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
		},
		CacheControllerState_SI_A,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var event = params.(*ReplacementEvent)

			event.OnStalledCallback()
		},
		CacheControllerState_SI_A,
	).OnCondition(
		CacheControllerEventType_INV,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*InvEvent)

			cacheControllerFsm.SendInvAckToRequester(event, event.Tag(), event.Requester)
		},
		CacheControllerState_II_A,
	).OnCondition(
		CacheControllerEventType_RECALL,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*RecallEvent)

			cacheControllerFsm.SendRecallAckToDir(event, event.Tag(), 8)
		},
		CacheControllerState_II_A,
	).OnCondition(
		CacheControllerEventType_PUT_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)

			cacheControllerFsm.Line().Access = nil
			cacheControllerFsm.Line().Tag = INVALID_TAG
		},
		CacheControllerState_I,
	)

	fsmFactory.InState(CacheControllerState_II_A).SetOnCompletedCallback(actionWhenStateChanged).OnCondition(
		CacheControllerEventType_LOAD,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*LoadEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
		},
		CacheControllerState_II_A,
	).OnCondition(
		CacheControllerEventType_STORE,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)
			var event = params.(*StoreEvent)

			cacheControllerFsm.Stall(event.OnStalledCallback)
		},
		CacheControllerState_II_A,
	).OnCondition(
		CacheControllerEventType_REPLACEMENT,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var event = params.(*ReplacementEvent)

			event.OnStalledCallback()
		},
		CacheControllerState_II_A,
	).OnCondition(
		CacheControllerEventType_PUT_ACK,
		func(fsm simutil.FiniteStateMachine, condition interface{}, params interface{}) {
			var cacheControllerFsm = fsm.(*CacheControllerFiniteStateMachine)

			cacheControllerFsm.Line().Access = nil
			cacheControllerFsm.Line().Tag = INVALID_TAG
		},
		CacheControllerState_I,
	)

	return fsmFactory
}
