package cpu

type L2PrefetchRequestQuality string

const (
	L2PrefetchRequestQuality_REDUNDANT_HIT_TO_TRANSIENT_TAG = L2PrefetchRequestQuality("REDUNDANT_HIT_TO_TRANSIENT_TAG")
	L2PrefetchRequestQuality_REDUNDANT_HIT_TO_CACHE = L2PrefetchRequestQuality("REDUNDANT_HIT_TO_CACHE")
	L2PrefetchRequestQuality_TIMELY = L2PrefetchRequestQuality("TIMELY")
	L2PrefetchRequestQuality_LATE = L2PrefetchRequestQuality("LATE")
	L2PrefetchRequestQuality_BAD = L2PrefetchRequestQuality("BAD")
	L2PrefetchRequestQuality_EARLY = L2PrefetchRequestQuality("EARLY")
	L2PrefetchRequestQuality_UGLY = L2PrefetchRequestQuality("UGLY")
	L2PrefetchRequestQuality_INVALID = L2PrefetchRequestQuality("INVALID")
)

func (quality L2PrefetchRequestQuality) Modifiable() bool {
	return quality == L2PrefetchRequestQuality_UGLY || quality == L2PrefetchRequestQuality_BAD
}

func (quality L2PrefetchRequestQuality) Useful() bool {
	return quality == L2PrefetchRequestQuality_TIMELY || quality == L2PrefetchRequestQuality_LATE
}

func (quality L2PrefetchRequestQuality) Polluting() bool {
	return quality == L2PrefetchRequestQuality_BAD
}