package uncore

type MemoryHierarchyAccessType uint32

const (
	MemoryHierarchyAccessType_IFETCH = MemoryHierarchyAccessType(0)
	MemoryHierarchyAccessType_LOAD = MemoryHierarchyAccessType(1)
	MemoryHierarchyAccessType_STORE = MemoryHierarchyAccessType(2)
	MemoryHierarchyAccessType_UNKNOWN = MemoryHierarchyAccessType(3)
)

func (accessType MemoryHierarchyAccessType) IsRead() bool {
	return accessType == MemoryHierarchyAccessType_IFETCH ||
		accessType == MemoryHierarchyAccessType_LOAD
}

func (accessType MemoryHierarchyAccessType) IsWrite() bool {
	return accessType == MemoryHierarchyAccessType_STORE
}

type MemoryHierarchyAccess struct {
	MemoryHierarchy     MemoryHierarchy
	Id                  int32
	AccessType          MemoryHierarchyAccessType

	ThreadId            int32
	VirtualPc           int32
	PhysicalAddress     uint32
	PhysicalTag         uint32

	OnCompletedCallback func()

	Aliases             []*MemoryHierarchyAccess

	BeginCycle          int64
	EndCycle            int64
}

func NewMemoryHierarchyAccess(memoryHierarchy MemoryHierarchy, accessType MemoryHierarchyAccessType, threadId int32, virtualPc int32, physicalAddress uint32, physicalTag uint32, onCompletedCallback func()) *MemoryHierarchyAccess {
	var access = &MemoryHierarchyAccess{
		MemoryHierarchy:memoryHierarchy,
		Id:memoryHierarchy.CurrentMemoryHierarchyAccessId(),
		AccessType:accessType,
		ThreadId:threadId,
		VirtualPc:virtualPc,
		PhysicalAddress:physicalAddress,
		PhysicalTag:physicalTag,
		OnCompletedCallback:onCompletedCallback,
		BeginCycle:memoryHierarchy.Driver().CycleAccurateEventQueue().CurrentCycle,
	}

	memoryHierarchy.SetCurrentMemoryHierarchyAccessId(
		memoryHierarchy.CurrentMemoryHierarchyAccessId() + 1,
	)

	return access
}

func (access *MemoryHierarchyAccess) NumCycles() uint32 {
	return uint32(access.EndCycle - access.BeginCycle)
}

func (access *MemoryHierarchyAccess) Complete() {
	access.EndCycle = access.MemoryHierarchy.Driver().CycleAccurateEventQueue().CurrentCycle
	access.OnCompletedCallback()
	access.OnCompletedCallback = nil
}
