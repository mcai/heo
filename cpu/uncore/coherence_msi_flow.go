package uncore

import (
	"fmt"
)

type CacheCoherenceFlow interface {
	Id() int32
	Generator() Controller
	ProducerFlow() CacheCoherenceFlow
	AncestorFlow() CacheCoherenceFlow
	SetAncestorFlow(ancestorFlow CacheCoherenceFlow)
	ChildFlows() []CacheCoherenceFlow
	SetChildFlows(childFlows []CacheCoherenceFlow)
	NumPendingDescendantFlows() int32
	SetNumPendingDescendantFlows(numPendingDescendantFlows int32)
	BeginCycle() int64
	EndCycle() int64
	Completed() bool
	Access() *MemoryHierarchyAccess
	Tag() uint32
	Complete()
}

type BaseCacheCoherenceFlow struct {
	id                        int32
	generator                 Controller
	producerFlow              CacheCoherenceFlow
	ancestorFlow              CacheCoherenceFlow
	childFlows                []CacheCoherenceFlow
	numPendingDescendantFlows int32
	beginCycle                int64
	endCycle                  int64
	completed                 bool
	access                    *MemoryHierarchyAccess
	tag                       uint32
}

func NewBaseCacheCoherenceFlow(generator Controller, producerFlow CacheCoherenceFlow, access *MemoryHierarchyAccess, tag uint32) *BaseCacheCoherenceFlow {
	var flow = &BaseCacheCoherenceFlow{
		id:generator.MemoryHierarchy().CurrentCacheCoherenceFlowId(),
		generator:generator,
		producerFlow:producerFlow,
		access:access,
		tag:tag,
	}

	generator.MemoryHierarchy().SetCurrentCacheCoherenceFlowId(
		generator.MemoryHierarchy().CurrentCacheCoherenceFlowId() + 1,
	)

	flow.beginCycle = generator.MemoryHierarchy().Driver().CycleAccurateEventQueue().CurrentCycle

	return flow
}

func (flow *BaseCacheCoherenceFlow) Id() int32 {
	return flow.id
}

func (flow *BaseCacheCoherenceFlow) Generator() Controller {
	return flow.generator
}

func (flow *BaseCacheCoherenceFlow) ProducerFlow() CacheCoherenceFlow {
	return flow.producerFlow
}

func (flow *BaseCacheCoherenceFlow) AncestorFlow() CacheCoherenceFlow {
	return flow.ancestorFlow
}

func (flow *BaseCacheCoherenceFlow) SetAncestorFlow(ancestorFlow CacheCoherenceFlow) {
	flow.ancestorFlow = ancestorFlow
}

func (flow *BaseCacheCoherenceFlow) ChildFlows() []CacheCoherenceFlow {
	return flow.childFlows
}

func (flow *BaseCacheCoherenceFlow) SetChildFlows(childFlows []CacheCoherenceFlow) {
	flow.childFlows = childFlows
}

func (flow *BaseCacheCoherenceFlow) NumPendingDescendantFlows() int32 {
	return flow.numPendingDescendantFlows
}

func (flow *BaseCacheCoherenceFlow) SetNumPendingDescendantFlows(numPendingDescendantFlows int32) {
	flow.numPendingDescendantFlows = numPendingDescendantFlows
}

func (flow *BaseCacheCoherenceFlow) BeginCycle() int64 {
	return flow.beginCycle
}

func (flow *BaseCacheCoherenceFlow) EndCycle() int64 {
	return flow.endCycle
}

func (flow *BaseCacheCoherenceFlow) Completed() bool {
	return flow.completed
}

func (flow *BaseCacheCoherenceFlow) Access() *MemoryHierarchyAccess {
	return flow.access
}

func (flow *BaseCacheCoherenceFlow) Tag() uint32 {
	return flow.tag
}

func (flow *BaseCacheCoherenceFlow) Complete() {
	flow.completed = true
	flow.endCycle = flow.generator.MemoryHierarchy().Driver().CycleAccurateEventQueue().CurrentCycle
	flow.ancestorFlow.SetNumPendingDescendantFlows(
		flow.ancestorFlow.NumPendingDescendantFlows() - 1)

	if flow.ancestorFlow.NumPendingDescendantFlows() == 0 {
		var pendingFlowsToReserve []CacheCoherenceFlow

		for _, pendingFlow := range flow.generator.MemoryHierarchy().PendingFlows() {
			if pendingFlow != flow.ancestorFlow {
				pendingFlowsToReserve = append(pendingFlowsToReserve, pendingFlow)
			}
		}

		flow.generator.MemoryHierarchy().SetPendingFlows(pendingFlowsToReserve)
	}
}

type LoadFlow struct {
	*BaseCacheCoherenceFlow
	OnCompletedCallback func()
}

func NewLoadFlow(generator *CacheController, access *MemoryHierarchyAccess, tag uint32, onCompletedCallback func()) *LoadFlow {
	var flow = &LoadFlow{
		BaseCacheCoherenceFlow:NewBaseCacheCoherenceFlow(generator, nil, access, tag),
	}

	flow.OnCompletedCallback = func() {
		onCompletedCallback()
		flow.Complete()
	}

	SetupCacheCoherenceFlowTree(flow)

	return flow
}

func (flow *LoadFlow) String() string {
	return fmt.Sprintf(
		"[%d] %s: LoadFlow{id=%d, tag=0x%08x}",
		flow.BeginCycle(),
		flow.Generator(),
		flow.Id(),
		flow.Tag(),
	)
}

type StoreFlow struct {
	*BaseCacheCoherenceFlow
	OnCompletedCallback func()
}

func NewStoreFlow(generator *CacheController, access *MemoryHierarchyAccess, tag uint32, onCompletedCallback func()) *StoreFlow {
	var flow = &StoreFlow{
		BaseCacheCoherenceFlow:NewBaseCacheCoherenceFlow(generator, nil, access, tag),
	}

	flow.OnCompletedCallback = func() {
		onCompletedCallback()
		flow.Complete()
	}

	SetupCacheCoherenceFlowTree(flow)

	return flow
}

func (flow *StoreFlow) String() string {
	return fmt.Sprintf(
		"[%d] %s: StoreFlow{id=%d, tag=0x%08x}",
		flow.BeginCycle(),
		flow.Generator(),
		flow.Id(),
		flow.Tag(),
	)
}