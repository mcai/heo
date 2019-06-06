package cpu

import (
	"fmt"
	"github.com/mcai/heo/cpu/uncore"
	"reflect"
)

type L2CacheRequestTracer struct {
	TraceFileName string
}

func NewL2CacheRequestTracer(experiment *CPUExperiment, traceFileName string) *L2CacheRequestTracer  {
	var l2CacheRequestTracer = &L2CacheRequestTracer{
		TraceFileName:traceFileName,
	}

	experiment.BlockingEventDispatcher().AddListener(reflect.TypeOf((*uncore.GeneralCacheControllerServiceNonblockingRequestEvent)(nil)), func(event interface{}) {
		var e = event.(*uncore.GeneralCacheControllerServiceNonblockingRequestEvent)

		if e.CacheController == experiment.MemoryHierarchy.L2Controller() {
			l2CacheRequestTracer.handleL2Request(e)
		}
	})

	return l2CacheRequestTracer
}

func (l2CacheRequestTracer *L2CacheRequestTracer) handleL2Request(event *uncore.GeneralCacheControllerServiceNonblockingRequestEvent) {
	var _type = "W"

	if event.Access.AccessType.IsRead() {
		_type = "R"
	}

	fmt.Printf("%d,%x,%s,%x\n", event.Access.ThreadId, event.Access.VirtualPc, _type, event.Access.PhysicalTag)
}