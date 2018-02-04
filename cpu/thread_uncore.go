package cpu

import "github.com/mcai/heo/cpu/uncore"

type MemoryHierarchyThread struct {
	*BaseThread
	FetchStalled         bool
	LastFetchedCacheLine int32
}

func NewMemoryHierarchyThread(core Core, num int32) *MemoryHierarchyThread {
	var memoryHierarchyThread = &MemoryHierarchyThread{
		BaseThread:           NewBaseThread(core, num),
		LastFetchedCacheLine: -1,
	}

	return memoryHierarchyThread
}

func (thread *MemoryHierarchyThread) Itlb() *uncore.TranslationLookasideBuffer {
	return thread.Core().Processor().Experiment.MemoryHierarchy.ITlbs()[thread.Id()]
}

func (thread *MemoryHierarchyThread) Dtlb() *uncore.TranslationLookasideBuffer {
	return thread.Core().Processor().Experiment.MemoryHierarchy.DTlbs()[thread.Id()]
}
