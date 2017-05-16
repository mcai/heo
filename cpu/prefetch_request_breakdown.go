package cpu

type L2PrefetchRequestBreakdown interface {
	NumL2DemandHits() int32
	NumL2DemandMisses() int32

	NumL2PrefetchHits() int32
	NumL2PrefetchMisses() int32

	NumRedundantHitToTransientTagL2PrefetchRequests() int32
	NumRedundantHitToCacheL2PrefetchRequests() int32

	NumTimelyL2PrefetchRequests() int32
	NumLateL2PrefetchRequests() int32

	NumBadL2PrefetchRequests() int32

	NumEarlyL2PrefetchRequests() int32

	NumUglyL2PrefetchRequests() int32
}

func NumTotalL2PrefetchRequests(breakdown L2PrefetchRequestBreakdown) int32 {
	return NumUsefulL2PrefetchRequests(breakdown) + NumUselessL2PrefetchRequests(breakdown)
}

func NumUsefulL2PrefetchRequests(breakdown L2PrefetchRequestBreakdown) int32 {
	return breakdown.NumLateL2PrefetchRequests() + breakdown.NumTimelyL2PrefetchRequests()
}

func NumUselessL2PrefetchRequests(breakdown L2PrefetchRequestBreakdown) int32 {
	return breakdown.NumBadL2PrefetchRequests() + breakdown.NumEarlyL2PrefetchRequests() + breakdown.NumUglyL2PrefetchRequests() +
		breakdown.NumRedundantHitToTransientTagL2PrefetchRequests() + breakdown.NumRedundantHitToCacheL2PrefetchRequests()
}
