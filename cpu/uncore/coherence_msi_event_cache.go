package uncore

type CacheControllerEventType string

const (
	CacheControllerEventType_LOAD = CacheControllerEventType("LOAD")
	CacheControllerEventType_STORE = CacheControllerEventType("STORE")
	CacheControllerEventType_REPLACEMENT = CacheControllerEventType("REPLACEMENT")
	CacheControllerEventType_FWD_GETS = CacheControllerEventType("FWD_GETS")
	CacheControllerEventType_FWD_GETM = CacheControllerEventType("FWD_GETM")
	CacheControllerEventType_INV = CacheControllerEventType("INV")
	CacheControllerEventType_RECALL = CacheControllerEventType("RECALL")
	CacheControllerEventType_PUT_ACK = CacheControllerEventType("PUT_ACK")
	CacheControllerEventType_DATA_FROM_DIR_ACKS_EQ_0 = CacheControllerEventType("DATA_FROM_DIR_ACKS_EQ_0")
	CacheControllerEventType_DATA_FROM_DIR_ACKS_GT_0 = CacheControllerEventType("DATA_FROM_DIR_ACKS_GT_0")
	CacheControllerEventType_DATA_FROM_OWNER = CacheControllerEventType("DATA_FROM_OWNER")
	CacheControllerEventType_INV_ACK = CacheControllerEventType("INV_ACK")
	CacheControllerEventType_LAST_INV_ACK = CacheControllerEventType("LAST_INV_ACK")
)

type CacheControllerEvent interface {
	ControllerEvent
	EventType() CacheControllerEventType
}

type BaseCacheControllerEvent struct {
	*BaseControllerEvent
	eventType CacheControllerEventType
}

func NewBaseCacheControllerEvent(generator *CacheController, producerFlow CacheCoherenceFlow, eventType CacheControllerEventType, access *MemoryHierarchyAccess, tag uint32) *BaseCacheControllerEvent {
	var event = &BaseCacheControllerEvent{
		BaseControllerEvent:NewBaseControllerEvent(generator, producerFlow, access, tag),
		eventType:eventType,
	}

	return event
}

func (event *BaseCacheControllerEvent) EventType() CacheControllerEventType {
	return event.eventType
}

type DataFromDirAcksEq0Event struct {
	*BaseCacheControllerEvent
	Sender Controller
}

func NewDataFromDirAcksEq0Event(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender Controller) *DataFromDirAcksEq0Event {
	var event = &DataFromDirAcksEq0Event{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_DATA_FROM_DIR_ACKS_EQ_0, access, tag),
		Sender:sender,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type DataFromDirAcksGt0Event struct {
	*BaseCacheControllerEvent
	Sender Controller
}

func NewDataFromDirAcksGt0Event(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender Controller) *DataFromDirAcksGt0Event {
	var event = &DataFromDirAcksGt0Event{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_DATA_FROM_DIR_ACKS_GT_0, access, tag),
		Sender:sender,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type DataFromOwnerEvent struct {
	*BaseCacheControllerEvent
	Sender Controller
}

func NewDataFromOwnerEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender Controller) *DataFromOwnerEvent {
	var event = &DataFromOwnerEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_DATA_FROM_OWNER, access, tag),
		Sender:sender,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type FwdGetMEvent struct {
	*BaseCacheControllerEvent
	Requester *CacheController
}

func NewFwdGetMEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *FwdGetMEvent {
	var event = &FwdGetMEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_FWD_GETM, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type FwdGetSEvent struct {
	*BaseCacheControllerEvent
	Requester *CacheController
}

func NewFwdGetSEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *FwdGetSEvent {
	var event = &FwdGetSEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_FWD_GETS, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type InvAckEvent struct {
	*BaseCacheControllerEvent
	Sender *CacheController
}

func NewInvAckEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender *CacheController) *InvAckEvent {
	var event = &InvAckEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_INV_ACK, access, tag),
		Sender:sender,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type InvEvent struct {
	*BaseCacheControllerEvent
	Requester *CacheController
}

func NewInvEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *InvEvent {
	var event = &InvEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_INV, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type LastInvAckEvent struct {
	*BaseCacheControllerEvent
}

func NewLastInvAckEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32) *LastInvAckEvent {
	var event = &LastInvAckEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_LAST_INV_ACK, access, tag),
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type LoadEvent struct {
	*BaseCacheControllerEvent
	Set                 uint32
	Way                 uint32
	OnCompletedCallback func()
	OnStalledCallback   func()
}

func NewLoadEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, set uint32, way uint32, onCompletedCallback func(), onStalledCallback func()) *LoadEvent {
	var event = &LoadEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_LOAD, access, tag),
		Set:set,
		Way:way,
		OnCompletedCallback:onCompletedCallback,
		OnStalledCallback:onStalledCallback,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type PutAckEvent struct {
	*BaseCacheControllerEvent
}

func NewPutAckEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32) *PutAckEvent {
	var event = &PutAckEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_PUT_ACK, access, tag),
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type RecallEvent struct {
	*BaseCacheControllerEvent
}

func NewRecallEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32) *RecallEvent {
	var event = &RecallEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_RECALL, access, tag),
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type ReplacementEvent struct {
	*BaseCacheControllerEvent
	CacheAccess         *CacheAccess
	Set                 uint32
	Way                 uint32
	OnCompletedCallback func()
	OnStalledCallback   func()
}

func NewReplacementEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, cacheAccess *CacheAccess, set uint32, way uint32, onCompletedCallback func(), onStalledCallback func()) *ReplacementEvent {
	var event = &ReplacementEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_REPLACEMENT, access, tag),
		CacheAccess:cacheAccess,
		Set:set,
		Way:way,
		OnCompletedCallback:onCompletedCallback,
		OnStalledCallback:onStalledCallback,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type StoreEvent struct {
	*BaseCacheControllerEvent
	Set                 uint32
	Way                 uint32
	OnCompletedCallback func()
	OnStalledCallback   func()
}

func NewStoreEvent(generator *CacheController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, set uint32, way uint32, onCompletedCallback func(), onStalledCallback func()) *StoreEvent {
	var event = &StoreEvent{
		BaseCacheControllerEvent:NewBaseCacheControllerEvent(generator, producerFlow, CacheControllerEventType_STORE, access, tag),
		Set:set,
		Way:way,
		OnCompletedCallback:onCompletedCallback,
		OnStalledCallback:onStalledCallback,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}