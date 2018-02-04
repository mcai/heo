package uncore

import "github.com/mcai/heo/cpu/mem"

type CacheController struct {
	*BaseCacheController
	Cache                     *EvictableCache
	NumReadPorts              uint32
	NumWritePorts             uint32
	hitLatency                uint32
	PendingAccesses           map[uint32]*MemoryHierarchyAccess
	NumPendingAccessesPerType map[MemoryHierarchyAccessType]uint32
	FsmFactory                *CacheControllerFiniteStateMachineFactory
}

func NewCacheController(memoryHierarchy MemoryHierarchy, name string, deviceType MemoryDeviceType, geometry *mem.Geometry, replacementPolicyType CacheReplacementPolicyType, numReadPorts uint32, numWritePorts uint32, hitLatency uint32) *CacheController {
	var cacheController = &CacheController{
		NumReadPorts:              numReadPorts,
		NumWritePorts:             numWritePorts,
		hitLatency:                hitLatency,
		PendingAccesses:           make(map[uint32]*MemoryHierarchyAccess),
		NumPendingAccessesPerType: make(map[MemoryHierarchyAccessType]uint32),
	}

	cacheController.Cache = NewEvictableCache(
		geometry,
		func(set uint32, way uint32) CacheLineStateProvider {
			return NewCacheControllerFiniteStateMachine(set, way, cacheController)
		},
		replacementPolicyType,
	)

	cacheController.BaseCacheController = NewBaseCacheController(
		memoryHierarchy,
		name,
		deviceType,
	)

	cacheController.NumPendingAccessesPerType[MemoryHierarchyAccessType_IFETCH] = 0
	cacheController.NumPendingAccessesPerType[MemoryHierarchyAccessType_LOAD] = 0
	cacheController.NumPendingAccessesPerType[MemoryHierarchyAccessType_STORE] = 0

	cacheController.FsmFactory = NewCacheControllerFiniteStateMachineFactory()

	return cacheController
}

func (cacheController *CacheController) FindAccess(physicalTag uint32) *MemoryHierarchyAccess {
	if pendingAccess, ok := cacheController.PendingAccesses[physicalTag]; ok {
		return pendingAccess
	} else {
		return nil
	}
}

func (cacheController *CacheController) HitLatency() uint32 {
	return cacheController.hitLatency
}

func (cacheController *CacheController) CanAccess(accessType MemoryHierarchyAccessType, physicalTag uint32) bool {
	var access = cacheController.FindAccess(physicalTag)

	if access == nil {
		if accessType == MemoryHierarchyAccessType_STORE {
			return cacheController.NumPendingAccessesPerType[accessType] < cacheController.NumWritePorts
		} else {
			return cacheController.NumPendingAccessesPerType[accessType] < cacheController.NumReadPorts
		}
	} else {
		return accessType != MemoryHierarchyAccessType_STORE
	}
}

func (cacheController *CacheController) BeginAccess(accessType MemoryHierarchyAccessType, threadId int32, virtualPc int32, physicalAddress uint32, physicalTag uint32, onCompletedCallback func()) *MemoryHierarchyAccess {
	var newAccess = NewMemoryHierarchyAccess(cacheController.MemoryHierarchy(), accessType, threadId, virtualPc, physicalAddress, physicalTag, onCompletedCallback)

	var access = cacheController.FindAccess(physicalTag)

	if access != nil {
		access.Aliases = append([]*MemoryHierarchyAccess{newAccess}, access.Aliases...)
	} else {
		cacheController.PendingAccesses[physicalTag] = newAccess
		cacheController.NumPendingAccessesPerType[accessType]++
	}

	return newAccess
}

func (cacheController *CacheController) EndAccess(physicalTag uint32) {
	var access = cacheController.FindAccess(physicalTag)

	access.Complete()

	for _, alias := range access.Aliases {
		alias.Complete()
	}

	cacheController.NumPendingAccessesPerType[access.AccessType]--

	delete(cacheController.PendingAccesses, physicalTag)
}

func (cacheController *CacheController) ReceiveIfetch(access *MemoryHierarchyAccess, onCompletedCallback func()) {
	cacheController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
		func() {
			cacheController.OnLoad(access, access.PhysicalTag, onCompletedCallback)
		},
		int(cacheController.HitLatency()),
	)
}

