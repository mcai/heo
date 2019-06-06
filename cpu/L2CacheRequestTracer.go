package cpu

import (
	"encoding/csv"
	"fmt"
	"github.com/mcai/heo/cpu/uncore"
	"os"
	"reflect"
)

type L2CacheRequestTracer struct {
	TraceFileName string
	writer        *csv.Writer
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

	file, err := os.Create(traceFileName)

	if err != nil {
		panic("Cannot create file")
	}

	experiment.BlockingEventDispatcher().AddListener(reflect.TypeOf((*CPUExperimentStoppedEvent)(nil)), func(event interface{}) {
		defer file.Close()

		defer l2CacheRequestTracer.writer.Flush()
	})

	l2CacheRequestTracer.writer = csv.NewWriter(file)

	return l2CacheRequestTracer
}

func (l2CacheRequestTracer *L2CacheRequestTracer) handleL2Request(event *uncore.GeneralCacheControllerServiceNonblockingRequestEvent) {
	if !event.HitInCache {
		var _type = "W"

		if event.Access.AccessType.IsRead() {
			_type = "R"
		}

		var line = []string{fmt.Sprintf("%d", event.Access.ThreadId), fmt.Sprintf("%x", event.Access.VirtualPc), _type, fmt.Sprintf("%x", event.Access.PhysicalTag)}

		if err := l2CacheRequestTracer.writer.Write(line); err != nil {
			panic("Cannot write file")
		}
	}
}