package uncore

import (
	"github.com/mcai/heo/cpu/mem"
)

const (
	INVALID_TAG = -1
)

const (
	INVALID_WAY = -1
)

type CacheLineStateProvider interface {
	State() interface{}
	Valid() bool
}

type BaseCacheLineStateProvider struct {
	state interface{}
	valid func(state interface{}) bool
}

func NewBaseCacheLineStateProvider(state interface{}, valid func(state interface{}) bool) *BaseCacheLineStateProvider {
	var stateProvider = &BaseCacheLineStateProvider{
		state:state,
		valid:valid,
	}

	return stateProvider
}

func (stateProvider *BaseCacheLineStateProvider) State() interface{} {
	return stateProvider.state
}

func (stateProvider *BaseCacheLineStateProvider) SetState(state interface{}) {
	stateProvider.state = state
}

func (stateProvider *BaseCacheLineStateProvider) Valid() bool {
	return stateProvider.valid(stateProvider.state)
}

type CacheLine struct {
	Cache         *Cache

	Set           uint32
	Way           uint32

	Tag           int32

	Access        *MemoryHierarchyAccess

	StateProvider CacheLineStateProvider
}

func newCacheLine(cache *Cache, set uint32, way uint32, stateProvider CacheLineStateProvider) *CacheLine {
	var cacheLine = &CacheLine{
		Cache:cache,
		Set:set,
		Way:way,
		Tag:INVALID_TAG,
		StateProvider:stateProvider,
	}

	return cacheLine
}

func (cacheLine *CacheLine) State() interface{} {
	return cacheLine.StateProvider.State()
}

func (cacheLine *CacheLine) Valid() bool {
	return cacheLine.StateProvider.Valid()
}

type CacheSet struct {
	Cache *Cache
	Lines []*CacheLine
	Num   uint32
}

func newCacheSet(cache *Cache, assoc uint32, num uint32) *CacheSet {
	var cacheSet = &CacheSet{
		Cache:cache,
		Num:num,
	}

	for i := uint32(0); i < assoc; i++ {
		cacheSet.Lines = append(cacheSet.Lines,
			newCacheLine(
				cache,
				num,
				i,
				cache.LineStateProviderFactory(num, i),
			),
		)
	}

	return cacheSet
}

type Cache struct {
	Geometry                 *mem.Geometry
	Sets                     []*CacheSet
	LineStateProviderFactory func(set uint32, way uint32) CacheLineStateProvider
}

func NewCache(geometry *mem.Geometry, lineStateProviderFactory func(set uint32, way uint32) CacheLineStateProvider) *Cache {
	var cache = &Cache{
		Geometry:geometry,
		LineStateProviderFactory:lineStateProviderFactory,
	}

	for i := uint32(0); i < geometry.NumSets; i++ {
		cache.Sets = append(cache.Sets, newCacheSet(cache, geometry.Assoc, i))
	}

	return cache
}

func (cache *Cache) GetTag(address uint32) uint32 {
	return cache.Geometry.GetTag(address)
}

func (cache *Cache) GetSet(address uint32) uint32 {
	return cache.Geometry.GetSet(address)
}

func (cache *Cache) NumSets() uint32 {
	return cache.Geometry.NumSets
}

func (cache *Cache) Assoc() uint32 {
	return cache.Geometry.Assoc
}

func (cache *Cache) LineSize() uint32 {
	return cache.Geometry.LineSize
}

func (cache *Cache) FindWay(address uint32) int32 {
	var tag = cache.GetTag(address)
	var set = cache.GetSet(address)

	for _, line := range cache.Sets[set].Lines {
		if line.Valid() && line.Tag == int32(tag) {
			return int32(line.Way)
		}
	}

	return INVALID_WAY
}

func (cache *Cache) FindLine(address uint32) *CacheLine {
	var set = cache.GetSet(address)
	var way = cache.FindWay(address)

	if way != INVALID_WAY {
		return cache.Sets[set].Lines[way]
	} else {
		return nil
	}
}

func (cache *Cache) OccupancyRatio() float64 {
	var numValidLines = 0

	for set := uint32(0); set < cache.Geometry.NumSets; set++ {
		for _, line := range cache.Sets[set].Lines {
			if line.Valid() {
				numValidLines++
			}
		}
	}

	return float64(numValidLines) / float64(cache.Geometry.NumLines)
}