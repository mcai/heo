package uncore

type ControllerEvent interface {
	CacheCoherenceFlow
}

type BaseControllerEvent struct {
	*BaseCacheCoherenceFlow
}

func NewBaseControllerEvent(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32) *BaseControllerEvent {
	var event = &BaseControllerEvent{
		BaseCacheCoherenceFlow:NewBaseCacheCoherenceFlow(generator, producerFlow, access, tag),
	}

	return event
}