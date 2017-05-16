package cpu

import (
	"github.com/mcai/heo/cpu/uncore"
	"reflect"
)

func DemandThreadId() int32 {
	return 0
}

func IsDemandThread(threadId int32) bool {
	return threadId == DemandThreadId()
}

func PrefetchThreadId() int32 {
	return 2
}

func IsPrefetchThread(threadId int32) bool {
	return threadId == PrefetchThreadId()
}

type L2PrefetchRequestProfiler struct {
	L2Controller                                    *uncore.DirectoryController

	L2PrefetchRequestStates                         map[int32](map[int32]*L2PrefetchRequestState)

	NumL2DemandHits                                 int32
	NumL2DemandMisses                               int32

	NumL2PrefetchHits                               int32
	NumL2PrefetchMisses                             int32

	NumRedundantHitToTransientTagL2PrefetchRequests int32
	NumRedundantHitToCacheL2PrefetchRequests        int32

	NumGoodL2PrefetchRequests                       int32

	NumTimelyL2PrefetchRequests                     int32
	NumLateL2PrefetchRequests                       int32

	NumBadL2PrefetchRequests                        int32

	NumEarlyL2PrefetchRequests                      int32
}

func NewL2PrefetchRequestProfiler(experiment *CPUExperiment) *L2PrefetchRequestProfiler {
	var l2PrefetchRequestProfiler = &L2PrefetchRequestProfiler{
		L2Controller:experiment.MemoryHierarchy.L2Controller(),
	}

	l2PrefetchRequestProfiler.L2PrefetchRequestStates = make(map[int32](map[int32]*L2PrefetchRequestState))

	for set := uint32(0); set < l2PrefetchRequestProfiler.L2Controller.Cache.NumSets(); set++ {
		l2PrefetchRequestProfiler.L2PrefetchRequestStates[int32(set)] = make(map[int32]*L2PrefetchRequestState)

		for way := uint32(0); way < l2PrefetchRequestProfiler.L2Controller.Cache.Assoc(); way++ {
			l2PrefetchRequestProfiler.L2PrefetchRequestStates[int32(set)][int32(way)] = NewL2PrefetchRequestState()
		}
	}

	experiment.BlockingEventDispatcher().AddListener(reflect.TypeOf((*uncore.GeneralCacheControllerServiceNonblockingRequestEvent)(nil)), func(event interface{}) {
		var e = event.(*uncore.GeneralCacheControllerServiceNonblockingRequestEvent)

		if e.CacheController == experiment.MemoryHierarchy.L2Controller() {
			var requesterIsPrefetch = IsPrefetchThread(e.Access.ThreadId)
			var lineFoundIsPrefetch = IsPrefetchThread(l2PrefetchRequestProfiler.L2PrefetchRequestStates[int32(e.Set)][int32(e.Way)].ThreadId)
			l2PrefetchRequestProfiler.handleL2Request(e, requesterIsPrefetch, lineFoundIsPrefetch)
		}
	})

	experiment.BlockingEventDispatcher().AddListener(reflect.TypeOf((*uncore.LastLevelCacheControllerLineInsertEvent)(nil)), func(event interface{}) {
		var e = event.(*uncore.LastLevelCacheControllerLineInsertEvent)

		if e.CacheController == experiment.MemoryHierarchy.L2Controller() {
			var lineFoundIsPrefetch = IsPrefetchThread(l2PrefetchRequestProfiler.L2PrefetchRequestStates[int32(e.Set)][int32(e.Way)].ThreadId)
			l2PrefetchRequestProfiler.handleL2LineInsert(e, lineFoundIsPrefetch)
		}
	})

	experiment.BlockingEventDispatcher().AddListener(reflect.TypeOf((*uncore.GeneralCacheControllerLastPutSOrPutMAndDataFromOwnerEvent)(nil)), func(event interface{}) {
		var e = event.(*uncore.GeneralCacheControllerLastPutSOrPutMAndDataFromOwnerEvent)

		if e.CacheController == experiment.MemoryHierarchy.L2Controller() {
			l2PrefetchRequestProfiler.markInvalid(e.Set, e.Way)
		}
	})

	experiment.BlockingEventDispatcher().AddListener(reflect.TypeOf((*uncore.GeneralCacheControllerNonblockingRequestHitToTransientTagEvent)(nil)), func(event interface{}) {
		var e = event.(*uncore.GeneralCacheControllerNonblockingRequestHitToTransientTagEvent)

		if e.CacheController == experiment.MemoryHierarchy.L2Controller() {
			var requesterIsPrefetch = IsPrefetchThread(e.Access.ThreadId)
			var lineFoundIsPrefetch = IsPrefetchThread(l2PrefetchRequestProfiler.L2PrefetchRequestStates[int32(e.Set)][int32(e.Way)].InFlightThreadId)

			if !requesterIsPrefetch && lineFoundIsPrefetch {
				l2PrefetchRequestProfiler.markLate(e.Set, e.Way, true)
			} else if requesterIsPrefetch && !lineFoundIsPrefetch {
				l2PrefetchRequestProfiler.markLate(e.Set, e.Way, true)
			}
		}
	})

	return l2PrefetchRequestProfiler
}

