package uncore

import "github.com/mcai/heo/cpu/mem"

type TranslationLookasideBuffer struct {
	MemoryHierarchy MemoryHierarchy
	Name            string
	Cache           *EvictableCache
	NumHits         int64
	NumMisses       int64
	NumEvictions    int64
}

func NewTranslationLookasideBuffer(memoryHierarchy MemoryHierarchy, name string) *TranslationLookasideBuffer {
	var tlb = &TranslationLookasideBuffer{
		MemoryHierarchy: memoryHierarchy,
		Name:            name,
		Cache: NewEvictableCache(
			mem.NewGeometry(
				memoryHierarchy.Config().TlbSize,
				memoryHierarchy.Config().TlbAssoc,
				memoryHierarchy.Config().TlbLineSize,
			),
			func(set uint32, way uint32) CacheLineStateProvider {
				return NewBaseCacheLineStateProvider(
					false,
					func(state interface{}) bool {
						return state != false
					},
				)
			},
			CacheReplacementPolicyType_LRU,
		),
	}

	return tlb
}

func (tlb *TranslationLookasideBuffer) NumAccesses() int64 {
	return tlb.NumHits + tlb.NumMisses
}

func (tlb *TranslationLookasideBuffer) HitRatio() float64 {
	if tlb.NumAccesses() == 0 {
		return 0
	} else {
		return float64(tlb.NumHits) / float64(tlb.NumAccesses())
	}
}

func (tlb *TranslationLookasideBuffer) OccupancyRatio() float64 {
	return tlb.Cache.OccupancyRatio()
}

func (tlb *TranslationLookasideBuffer) HitLatency() uint32 {
	return tlb.MemoryHierarchy.Config().TlbHitLatency
}

func (tlb *TranslationLookasideBuffer) MissLatency() uint32 {
	return tlb.MemoryHierarchy.Config().TlbMissLatency
}

func (tlb *TranslationLookasideBuffer) Access(access *MemoryHierarchyAccess, onCompletedCallback func()) {
	var set = tlb.Cache.GetSet(access.PhysicalAddress)
	var cacheAccess = tlb.Cache.NewAccess(access, access.PhysicalAddress)

	if cacheAccess.HitInCache {
		tlb.Cache.ReplacementPolicy.HandlePromotionOnHit(access, set, cacheAccess.Way)

		tlb.NumHits++
	} else {
		if cacheAccess.Replacement {
			tlb.NumEvictions++
		}

		var line = tlb.Cache.Sets[set].Lines[cacheAccess.Way]
		line.StateProvider.(*BaseCacheLineStateProvider).SetState(true)
		line.Access = access
		line.Tag = int32(access.PhysicalTag)
		tlb.Cache.ReplacementPolicy.HandleInsertionOnMiss(access, set, cacheAccess.Way)

		tlb.NumMisses++
	}

	var delay uint32

	if cacheAccess.HitInCache {
		delay = tlb.HitLatency()
	} else {
		delay = tlb.MissLatency()
	}

	tlb.MemoryHierarchy.Driver().CycleAccurateEventQueue().Schedule(onCompletedCallback, int(delay))
}
