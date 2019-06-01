package uncore

type LRUPolicy struct {
	*BaseCacheReplacementPolicy
	lruStack *Cache
}

func NewLRUPolicy(cache *EvictableCache) *LRUPolicy {
	var lruPolicy = &LRUPolicy{
		BaseCacheReplacementPolicy: NewBaseCacheReplacementPolicy(cache),
		lruStack: NewCache(
			cache.Geometry,
			func(set uint32, way uint32) CacheLineStateProvider {
				return NewBaseCacheLineStateProvider(
					way,
					func(state interface{}) bool {
						return true
					},
				)
			},
		),
	}

	return lruPolicy
}

func (lruPolicy *LRUPolicy) GetStackPositionOfWay(set uint32, way uint32) uint32 {
	for i, lruStackEntry := range lruPolicy.lruStack.Sets[set].Lines {
		if lruStackEntry.State().(uint32) == way {
			return uint32(i)
		}
	}

	panic("Impossible")
}

func (lruPolicy *LRUPolicy) GetWayInStackPosition(set uint32, stackPosition uint32) uint32 {
	return lruPolicy.lruStack.Sets[set].Lines[stackPosition].State().(uint32)
}

func (lruPolicy *LRUPolicy) SetMRU(set uint32, way uint32) {
	var oldStackPosition = lruPolicy.GetStackPositionOfWay(set, way)

	var stackEntry = lruPolicy.lruStack.Sets[set].Lines[oldStackPosition]

	var lines = lruPolicy.lruStack.Sets[set].Lines

	lines = append(lines[:oldStackPosition], lines[oldStackPosition+1:]...)
	lines = append([]*CacheLine{stackEntry}, lines...)

	lruPolicy.lruStack.Sets[set].Lines = lines
}

func (lruPolicy *LRUPolicy) GetLRU(set uint32) uint32 {
	return lruPolicy.GetWayInStackPosition(set, lruPolicy.cache.Assoc()-1)
}

func (lruPolicy *LRUPolicy) HandleReplacement(access *MemoryHierarchyAccess, set uint32, tag uint32) *CacheAccess {
	return NewCacheAccess(lruPolicy.cache, access, set, lruPolicy.GetLRU(set), tag)
}

func (lruPolicy *LRUPolicy) HandlePromotionOnHit(access *MemoryHierarchyAccess, set uint32, way uint32) {
	lruPolicy.SetMRU(set, way)
}

func (lruPolicy *LRUPolicy) HandleInsertionOnMiss(access *MemoryHierarchyAccess, set uint32, way uint32) {
	lruPolicy.SetMRU(set, way)
}
