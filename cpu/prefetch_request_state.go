package cpu

type L2PrefetchRequestState struct {
	InFlightThreadId  int32
	ThreadId          int32
	Pc                int32
	Used              bool
	HitToTransientTag bool

	VictimThreadId    int32
	VictimPc          int32
	VictimTag         int32
}

func NewL2PrefetchRequestState() *L2PrefetchRequestState {
	var l2PrefetchRequestState = &L2PrefetchRequestState{
		InFlightThreadId:-1,
		ThreadId:-1,
		Pc:-1,
		Used:false,

		VictimThreadId:-1,
		VictimPc:-1,
		VictimTag:-1,
	}

	return l2PrefetchRequestState
}