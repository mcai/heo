package uncore

type CacheControllerState string

const (
	CacheControllerState_I     = CacheControllerState("I")
	CacheControllerState_IS_D  = CacheControllerState("IS_D")
	CacheControllerState_IM_AD = CacheControllerState("IM_AD")
	CacheControllerState_IM_A  = CacheControllerState("IM_A")
	CacheControllerState_S     = CacheControllerState("S")
	CacheControllerState_SM_AD = CacheControllerState("SM_AD")
	CacheControllerState_SM_A  = CacheControllerState("SM_A")
	CacheControllerState_M     = CacheControllerState("M")
	CacheControllerState_MI_A  = CacheControllerState("MI_A")
	CacheControllerState_SI_A  = CacheControllerState("SI_A")
	CacheControllerState_II_A  = CacheControllerState("II_A")
)

func (state CacheControllerState) Stable() bool {
	return state == CacheControllerState_I ||
		state == CacheControllerState_S ||
		state == CacheControllerState_M
}

func (state CacheControllerState) Transient() bool {
	return !state.Stable()
}

type DirectoryControllerState string

const (
	DirectoryControllerState_I    = DirectoryControllerState("I")
	DirectoryControllerState_IS_D = DirectoryControllerState("IS_D")
	DirectoryControllerState_IM_D = DirectoryControllerState("IM_D")
	DirectoryControllerState_S    = DirectoryControllerState("S")
	DirectoryControllerState_M    = DirectoryControllerState("M")
	DirectoryControllerState_S_D  = DirectoryControllerState("S_D")
	DirectoryControllerState_MI_A = DirectoryControllerState("MI_A")
	DirectoryControllerState_SI_A = DirectoryControllerState("SI_A")
)

func (state DirectoryControllerState) Stable() bool {
	return state == DirectoryControllerState_I ||
		state == DirectoryControllerState_S ||
		state == DirectoryControllerState_M
}

func (state DirectoryControllerState) Transient() bool {
	return !state.Stable()
}
