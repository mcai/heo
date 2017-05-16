package uncore

type CoherenceMessageType string

const (
	CoherenceMessageType_GETS = CoherenceMessageType("GETS")
	CoherenceMessageType_GETM = CoherenceMessageType("GETM")
	CoherenceMessageType_PUTS = CoherenceMessageType("PUTS")
	CoherenceMessageType_PUTM_AND_DATA = CoherenceMessageType("PUTM_AND_DATA")
	CoherenceMessageType_FWD_GETS = CoherenceMessageType("FWD_GETS")
	CoherenceMessageType_FWD_GETM = CoherenceMessageType("FWD_GETM")
	CoherenceMessageType_INV = CoherenceMessageType("INV")
	CoherenceMessageType_RECALL = CoherenceMessageType("RECALL")
	CoherenceMessageType_PUT_ACK = CoherenceMessageType("PUT_ACK")
	CoherenceMessageType_DATA = CoherenceMessageType("DATA")
	CoherenceMessageType_INV_ACK = CoherenceMessageType("INV_ACK")
	CoherenceMessageType_RECALL_ACK = CoherenceMessageType("RECALL_ACK")
)

type CoherenceMessage interface {
	CacheCoherenceFlow
	MessageType() CoherenceMessageType
	DestArrived() bool
	SetDestArrived(destArrived bool)
}

type BaseCoherenceMessage struct {
	*BaseCacheCoherenceFlow
	messageType CoherenceMessageType
	destArrived bool
}

func NewBaseCoherenceMessage(generator Controller, producerFlow CacheCoherenceFlow, messageType CoherenceMessageType, access *MemoryHierarchyAccess, tag uint32) *BaseCoherenceMessage {
	var message = &BaseCoherenceMessage{
		BaseCacheCoherenceFlow:NewBaseCacheCoherenceFlow(generator, producerFlow, access, tag),
		messageType:messageType,
	}

	return message
}

func (message *BaseCoherenceMessage) MessageType() CoherenceMessageType {
	return message.messageType
}

func (message *BaseCoherenceMessage) DestArrived() bool {
	return message.destArrived
}

func (message *BaseCoherenceMessage) SetDestArrived(destArrived bool) {
	message.destArrived = destArrived
}

type DataMessage struct {
	*BaseCoherenceMessage
	Sender     Controller
	NumInvAcks int32
}

func NewDataMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender Controller, numInvAcks int32) *DataMessage {
	var message = &DataMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_DATA, access, tag),
		Sender:sender,
		NumInvAcks:numInvAcks,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type FwdGetMMessage struct {
	*BaseCoherenceMessage
	Requester *CacheController
}

func NewFwdGetMMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *FwdGetMMessage {
	var message = &FwdGetMMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_FWD_GETM, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type FwdGetSMessage struct {
	*BaseCoherenceMessage
	Requester *CacheController
}

func NewFwdGetSMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *FwdGetSMessage {
	var message = &FwdGetSMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_FWD_GETS, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type GetMMessage struct {
	*BaseCoherenceMessage
	Requester *CacheController
}

func NewGetMMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *GetMMessage {
	var message = &GetMMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_GETM, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type GetSMessage struct {
	*BaseCoherenceMessage
	Requester *CacheController
}

func NewGetSMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *GetSMessage {
	var message = &GetSMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_GETS, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type InvAckMessage struct {
	*BaseCoherenceMessage
	Sender *CacheController
}

func NewInvAckMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender *CacheController) *InvAckMessage {
	var message = &InvAckMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_INV_ACK, access, tag),
		Sender:sender,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type InvMessage struct {
	*BaseCoherenceMessage
	Requester *CacheController
}

func NewInvMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *InvMessage {
	var message = &InvMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_INV, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type PutAckMessage struct {
	*BaseCoherenceMessage
}

func NewPutAckMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32) *PutAckMessage {
	var message = &PutAckMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_PUT_ACK, access, tag),
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type PutMAndDataMessage struct {
	*BaseCoherenceMessage
	Requester *CacheController
}

func NewPutMAndDataMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *PutMAndDataMessage {
	var message = &PutMAndDataMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_PUTM_AND_DATA, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type PutSMessage struct {
	*BaseCoherenceMessage
	Requester *CacheController
}

func NewPutSMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, requester *CacheController) *PutSMessage {
	var message = &PutSMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_PUTS, access, tag),
		Requester:requester,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type RecallAckMessage struct {
	*BaseCoherenceMessage
	Sender *CacheController
}

func NewRecallAckMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32, sender *CacheController) *RecallAckMessage {
	var message = &RecallAckMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_RECALL_ACK, access, tag),
		Sender:sender,
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}

type RecallMessage struct {
	*BaseCoherenceMessage
}

func NewRecallMessage(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32) *RecallMessage {
	var message = &RecallMessage{
		BaseCoherenceMessage:NewBaseCoherenceMessage(generator, producerFlow, CoherenceMessageType_RECALL, access, tag),
	}

	SetupCacheCoherenceFlowTree(message)

	return message
}