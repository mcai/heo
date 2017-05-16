package cpu

type SysctlArgs struct {
	Name    uint32
	Nlen    uint32
	Oldval  uint32
	Oldlenp uint32
	Newval  uint32
	Newlen  uint32
}

func NewSysctlArgs() *SysctlArgs {
	var sysctlArgs = &SysctlArgs{
	}

	return sysctlArgs
}
