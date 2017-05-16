package uncore

import "github.com/mcai/heo/cpu/mem"

type DirectoryController struct {
	*BaseCacheController
	Cache                    *EvictableCache
	CacheControllers         []*CacheController
	NumPendingMemoryAccesses uint32
	FsmFactory               *DirectoryControllerFiniteStateMachineFactory
}

func NewDirectoryController(memoryHierarchy MemoryHierarchy, name string) *DirectoryController {
	var directoryController = &DirectoryController{
	}

	directoryController.Cache = NewEvictableCache(
		mem.NewGeometry(
			memoryHierarchy.Config().L2Size,
			memoryHierarchy.Config().L2Assoc,
			memoryHierarchy.Config().L2LineSize,
		),
		func(set uint32, way uint32) CacheLineStateProvider {
			return NewDirectoryControllerFiniteStateMachine(set, way, directoryController)
		},
		memoryHierarchy.Config().L2ReplacementPolicy,
	)

	directoryController.BaseCacheController = NewBaseCacheController(
		memoryHierarchy,
		name,
		MemoryDeviceType_L2_CONTROLLER,
	)

	directoryController.FsmFactory = NewDirectoryControllerFiniteStateMachineFactory()

	return directoryController
}

func (directoryController *DirectoryController) HitLatency() uint32 {
	return directoryController.MemoryHierarchy().Config().L2HitLatency
}

func (directoryController *DirectoryController) SendPutAckToRequester(producerFlow CacheCoherenceFlow, tag uint32, requester *CacheController) {
	directoryController.TransferMessage(requester, 8, NewPutAckMessage(directoryController, producerFlow, producerFlow.Access(), tag))
}

func (directoryController *DirectoryController) ReceiveMessage(message CoherenceMessage) {
	switch message.MessageType() {
	case CoherenceMessageType_GETS:
		directoryController.onGetS(message.(*GetSMessage))
	case CoherenceMessageType_GETM:
		directoryController.onGetM(message.(*GetMMessage))
	case CoherenceMessageType_RECALL_ACK:
		directoryController.onRecallAck(message.(*RecallAckMessage))
	case CoherenceMessageType_PUTS:
		directoryController.onPutS(message.(*PutSMessage))
	case CoherenceMessageType_PUTM_AND_DATA:
		directoryController.onPutMAndData(message.(*PutMAndDataMessage))
	case CoherenceMessageType_DATA:
		directoryController.onData(message.(*DataMessage))
	default:
		panic("Impossible")
	}
}

func (directoryController *DirectoryController) access(producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController, onReplacementCompletedCallback func(set uint32, way uint32), onReplacementStalledCallback func()) {
	var set = directoryController.Cache.GetSet(tag)

	for _, line := range directoryController.Cache.Sets[set].Lines {
		var directoryControllerFsm = line.StateProvider.(*DirectoryControllerFiniteStateMachine)

		if line.State() == DirectoryControllerState_MI_A ||
			line.State() == DirectoryControllerState_SI_A && directoryControllerFsm.EvicterTag == int32(tag) {
			directoryControllerFsm.Stall(onReplacementStalledCallback)
			return
		}
	}

	var cacheAccess = directoryController.Cache.NewAccess(access, tag)
	if cacheAccess.HitInCache {
		onReplacementCompletedCallback(set, cacheAccess.Way)
	} else {
		if cacheAccess.Replacement {
			var line = directoryController.Cache.Sets[set].Lines[cacheAccess.Way]
			var directoryControllerFsm = line.StateProvider.(*DirectoryControllerFiniteStateMachine)

			directoryControllerFsm.OnEventReplacement(
				producerFlow,
				tag,
				cacheAccess,
				requester,
				func() {
					onReplacementCompletedCallback(set, cacheAccess.Way)
				},
				func() {
					directoryController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
						onReplacementStalledCallback,
						1,
					)
				},
			)
		} else {
			onReplacementCompletedCallback(set, cacheAccess.Way)
		}
	}
}

func (directoryController *DirectoryController) onGetS(message *GetSMessage) {
	var onStalledCallback = func() {
		directoryController.onGetS(message)
	}

	directoryController.access(
		message,
		message.Access(),
		message.Tag(),
		message.Requester,
		func(set uint32, way uint32) {
			var line = directoryController.Cache.Sets[set].Lines[way]
			var directoryControllerFsm = line.StateProvider.(*DirectoryControllerFiniteStateMachine)

			directoryControllerFsm.OnEventGetS(message, message.Tag(), message.Requester, onStalledCallback)
		},
		onStalledCallback,
	)
}

func (directoryController *DirectoryController) onGetM(message *GetMMessage) {
	var onStalledCallback = func() {
		directoryController.onGetM(message)
	}

	directoryController.access(
		message,
		message.Access(),
		message.Tag(),
		message.Requester,
		func(set uint32, way uint32) {
			var line = directoryController.Cache.Sets[set].Lines[way]
			var directoryControllerFsm = line.StateProvider.(*DirectoryControllerFiniteStateMachine)

			directoryControllerFsm.OnEventGetM(message, message.Tag(), message.Requester, onStalledCallback)
		},
		onStalledCallback,
	)
}

func (directoryController *DirectoryController) onRecallAck(message *RecallAckMessage) {
	var sender = message.Sender
	var tag = message.Tag()

	var way = directoryController.Cache.FindWay(tag)
	var line = directoryController.Cache.Sets[directoryController.Cache.GetSet(tag)].Lines[way]
	var directoryControllerFsm = line.StateProvider.(*DirectoryControllerFiniteStateMachine)
	directoryControllerFsm.OnEventRecallAck(message, tag, sender)
}

func (directoryController *DirectoryController) onPutS(message *PutSMessage) {
	var requester = message.Requester
	var tag = message.Tag()

	var way = directoryController.Cache.FindWay(tag)

	if way == INVALID_WAY {
		directoryController.SendPutAckToRequester(message, tag, requester)
	} else {
		var line = directoryController.Cache.Sets[directoryController.Cache.GetSet(tag)].Lines[way]
		var directoryControllerFsm = line.StateProvider.(*DirectoryControllerFiniteStateMachine)
		directoryControllerFsm.OnEventPutS(message, tag, requester)
	}
}

func (directoryController *DirectoryController) onPutMAndData(message *PutMAndDataMessage) {
	var requester = message.Requester
	var tag = message.Tag()

	var way = directoryController.Cache.FindWay(tag)

	if way == INVALID_WAY {
		directoryController.SendPutAckToRequester(message, tag, requester)
	} else {
		var line = directoryController.Cache.Sets[directoryController.Cache.GetSet(tag)].Lines[way]
		var directoryControllerFsm = line.StateProvider.(*DirectoryControllerFiniteStateMachine)
		directoryControllerFsm.OnEventPutMAndData(message, tag, requester)
	}
}

func (directoryController *DirectoryController) onData(message *DataMessage) {
	var sender = message.Sender.(*CacheController)
	var tag = message.Tag()

	var way = directoryController.Cache.FindWay(tag)

	var line = directoryController.Cache.Sets[directoryController.Cache.GetSet(tag)].Lines[way]
	var directoryControllerFsm = line.StateProvider.(*DirectoryControllerFiniteStateMachine)
	directoryControllerFsm.OnEventData(message, tag, sender)
}