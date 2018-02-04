package uncore

type CacheReplacementPolicy interface {
	Cache() *EvictableCache
	HandleReplacement(access *MemoryHierarchyAccess, set uint32, tag uint32) *CacheAccess
	HandlePromotionOnHit(access *MemoryHierarchyAccess, set uint32, way uint32)
	HandleInsertionOnMiss(access *MemoryHierarchyAccess, set uint32, way uint32)
}

type BaseCacheReplacementPolicy struct {
	cache *EvictableCache
}

func NewBaseCacheReplacementPolicy(cache *EvictableCache) *BaseCacheReplacementPolicy {
	var replacementPolicy = &BaseCacheReplacementPolicy{
		cache: cache,
	}

	return replacementPolicy
}

func (replacementPolicy *BaseCacheReplacementPolicy) Cache() *EvictableCache {
	return replacementPolicy.cache
}

func NewMiss(replacementPolicy CacheReplacementPolicy, access *MemoryHierarchyAccess, set uint32, address uint32) *CacheAccess {
	var tag = replacementPolicy.Cache().GetTag(address)

	for way := uint32(0); way < replacementPolicy.Cache().Assoc(); way++ {
		var line = replacementPolicy.Cache().Sets[set].Lines[way]
		if !line.Valid() {
			return NewCacheAccess(replacementPolicy.Cache(), access, set, way, tag)
		}
	}

	return replacementPolicy.HandleReplacement(access, set, tag)
}
