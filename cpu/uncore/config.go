package uncore

import (
	"github.com/mcai/heo/cpu/uncore/uncoreutil"
	"github.com/mcai/heo/simutil"
)

type CacheReplacementPolicyType string

const (
	CacheReplacementPolicyType_LRU = CacheReplacementPolicyType("LRU")
)

type UncoreConfig struct {
	NumCores          int32
	NumThreadsPerCore int32

	TlbSize        uint32
	TlbAssoc       uint32
	TlbLineSize    uint32
	TlbHitLatency  uint32
	TlbMissLatency uint32

	L1ISize              uint32
	L1IAssoc             uint32
	L1ILineSize          uint32
	L1IHitLatency        uint32
	L1INumReadPorts      uint32
	L1INumWritePorts     uint32
	L1IReplacementPolicy CacheReplacementPolicyType

	L1DSize              uint32
	L1DAssoc             uint32
	L1DLineSize          uint32
	L1DHitLatency        uint32
	L1DNumReadPorts      uint32
	L1DNumWritePorts     uint32
	L1DReplacementPolicy CacheReplacementPolicyType

	L2Size              uint32
	L2Assoc             uint32
	L2LineSize          uint32
	L2HitLatency        uint32
	L2ReplacementPolicy CacheReplacementPolicyType

	MemoryControllerLineSize uint32
	MemoryControllerLatency  uint32
}

func NewUncoreConfig(numCores int32, numThreadsPerCore int32) *UncoreConfig {
	var uncoreConfig = &UncoreConfig{
		NumCores:          numCores,
		NumThreadsPerCore: numThreadsPerCore,

		TlbSize:        32768,
		TlbAssoc:       4,
		TlbLineSize:    64,
		TlbHitLatency:  2,
		TlbMissLatency: 30,

		L1ISize:              64 * uncoreutil.KB,
		L1IAssoc:             4,
		L1ILineSize:          64,
		L1IHitLatency:        1,
		L1INumReadPorts:      128,
		L1INumWritePorts:     128,
		L1IReplacementPolicy: CacheReplacementPolicyType_LRU,

		L1DSize:              64 * uncoreutil.KB,
		L1DAssoc:             4,
		L1DLineSize:          64,
		L1DHitLatency:        1,
		L1DNumReadPorts:      128,
		L1DNumWritePorts:     128,
		L1DReplacementPolicy: CacheReplacementPolicyType_LRU,

		L2Size:              512 * uncoreutil.KB,
		L2Assoc:             16,
		L2LineSize:          64,
		L2HitLatency:        10,
		L2ReplacementPolicy: CacheReplacementPolicyType_LRU,

		MemoryControllerLineSize: 64,
		MemoryControllerLatency:  200,
	}

	return uncoreConfig
}

func (uncoreConfig *UncoreConfig) Dump(outputDirectory string) {
	simutil.WriteJsonFile(uncoreConfig, outputDirectory, simutil.UNCORE_CONFIG_JSON_FILE_NAME)
}
