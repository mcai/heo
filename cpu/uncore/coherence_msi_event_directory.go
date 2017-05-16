package uncore

type DirectoryControllerEventType string

const (
	DirectoryControllerEventType_GETS = DirectoryControllerEventType("GETS")
	DirectoryControllerEventType_GETM = DirectoryControllerEventType("GETM")
	DirectoryControllerEventType_DIR_REPLACEMENT = DirectoryControllerEventType("DIR_REPLACEMENT")
	DirectoryControllerEventType_RECALL_ACK = DirectoryControllerEventType("RECALL_ACK")
	DirectoryControllerEventType_LAST_RECALL_ACK = DirectoryControllerEventType("LAST_RECALL_ACK")
	DirectoryControllerEventType_PUTS_NOT_LAST = DirectoryControllerEventType("PUTS_NOT_LAST")
	DirectoryControllerEventType_PUTS_LAST = DirectoryControllerEventType("PUTS_LAST")
	DirectoryControllerEventType_PUTM_AND_DATA_FROM_OWNER = DirectoryControllerEventType("PUTM_AND_DATA_FROM_OWNER")
	DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER = DirectoryControllerEventType("PUTM_AND_DATA_FROM_NONOWNER")
	DirectoryControllerEventType_DATA = DirectoryControllerEventType("DATA")
	DirectoryControllerEventType_DATA_FROM_MEM = DirectoryControllerEventType("DATA_FROM_MEM")
)

type DirectoryControllerEvent interface {
	ControllerEvent
	EventType() DirectoryControllerEventType
}

type BaseDirectoryControllerEvent struct {
	*BaseControllerEvent
	eventType DirectoryControllerEventType
}

func NewBaseDirectoryControllerEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, eventType DirectoryControllerEventType, access *MemoryHierarchyAccess, tag uint32) *BaseDirectoryControllerEvent {
	var event = &BaseDirectoryControllerEvent{
		BaseControllerEvent:NewBaseControllerEvent(generator, producerFlow, access, tag),
		eventType:eventType,
	}

	return event
}

func (event *BaseDirectoryControllerEvent) EventType() DirectoryControllerEventType {
	return event.eventType
}

type DataEvent struct {
	*BaseDirectoryControllerEvent
	Sender *CacheController
}

func NewDataEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender *CacheController) *DataEvent {
	var event = &DataEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_DATA, access, tag),
		Sender:sender,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type DataFromMemEvent struct {
	*BaseDirectoryControllerEvent
	Requester *CacheController
}

func NewDataFromMemEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *DataFromMemEvent {
	var event = &DataFromMemEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_DATA_FROM_MEM, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type GetMEvent struct {
	*BaseDirectoryControllerEvent
	Requester         *CacheController
	Set               uint32
	Way               uint32
	OnStalledCallback func()
}

func NewGetMEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController, set uint32, way uint32, onStalledCallback func()) *GetMEvent {
	var event = &GetMEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_GETM, access, tag),
		Requester:requester,
		Set:set,
		Way:way,
		OnStalledCallback:onStalledCallback,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type GetSEvent struct {
	*BaseDirectoryControllerEvent
	Requester         *CacheController
	Set               uint32
	Way               uint32
	OnStalledCallback func()
}

func NewGetSEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController, set uint32, way uint32, onStalledCallback func()) *GetSEvent {
	var event = &GetSEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_GETS, access, tag),
		Requester:requester,
		Set:set,
		Way:way,
		OnStalledCallback:onStalledCallback,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type LastRecallAckEvent struct {
	*BaseDirectoryControllerEvent
}

func NewLastRecallAckEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32) *LastRecallAckEvent {
	var event = &LastRecallAckEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_LAST_RECALL_ACK, access, tag),
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type PutMAndDataFromNonOwnerEvent struct {
	*BaseDirectoryControllerEvent
	Requester *CacheController
}

func NewPutMAndDataFromNonOwnerEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *PutMAndDataFromNonOwnerEvent {
	var event = &PutMAndDataFromNonOwnerEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_PUTM_AND_DATA_FROM_NONOWNER, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type PutMAndDataFromOwnerEvent struct {
	*BaseDirectoryControllerEvent
	Requester *CacheController
}

func NewPutMAndDataFromOwnerEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *PutMAndDataFromOwnerEvent {
	var event = &PutMAndDataFromOwnerEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_PUTM_AND_DATA_FROM_OWNER, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type PutSLastEvent struct {
	*BaseDirectoryControllerEvent
	Requester *CacheController
}

func NewPutSLastEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *PutSLastEvent {
	var event = &PutSLastEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_PUTS_LAST, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type PutSNotLastEvent struct {
	*BaseDirectoryControllerEvent
	Requester *CacheController
}

func NewPutSNotLastEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *PutSNotLastEvent {
	var event = &PutSNotLastEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_PUTS_NOT_LAST, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type RecallAckEvent struct {
	*BaseDirectoryControllerEvent
	Sender *CacheController
}

func NewRecallAckEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender *CacheController) *RecallAckEvent {
	var event = &RecallAckEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_RECALL_ACK, access, tag),
		Sender:sender,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}

type DirReplacementEvent struct {
	*BaseDirectoryControllerEvent
	CacheAccess         *CacheAccess
	Set                 uint32
	Way                 uint32
	OnCompletedCallback func()
	OnStalledCallback   func()
}

func NewDirReplacementEvent(generator *DirectoryController, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, cacheAccess *CacheAccess, set uint32, way uint32, onCompletedCallback func(), onStalledCallback func()) *DirReplacementEvent {
	var event = &DirReplacementEvent{
		BaseDirectoryControllerEvent:NewBaseDirectoryControllerEvent(generator, producerFlow, DirectoryControllerEventType_DIR_REPLACEMENT, access, tag),
		CacheAccess:cacheAccess,
		Set:set,
		Way:way,
		OnCompletedCallback:onCompletedCallback,
		OnStalledCallback:onStalledCallback,
	}

	SetupCacheCoherenceFlowTree(event)

	return event
}