func (cacheController *CacheController) ReceiveLoad(access *MemoryHierarchyAccess, onCompletedCallback func()) {
	cacheController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
		func() {
			cacheController.OnLoad(access, access.PhysicalTag, onCompletedCallback)
		},
		int(cacheController.HitLatency()),
	)
}

func (cacheController *CacheController) ReceiveStore(access *MemoryHierarchyAccess, onCompletedCallback func()) {
	cacheController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
		func() {
			cacheController.OnStore(access, access.PhysicalTag, onCompletedCallback)
		},
		int(cacheController.HitLatency()),
	)
}

func (cacheController *CacheController) SetNext(next MemoryDevice) {
	next.(*DirectoryController).CacheControllers = append(
		next.(*DirectoryController).CacheControllers,
		cacheController,
	)

	cacheController.BaseController.SetNext(next)
}

func (cacheController *CacheController) ReceiveMessage(message CoherenceMessage) {
	switch message.MessageType() {
	case CoherenceMessageType_FWD_GETS:
		cacheController.onFwdGetS(message.(*FwdGetSMessage))
	case CoherenceMessageType_FWD_GETM:
		cacheController.onFwdGetM(message.(*FwdGetMMessage))
	case CoherenceMessageType_INV:
		cacheController.onInv(message.(*InvMessage))
	case CoherenceMessageType_RECALL:
		cacheController.onRecall(message.(*RecallMessage))
	case CoherenceMessageType_PUT_ACK:
		cacheController.onPutAck(message.(*PutAckMessage))
	case CoherenceMessageType_DATA:
		cacheController.onData(message.(*DataMessage))
	case CoherenceMessageType_INV_ACK:
		cacheController.onInvAck(message.(*InvAckMessage))
	default:
		panic("Impossible")
	}
}

