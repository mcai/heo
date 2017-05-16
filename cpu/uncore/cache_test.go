package uncore

import (
	"testing"
	"github.com/mcai/heo/cpu/mem"
	"github.com/mcai/heo/cpu/uncore/uncoreutil"
	"fmt"
)

func TestCache(t *testing.T) {
	var geometry = mem.NewGeometry(32 * uncoreutil.KB, 16, 64)

	var cache = NewCache(
		geometry,
		func(set uint32, way uint32) CacheLineStateProvider {
			return NewBaseCacheLineStateProvider("init_state")
		},
	)

	cache.Sets[0].Lines[0].StateProvider.(*BaseCacheLineStateProvider).SetState("test_state")

	fmt.Printf("len(cache.Sets): %d\n", len(cache.Sets))
	fmt.Printf("cache.Sets[0].Lines[0].State: %s\n", cache.Sets[0].Lines[0].State())
}
