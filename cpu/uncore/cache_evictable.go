package uncore

import (
	"github.com/mcai/heo/cpu/mem"
	"fmt"
)

type EvictableCache struct {
	*Cache
	ReplacementPolicy CacheReplacementPolicy
}

func NewEvictableCache(geometry *mem.Geometry, lineStateProviderFactory func(set uint32, way uint32) CacheLineStateProvider, replacementPolicyType CacheReplacementPolicyType) *EvictableCache {
	var evictableCache = &EvictableCache{
		Cache: NewCache(geometry, lineStateProviderFactory),
	}

	switch replacementPolicyType {
	case CacheReplacementPolicyType_LRU:
		evictableCache.ReplacementPolicy = NewLRUPolicy(evictableCache)
	default:
		panic(fmt.Sprintf("cache replacement policy type %s is not supported", replacementPolicyType))
	}

	return evictableCache
}

func (evictableCache *EvictableCache) newHit(access *MemoryHierarchyAccess, set uint32, address uint32, way uint32) *CacheAccess {
	return NewCacheAccess(evictableCache, access, set, way, evictableCache.Cache.GetTag(address))
}

func (evictableCache *EvictableCache) newMiss(access *MemoryHierarchyAccess, set uint32, address uint32) *CacheAccess {
	return NewMiss(evictableCache.ReplacementPolicy, access, set, address)
}

func (evictableCache *EvictableCache) NewAccess(access *MemoryHierarchyAccess, address uint32) *CacheAccess {
	var line = evictableCache.Cache.FindLine(address)

	var set = evictableCache.Cache.GetSet(address)

	if line != nil {
		return evictableCache.newHit(access, set, address, line.Way)
	} else {
		return evictableCache.newMiss(access, set, address)
	}
}