func (profiler *L2PrefetchRequestProfiler) handleL2Request(event *uncore.GeneralCacheControllerServiceNonblockingRequestEvent, requesterIsPrefetch bool, lineFoundIsPrefetch bool) {
	var victimWay = profiler.findWayOfL2LineByVictimTag(event.Set, event.Way)

	var victimLine *uncore.CacheLine = nil
	var victimLineState *L2PrefetchRequestState = nil

	if victimWay != -1 {
		victimLine = profiler.L2Controller.Cache.Sets[int32(event.Set)].Lines[victimWay]
		victimLineState = profiler.L2PrefetchRequestStates[int32(event.Set)][victimWay]
	}

	var victimHit = victimLine != nil
	var victimEvicterDemandHit = victimHit && IsDemandThread(victimLineState.ThreadId)
	var victimEvicterPrefetchHit = victimHit && !victimLineState.Used && IsPrefetchThread(victimLineState.ThreadId)
	var victimDemandHit = victimHit && IsDemandThread(victimLineState.VictimThreadId)
	var victimPrefetchHit = victimHit && IsPrefetchThread(victimLineState.VictimThreadId)

	var l2Line = profiler.L2Controller.Cache.Sets[int32(event.Set)].Lines[int32(event.Way)]
	var l2LineState = profiler.L2PrefetchRequestStates[int32(event.Set)][int32(event.Way)]

	var demandHit = event.HitInCache && !requesterIsPrefetch && !lineFoundIsPrefetch
	var prefetchHit = event.HitInCache && !l2LineState.Used && !requesterIsPrefetch && lineFoundIsPrefetch

	if !requesterIsPrefetch {
		if event.HitInCache {
			profiler.NumL2DemandHits++

			if lineFoundIsPrefetch && !l2LineState.Used {
				profiler.NumGoodL2PrefetchRequests++
			}
		} else {
			profiler.NumL2DemandMisses++
		}
	} else {
		if event.HitInCache {
			profiler.NumL2PrefetchHits++
		} else {
			profiler.NumL2PrefetchMisses++
		}

		if event.HitInCache && !lineFoundIsPrefetch {
			profiler.redundant(event, l2LineState)
		}
	}

	if !event.HitInCache {
		profiler.setL2LineBroughterThreadId(event.Set, event.Way, event.Access.ThreadId, event.Access.VirtualPc, true)
	}

	if !requesterIsPrefetch {
		if !demandHit && !prefetchHit {
			if !victimHit {
				//No action.
			} else if victimEvicterDemandHit && victimDemandHit {
				//No action.
			} else if victimEvicterPrefetchHit && victimDemandHit {
				profiler.bad(event, victimLine, victimLineState)
			} else if victimEvicterDemandHit && victimPrefetchHit {
				profiler.early(event, victimLineState)
			} else if victimEvicterPrefetchHit && victimPrefetchHit {
				//Ugly.
			}
		} else if prefetchHit {
			if !victimHit {
				profiler.good(event, l2Line, l2LineState)
			} else if victimEvicterDemandHit && victimDemandHit {
				profiler.good(event, l2Line, l2LineState)
			} else if victimEvicterPrefetchHit && victimDemandHit {
				profiler.good(event, l2Line, l2LineState)
				profiler.bad(event, victimLine, victimLineState)
			} else if victimEvicterDemandHit && victimPrefetchHit {
				profiler.good(event, l2Line, l2LineState)
				profiler.early(event, victimLineState)
			} else if victimEvicterPrefetchHit && victimPrefetchHit {
				profiler.good(event, l2Line, l2LineState)
			}
		} else {
			if !victimHit {
				//No action.
			} else if victimEvicterDemandHit && victimDemandHit {
				//No action.
			} else if victimEvicterPrefetchHit && victimDemandHit {
				//Bandwidth waste.
			} else if victimEvicterDemandHit && victimPrefetchHit {
				//Bandwidth waste.
			} else if victimEvicterPrefetchHit && victimPrefetchHit {
				//Bandwidth waste.
			}
		}
	}

	if event.HitInCache {
		l2LineState.VictimThreadId = l2LineState.ThreadId
		l2LineState.VictimPc = l2LineState.Pc
		l2LineState.VictimTag = l2Line.Tag
		profiler.setL2LineBroughterThreadId(event.Set, event.Way, event.Access.ThreadId, event.Access.VirtualPc, false)

		l2LineState.Used = requesterIsPrefetch && !lineFoundIsPrefetch
	}

	if victimHit {
		victimLineState.VictimThreadId = -1
		victimLineState.VictimPc = -1
		victimLineState.VictimTag = uncore.INVALID_TAG
	}
}

func (profiler *L2PrefetchRequestProfiler) redundant(event *uncore.GeneralCacheControllerServiceNonblockingRequestEvent, l2LineState *L2PrefetchRequestState) {
	//Redundant.
	if l2LineState.HitToTransientTag {
		profiler.NumRedundantHitToTransientTagL2PrefetchRequests++
	} else {
		profiler.NumRedundantHitToCacheL2PrefetchRequests++
	}
}