func (cacheController *CacheController) access(producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, onReplacementCompletedCallback func(uint32, uint32), onReplacementStalledCallback func()) {
	var set = cacheController.Cache.GetSet(tag)

	var cacheAccess = cacheController.Cache.NewAccess(access, tag)

	if cacheAccess.HitInCache {
		onReplacementCompletedCallback(set, cacheAccess.Way)
	} else {
		if cacheAccess.Replacement {
			var line = cacheController.Cache.Sets[set].Lines[cacheAccess.Way]
			var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
			cacheControllerFsm.OnEventReplacement(
				producerFlow,
				tag,
				cacheAccess,
				func() {
					onReplacementCompletedCallback(set, cacheAccess.Way)
				},
				func() {
					cacheController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
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

func (cacheController *CacheController) _onLoad(access *MemoryHierarchyAccess, tag uint32, loadFlow *LoadFlow) {
	var onStalledCallback = func() {
		cacheController._onLoad(access, tag, loadFlow)
	}

	cacheController.access(
		loadFlow,
		access,
		tag,
		func(set uint32, way uint32) {
			var line = cacheController.Cache.Sets[set].Lines[way]
			var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
			cacheControllerFsm.OnEventLoad(
				loadFlow,
				tag,
				loadFlow.OnCompletedCallback,
				onStalledCallback,
			)
		},
		onStalledCallback,
	)
}

func (cacheController *CacheController) OnLoad(access *MemoryHierarchyAccess, tag uint32, onCompletedCallback func()) {
	cacheController._onLoad(access, tag, NewLoadFlow(cacheController, access, tag, onCompletedCallback))
}

func (cacheController *CacheController) _onStore(access *MemoryHierarchyAccess, tag uint32, storeFlow *StoreFlow) {
	var onStalledCallback = func() {
		cacheController._onStore(access, tag, storeFlow)
	}

	cacheController.access(
		storeFlow,
		access,
		tag,
		func(set uint32, way uint32) {
			var line = cacheController.Cache.Sets[set].Lines[way]
			var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
			cacheControllerFsm.OnEventStore(
				storeFlow,
				tag,
				storeFlow.OnCompletedCallback,
				onStalledCallback,
			)
		},
		onStalledCallback,
	)
}

func (cacheController *CacheController) OnStore(access *MemoryHierarchyAccess, tag uint32, onCompletedCallback func()) {
	cacheController._onStore(access, tag, NewStoreFlow(cacheController, access, tag, onCompletedCallback))
}

func (cacheController *CacheController) onFwdGetS(message *FwdGetSMessage) {
	var way = cacheController.Cache.FindWay(message.Tag())
	var line = cacheController.Cache.Sets[cacheController.Cache.GetSet(message.Tag())].Lines[way]
	var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
	cacheControllerFsm.OnEventFwdGetS(message, message.Tag(), message.Requester)
}

func (cacheController *CacheController) onFwdGetM(message *FwdGetMMessage) {
	var way = cacheController.Cache.FindWay(message.Tag())
	var line = cacheController.Cache.Sets[cacheController.Cache.GetSet(message.Tag())].Lines[way]
	var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
	cacheControllerFsm.OnEventFwdGetM(message, message.Tag(), message.Requester)
}

func (cacheController *CacheController) onInv(message *InvMessage) {
	var way = cacheController.Cache.FindWay(message.Tag())
	var line = cacheController.Cache.Sets[cacheController.Cache.GetSet(message.Tag())].Lines[way]
	var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
	cacheControllerFsm.OnEventInv(message, message.Tag(), message.Requester)
}

func (cacheController *CacheController) onRecall(message *RecallMessage) {
	var way = cacheController.Cache.FindWay(message.Tag())
	var line = cacheController.Cache.Sets[cacheController.Cache.GetSet(message.Tag())].Lines[way]
	var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
	cacheControllerFsm.OnEventRecall(message, message.Tag())
}

func (cacheController *CacheController) onPutAck(message *PutAckMessage) {
	var way = cacheController.Cache.FindWay(message.Tag())
	var line = cacheController.Cache.Sets[cacheController.Cache.GetSet(message.Tag())].Lines[way]
	var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
	cacheControllerFsm.OnEventPutAck(message, message.Tag())
}

func (cacheController *CacheController) onData(message *DataMessage) {
	var way = cacheController.Cache.FindWay(message.Tag())
	var line = cacheController.Cache.Sets[cacheController.Cache.GetSet(message.Tag())].Lines[way]
	var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
	cacheControllerFsm.OnEventData(message, message.Tag(), message.Sender, message.NumInvAcks)
}

func (cacheController *CacheController) onInvAck(message *InvAckMessage) {
	var way = cacheController.Cache.FindWay(message.Tag())
	var line = cacheController.Cache.Sets[cacheController.Cache.GetSet(message.Tag())].Lines[way]
	var cacheControllerFsm = line.StateProvider.(*CacheControllerFiniteStateMachine)
	cacheControllerFsm.OnEventInvAck(message, message.Tag(), message.Sender)
}

type L1IController struct {
	*CacheController
}

func NewL1IController(memoryHierarchy MemoryHierarchy, name string) *L1IController {
	var l1IController = &L1IController{
		CacheController: NewCacheController(
			memoryHierarchy,
			name,
			MemoryDeviceType_L1I_CONTROLLER,
			mem.NewGeometry(
				memoryHierarchy.Config().L1ISize,
				memoryHierarchy.Config().L1IAssoc,
				memoryHierarchy.Config().L1ILineSize,
			),
			memoryHierarchy.Config().L1IReplacementPolicy,
			memoryHierarchy.Config().L1INumReadPorts,
			memoryHierarchy.Config().L1INumWritePorts,
			memoryHierarchy.Config().L1IHitLatency,
		),
	}

	return l1IController
}

type L1DController struct {
	*CacheController
}

func NewL1DController(memoryHierarchy MemoryHierarchy, name string) *L1DController {
	var l1DController = &L1DController{
		CacheController: NewCacheController(
			memoryHierarchy,
			name,
			MemoryDeviceType_L1D_CONTROLLER,
			mem.NewGeometry(
				memoryHierarchy.Config().L1DSize,
				memoryHierarchy.Config().L1DAssoc,
				memoryHierarchy.Config().L1DLineSize,
			),
			memoryHierarchy.Config().L1DReplacementPolicy,
			memoryHierarchy.Config().L1DNumReadPorts,
			memoryHierarchy.Config().L1DNumWritePorts,
			memoryHierarchy.Config().L1DHitLatency,
		),
	}

	return l1DController
}
