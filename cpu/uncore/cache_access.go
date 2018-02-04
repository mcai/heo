package uncore

type CacheAccess struct {
	Cache       *EvictableCache
	Access      *MemoryHierarchyAccess
	Set         uint32
	Way         uint32
	Line        *CacheLine
	HitInCache  bool
	Replacement bool
}

func NewCacheAccess(cache *EvictableCache, access *MemoryHierarchyAccess, set uint32, way uint32, tag uint32) *CacheAccess {
	var cacheAccess = &CacheAccess{
		Cache:  cache,
		Access: access,
		Set:    set,
		Way:    way,
		Line:   cache.Sets[set].Lines[way],
	}

	cacheAccess.HitInCache = cacheAccess.Line.Tag == int32(tag)
	cacheAccess.Replacement = cacheAccess.Line.Valid()

	return cacheAccess
}