func (profiler *L2PrefetchRequestProfiler) good(event *uncore.GeneralCacheControllerServiceNonblockingRequestEvent, l2Line *uncore.CacheLine, l2LineState *L2PrefetchRequestState) {
	//Good.
	if l2LineState.HitToTransientTag {
		profiler.NumLateL2PrefetchRequests++
	} else {
		profiler.NumTimelyL2PrefetchRequests++
	}
}

func (profiler *L2PrefetchRequestProfiler) bad(event *uncore.GeneralCacheControllerServiceNonblockingRequestEvent, victimLine *uncore.CacheLine, victimLineState *L2PrefetchRequestState) {
	//Bad.
	profiler.NumBadL2PrefetchRequests++
	victimLineState.Used = true
}

func (profiler *L2PrefetchRequestProfiler) early(event *uncore.GeneralCacheControllerServiceNonblockingRequestEvent, victimLineState *L2PrefetchRequestState) {
	//Early.
	profiler.NumEarlyL2PrefetchRequests++
}

func (profiler *L2PrefetchRequestProfiler) handleL2LineInsert(event *uncore.LastLevelCacheControllerLineInsertEvent, lineFoundIsPrefetch bool) {
	var l2LineState = profiler.L2PrefetchRequestStates[int32(event.Set)][int32(event.Way)]

	if !lineFoundIsPrefetch && l2LineState.Used {
		panic("Impossible")
	}

	if !event.Eviction() {
		l2LineState.VictimThreadId = -1
		l2LineState.VictimPc = -1
		l2LineState.VictimTag = uncore.INVALID_TAG
	} else {
		l2LineState.VictimThreadId = l2LineState.ThreadId
		l2LineState.VictimPc = l2LineState.Pc
		l2LineState.VictimTag = event.VictimTag
	}

	profiler.setL2LineBroughterThreadId(event.Set, event.Way, event.Access.ThreadId, event.Access.VirtualPc, false)
	l2LineState.Used = false
}

func (profiler *L2PrefetchRequestProfiler) findWayOfL2LineByVictimTag(set uint32, victimTag uint32) int32 {
	for way := uint32(0); way < profiler.L2Controller.Cache.Assoc(); way++ {
		var state = profiler.L2PrefetchRequestStates[int32(set)][int32(way)]
		if state.VictimTag == int32(victimTag) {
			return int32(way)
		}
	}

	return -1
}

func (profiler *L2PrefetchRequestProfiler) markInvalid(set uint32, way uint32) {
	var l2LineState = profiler.L2PrefetchRequestStates[int32(set)][int32(way)]

	profiler.setL2LineBroughterThreadId(set, way, -1, -1, false)

	l2LineState.Pc = -1
	l2LineState.VictimThreadId = -1
	l2LineState.VictimPc = -1
	l2LineState.VictimTag = uncore.INVALID_TAG
	l2LineState.Used = false

	profiler.markLate(set, way, false)
}

func (profiler *L2PrefetchRequestProfiler) setL2LineBroughterThreadId(set uint32, way uint32, l2LineBroughterThreadId int32, pc int32, inflight bool) {
	var l2LineState = profiler.L2PrefetchRequestStates[int32(set)][int32(way)]

	if inflight {
		l2LineState.InFlightThreadId = l2LineBroughterThreadId
		l2LineState.Pc = pc
	} else {
		l2LineState.InFlightThreadId = -1
		l2LineState.ThreadId = l2LineBroughterThreadId
	}
}

func (profiler *L2PrefetchRequestProfiler) markLate(set uint32, way uint32, late bool) {
	var l2LineState = profiler.L2PrefetchRequestStates[int32(set)][int32(way)]
	l2LineState.HitToTransientTag = late
}

func (profiler *L2PrefetchRequestProfiler) NumUglyL2PrefetchRequests() int32 {
	return profiler.NumL2PrefetchHits + profiler.NumL2PrefetchMisses -
		profiler.NumRedundantHitToCacheL2PrefetchRequests -
		profiler.NumRedundantHitToTransientTagL2PrefetchRequests -
		profiler.NumTimelyL2PrefetchRequests - profiler.NumLateL2PrefetchRequests -
		profiler.NumBadL2PrefetchRequests - profiler.NumEarlyL2PrefetchRequests
}

func (profiler *L2PrefetchRequestProfiler) ResetStats() {
	profiler.NumL2DemandHits = 0
	profiler.NumL2DemandMisses = 0

	profiler.NumL2PrefetchHits = 0
	profiler.NumL2PrefetchMisses = 0

	profiler.NumRedundantHitToTransientTagL2PrefetchRequests = 0
	profiler.NumRedundantHitToCacheL2PrefetchRequests = 0

	profiler.NumGoodL2PrefetchRequests = 0

	profiler.NumTimelyL2PrefetchRequests = 0
	profiler.NumLateL2PrefetchRequests = 0

	profiler.NumBadL2PrefetchRequests = 0

	profiler.NumEarlyL2PrefetchRequests = 0
}