package cpu

import (
	"fmt"
	"math"
	"github.com/mcai/heo/cpu/regs"
	"github.com/mcai/heo/cpu/native"
	"github.com/mcai/heo/cpu/mem"
)

type ErrNo uint32

const (
	EPERM ErrNo = 1

	ENOENT ErrNo = 2

	ESRCH ErrNo = 3

	EINTR ErrNo = 4

	EIO ErrNo = 5

	ENXIO ErrNo = 6

	E2BIG ErrNo = 7

	ENOEXEC ErrNo = 8

	EBADF ErrNo = 9

	ECHILD ErrNo = 10

	EAGAIN ErrNo = 11

	ENOMEM ErrNo = 12

	EACCES ErrNo = 13

	EFAULT ErrNo = 14

	ENOTBLK ErrNo = 15

	EBUSY ErrNo = 16

	EEXIST ErrNo = 17

	EXDEV ErrNo = 18

	ENODEV ErrNo = 19

	ENOTDIR ErrNo = 20

	EISDIR ErrNo = 21

	EINVAL ErrNo = 22

	ENFILE ErrNo = 23

	EMFILE ErrNo = 24

	ENOTTY ErrNo = 25

	ETXTBSY ErrNo = 26

	EFBIG ErrNo = 27

	ENOSPC ErrNo = 28

	ESPIPE ErrNo = 29

	EROFS ErrNo = 30

	EMLINK ErrNo = 31

	EPIPE ErrNo = 32

	EDOM ErrNo = 33

	ERANGE ErrNo = 34
)

const (
	MAX_BUFFER_SIZE = 1024
)

type TargetOpenFlag uint32

const (
	TargetOpenFlag_O_RDONLY TargetOpenFlag = 0

	TargetOpenFlag_O_WRONLY TargetOpenFlag = 1

	TargetOpenFlag_O_RDWR TargetOpenFlag = 2

	TargetOpenFlag_O_CREAT TargetOpenFlag = 0x100

	TargetOpenFlag_O_EXCL TargetOpenFlag = 0x400

	TargetOpenFlag_O_NOCTTY TargetOpenFlag = 0x800

	TargetOpenFlag_O_TRUNC TargetOpenFlag = 0x200

	TargetOpenFlag_O_APPEND TargetOpenFlag = 8

	TargetOpenFlag_O_NONBLOCK TargetOpenFlag = 0x80

	TargetOpenFlag_O_SYNC TargetOpenFlag = 0x10
)

type OpenFlag uint32

const (
	OpenFlag_O_RDONLY OpenFlag = 0x00000000

	OpenFlag_O_WRONLY OpenFlag = 0x00000001

	OpenFlag_O_RDWR OpenFlag = 0x00000002

	OpenFlag_O_CREAT OpenFlag = 0x00000040

	OpenFlag_O_EXCL OpenFlag = 0x00000080

	OpenFlag_O_NOCTTY OpenFlag = 0x00000100

	OpenFlag_O_TRUNC OpenFlag = 0x00000200

	OpenFlag_O_APPEND OpenFlag = 0x00000400

	OpenFlag_O_NONBLOCK OpenFlag = 0x00000800

	OpenFlag_O_SYNC OpenFlag = 0x00001000
)

type OpenFlagMapping struct {
	TargetFlag TargetOpenFlag
	HostFlag   OpenFlag
}

func NewOpenFlagMapping(targetFlag TargetOpenFlag, hostFlag OpenFlag) *OpenFlagMapping {
	var openFlagMapping = &OpenFlagMapping{
		TargetFlag:targetFlag,
		HostFlag:hostFlag,
	}

	return openFlagMapping
}

type SyscallHandler struct {
	Index uint32
	Name  string
	Run   func(context *Context)
}

func NewSyscallHandler(index uint32, name string, run func(context *Context)) *SyscallHandler {
	var syscallHandler = &SyscallHandler{
		Index:index,
		Name:name,
		Run:run,
	}

	return syscallHandler
}

type SyscallEmulation struct {
	Handlers         map[uint32]*SyscallHandler
	StackLimit       uint32
	Error            bool
	OpenFlagMappings []*OpenFlagMapping
}

func NewSyscallEmulation() *SyscallEmulation {
	var syscallEmulation = &SyscallEmulation{
		Handlers:make(map[uint32]*SyscallHandler),
		StackLimit: 0x800000,
	}

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_RDONLY, OpenFlag_O_RDONLY))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_WRONLY, OpenFlag_O_WRONLY))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_RDWR, OpenFlag_O_RDWR))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_APPEND, OpenFlag_O_APPEND))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_SYNC, OpenFlag_O_SYNC))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_CREAT, OpenFlag_O_CREAT))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_TRUNC, OpenFlag_O_TRUNC))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_EXCL, OpenFlag_O_EXCL))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag_O_NOCTTY, OpenFlag_O_NOCTTY))

	syscallEmulation.OpenFlagMappings = append(
		syscallEmulation.OpenFlagMappings,
		NewOpenFlagMapping(TargetOpenFlag(0x2000), OpenFlag(0)))

	syscallEmulation.registerHandler(NewSyscallHandler(1, "exit", syscallEmulation.exit_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(3, "read", syscallEmulation.read_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(4, "write", syscallEmulation.write_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(5, "open", syscallEmulation.open_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(6, "close", syscallEmulation.close_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(7, "waitpid", syscallEmulation.waitpid_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(20, "getpid", syscallEmulation.getpid_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(24, "getuid", syscallEmulation.getuid_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(37, "kill", syscallEmulation.kill_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(42, "pipe", syscallEmulation.pipe_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(45, "brk", syscallEmulation.brk_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(47, "getgid", syscallEmulation.getgid_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(49, "geteuid", syscallEmulation.geteuid_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(50, "getegid", syscallEmulation.getegid_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(54, "ioctl", syscallEmulation.ioctl_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(64, "getppid", syscallEmulation.getppid_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(75, "setrlimit", syscallEmulation.setrlimit_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(76, "getrlimit", syscallEmulation.getrlimit_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(90, "mmap", syscallEmulation.mmap_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(91, "munmap", syscallEmulation.munmap_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(120, "clone", syscallEmulation.clone_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(122, "uname", syscallEmulation.uname_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(125, "mprotect", syscallEmulation.mprotect_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(140, "_llseek", syscallEmulation._llseek_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(153, "_sysctl", syscallEmulation._sysctl_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(166, "nanosleep", syscallEmulation.nanosleep_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(167, "mremap", syscallEmulation.mremap_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(188, "poll", syscallEmulation.poll_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(194, "rt_sigaction", syscallEmulation.rt_sigaction_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(195, "rt_sigprocmask", syscallEmulation.rt_sigprocmask_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(199, "rt_sigsuspend", syscallEmulation.rt_sigsuspend_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(215, "fstat64", syscallEmulation.fstat64_impl))
	syscallEmulation.registerHandler(NewSyscallHandler(246, "exit_group", syscallEmulation.exit_group_impl))

	return syscallEmulation
}

func (syscallEmulation *SyscallEmulation) registerHandler(handler *SyscallHandler) {
	syscallEmulation.Handlers[handler.Index] = handler
}

func (syscallEmulation *SyscallEmulation) findAndRunSyscallHandler(syscallIndex uint32, context *Context) bool {
	if handler, ok := syscallEmulation.Handlers[syscallIndex]; ok {
		handler.Run(context)

		context.Kernel.ProcessSystemEvents()
		context.Kernel.ProcessSignals()

		return true
	}

	return false
}

func (syscallEmulation *SyscallEmulation) checkSyscallError(context *Context) bool {
	if context.Regs().Sgpr(regs.REGISTER_V0) != -1 {
		context.Regs().Gpr[regs.REGISTER_A3] = 0
		return false
	} else {
		context.Regs().Gpr[regs.REGISTER_V0] = 0
		context.Regs().Gpr[regs.REGISTER_A3] = 1
		return true
	}
}

func (syscallEmulation *SyscallEmulation) DoSyscall(callNum uint32, context *Context) {
	var syscallIndex = callNum - 4000

	if !syscallEmulation.findAndRunSyscallHandler(syscallIndex, context) {
		panic(fmt.Sprintf("ctx-%d: syscall %d (%d) not implemented", context.Id, callNum, syscallIndex))
	}
}

func (syscallEmulation *SyscallEmulation) exit_impl(context *Context) {
	context.Finish()
}

func (syscallEmulation *SyscallEmulation) read_impl(context *Context) {
	var readMaxSize = uint32(1 << 25)

	var fd = context.Process.TranslateFileDescriptor(int32(context.Regs().Gpr[regs.REGISTER_A0]))
	var bufAddr = context.Regs().Gpr[regs.REGISTER_A1]
	var size = uint32(math.Min(float64(readMaxSize), float64(context.Regs().Gpr[regs.REGISTER_A2])))

	var ret uint32
	var buf []byte

	var buffer = context.Kernel.GetReadBuffer(fd)
	if buffer != nil {
		if buffer.IsEmpty() {
			var e = NewReadEvent(context)
			e.WaitForFileDescriptorCriterion.Buffer = buffer
			e.WaitForFileDescriptorCriterion.Address = bufAddr
			e.WaitForFileDescriptorCriterion.Size = size
			context.Kernel.SystemEvents = append(context.Kernel.SystemEvents, e)
			context.Suspend()
			return
		} else {
			buf = buffer.Read(size)
			ret = uint32(len(buf))
		}
	} else {
		buf = make([]byte, size)
		ret = uint32(native.Read(fd, buf))
	}

	if ret >= readMaxSize {
		panic("Impossible")
	}

	context.Regs().Gpr[regs.REGISTER_V0] = uint32(ret)
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)

	context.Process.Memory().WriteBlockAt(bufAddr, ret, buf)
}

func (syscallEmulation *SyscallEmulation) write_impl(context *Context) {
	var fd = context.Process.TranslateFileDescriptor(int32(context.Regs().Gpr[regs.REGISTER_A0]))
	var bufAddr = context.Regs().Gpr[regs.REGISTER_A1]
	var size = context.Regs().Gpr[regs.REGISTER_A2]

	var buf = context.Process.Memory().ReadBlockAt(bufAddr, size)

	var ret uint32

	var buffer = context.Kernel.GetWriteBuffer(fd)
	if buffer != nil {
		buffer.Write(buf)
		ret = size
	} else {
		ret = uint32(native.Write(fd, buf))
	}

	context.Regs().Gpr[regs.REGISTER_V0] = ret
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) open_impl(context *Context) {
	var addr = context.Regs().Gpr[regs.REGISTER_A0]
	var targetFlags = context.Regs().Gpr[regs.REGISTER_A1]
	var mode = int32(context.Regs().Gpr[regs.REGISTER_A2])

	var hostFlags = uint32(0)
	for _, mapping := range syscallEmulation.OpenFlagMappings {
		if (targetFlags & uint32(mapping.TargetFlag)) != 0 {
			targetFlags &= ^uint32(mapping.TargetFlag)
			hostFlags |= uint32(mapping.HostFlag)
		}
	}

	if targetFlags != 0 {
		panic(fmt.Sprintf("syscall open: cannot decode flags 0x%08x", targetFlags))
	}

	var path = context.Process.Memory().ReadStringAt(addr, MAX_BUFFER_SIZE)

	var ret = native.Open(path, mode, hostFlags)

	context.Regs().Gpr[regs.REGISTER_V0] = uint32(ret)
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) close_impl(context *Context) {
	var fd = int32(context.Regs().Gpr[regs.REGISTER_A0])

	if fd == 0 || fd == 1 || fd == 2 {
		context.Regs().Gpr[regs.REGISTER_A3] = 0
		return
	}

	var ret = native.Close(fd)

	context.Regs().Gpr[regs.REGISTER_V0] = uint32(ret)
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) waitpid_impl(context *Context) {
	var pid = int32(context.Regs().Gpr[regs.REGISTER_A0])
	var pstatus = context.Regs().Gpr[regs.REGISTER_A1]

	if pid < 1 {
		panic("Impossible")
	}

	if context.Kernel.GetContextFromProcessId(pid) == nil {
		context.Regs().Gpr[regs.REGISTER_A3] = uint32(ECHILD)
		context.Regs().SetSgpr(regs.REGISTER_V0, -1)
		return
	}

	var e = NewWaitEvent(context, pid)
	context.Kernel.SystemEvents = append(context.Kernel.SystemEvents, e)
	context.Suspend()

	if pstatus != 0 {
		panic("Impossible")
	}
}

func (syscallEmulation *SyscallEmulation) getpid_impl(context *Context) {
	context.Regs().Gpr[regs.REGISTER_V0] = uint32(context.ProcessId)
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) getuid_impl(context *Context) {
	context.Regs().Gpr[regs.REGISTER_V0] = uint32(context.UserId)
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (SyscallEmulation *SyscallEmulation) kill_impl(context *Context) {
	var pid = int32(context.Regs().Gpr[regs.REGISTER_A0])
	var sig = context.Regs().Gpr[regs.REGISTER_A1]
	if pid < 0 {
		panic("Impossible")
	}

	var destContext = context.Kernel.GetContextFromProcessId(pid)
	if destContext == nil {
		context.Regs().Gpr[regs.REGISTER_A3] = uint32(ESRCH)
		context.Regs().SetSgpr(regs.REGISTER_V0, -1)
		return
	}

	destContext.SignalMasks.Pending.Set(sig)
	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) pipe_impl(context *Context) {
	var fileDescriptors = context.Kernel.CreatePipe()

	context.Regs().Gpr[regs.REGISTER_V0] = uint32(fileDescriptors[0])
	context.Regs().Gpr[regs.REGISTER_V1] = uint32(fileDescriptors[1])

	context.Regs().Gpr[regs.REGISTER_A3] = 0
}

func (syscallEmulation *SyscallEmulation) brk_impl(context *Context) {
	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (SyscallEmulation *SyscallEmulation) getgid_impl(context *Context) {
	context.Regs().Gpr[regs.REGISTER_V0] = uint32(context.GroupId)
	SyscallEmulation.Error = SyscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) geteuid_impl(context *Context) {
	context.Regs().Gpr[regs.REGISTER_V0] = uint32(context.EffectiveUserId)
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) getegid_impl(context *Context) {
	context.Regs().Gpr[regs.REGISTER_V0] = uint32(context.EffectiveGroupId)
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) ioctl_impl(context *Context) {
	var buf = make([]byte, 128)

	if context.Regs().Gpr[regs.REGISTER_A2] != 0 {
		buf = context.Process.Memory().ReadBlockAt(context.Regs().Gpr[regs.REGISTER_A2], 128)
	}

	var fd = context.Process.TranslateFileDescriptor(int32(context.Regs().Gpr[regs.REGISTER_A0]))
	if fd < 3 {
		context.Regs().Gpr[regs.REGISTER_V0] = uint32(native.Ioctl(fd, int32(context.Regs().Gpr[regs.REGISTER_A1]), buf))

		syscallEmulation.Error = syscallEmulation.checkSyscallError(context)

		if context.Regs().Gpr[regs.REGISTER_A2] != 0 {
			context.Process.Memory().WriteBlockAt(context.Regs().Gpr[regs.REGISTER_A2], 128, buf)
		}
	} else {
		context.Regs().Gpr[regs.REGISTER_A3] = 0
		context.Regs().Gpr[regs.REGISTER_V0] = 0
	}
}

func (syscallEmulation *SyscallEmulation) getppid_impl(context *Context) {
	context.Regs().Gpr[regs.REGISTER_V0] = uint32(context.GetParentProcessId())
	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) setrlimit_impl(context *Context) {
	if context.Regs().Gpr[regs.REGISTER_A0] != 3 {
		panic("Impossbile")
	}

	syscallEmulation.StackLimit = context.Regs().Gpr[regs.REGISTER_A1]

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) getrlimit_impl(context *Context) {
	var prlimit = context.Regs().Gpr[regs.REGISTER_A1]

	if context.Regs().Gpr[regs.REGISTER_A0] != 3 {
		panic("Impossible")
	}

	context.Process.Memory().WriteWordAt(prlimit, syscallEmulation.StackLimit)
	context.Process.Memory().WriteWordAt(prlimit + 4, 0xffffffff)

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) mmap_impl(context *Context) {
	var start = context.Regs().Gpr[regs.REGISTER_A0]
	var length = context.Regs().Gpr[regs.REGISTER_A1]

	var fd = int32(context.Process.Memory().ReadWordAt(context.Regs().Gpr[regs.REGISTER_SP] + 16))

	if fd != -1 {
		panic(fmt.Sprintf("syscall mmap: syscall is only supported with fd = -1 (%d)", fd))
	}

	if start == 0 {
		start = context.Process.HeapTop
	}

	start = context.Process.Memory().Map(start, length)

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = start
}

func (syscallEmulation *SyscallEmulation) munmap_impl(context *Context) {
	var start = context.Regs().Gpr[regs.REGISTER_A0]
	var length = context.Regs().Gpr[regs.REGISTER_A1]

	context.Process.Memory().Unmap(start, length)

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) clone_impl(context *Context) {
	var cloneFlags = context.Regs().Gpr[regs.REGISTER_A0]
	var newSp = context.Regs().Gpr[regs.REGISTER_A1]

	var newContext *Context = NewContextFromParent(context, context.Regs().Clone(), cloneFlags & 0xff)

	if !context.Kernel.Map(newContext, func(candidateThreadId int32) bool {
		return true
	}) {
		panic("Impossible")
	}

	context.Kernel.Contexts = append(context.Kernel.Contexts, newContext)

	newContext.Regs().Gpr[regs.REGISTER_SP] = newSp
	newContext.Regs().Gpr[regs.REGISTER_A3] = 0
	newContext.Regs().Gpr[regs.REGISTER_V0] = 0

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = uint32(newContext.ProcessId)
}

func (syscallEmulation *SyscallEmulation) uname_impl(context *Context) {
	var un = NewUtsname()
	un.Sysname = "Linux"
	un.Nodename = "sim"
	un.Release = "2.6"
	un.Version = "Tue Apr 5 12:21:57 UTC 2005"
	un.Machine = "mips"

	var un_buf = un.GetBytes(context.Process.LittleEndian)
	context.Process.Memory().WriteBlockAt(context.Regs().Gpr[regs.REGISTER_A0], uint32(len(un_buf)), un_buf)

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) mprotect_impl(context *Context) {
	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) _llseek_impl(context *Context) {
	//TODO: correct?
	var fd = context.Process.TranslateFileDescriptor(int32(context.Regs().Gpr[regs.REGISTER_A0]))
	var offset = int64(context.Regs().Gpr[regs.REGISTER_A1])
	var whence = int32(context.Regs().Gpr[regs.REGISTER_A2])

	var ret = native.Seek(fd, offset, whence)

	context.Regs().Gpr[regs.REGISTER_V0] = uint32(ret)

	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) _sysctl_impl(context *Context) {
	var buf = context.Process.Memory().ReadBlockAt(context.Regs().Gpr[regs.REGISTER_A0], 4 * 6)
	var memory = mem.NewSimpleMemory(context.Process.LittleEndian, buf)

	var args = NewSysctlArgs()
	args.Name = memory.ReadWord()
	args.Nlen = memory.ReadWord()
	args.Oldval = memory.ReadWord()
	args.Oldlenp = memory.ReadWord()
	args.Newval = memory.ReadWord()
	args.Newlen = memory.ReadWord()

	var buf2 = context.Process.Memory().ReadBlockAt(args.Name, 4 * 10)
	var memory2 = mem.NewSimpleMemory(context.Process.LittleEndian, buf2)

	var name = make([]uint32, 10)

	for i := 0; i < len(name); i++ {
		name[i] = memory2.ReadWord()
	}

	context.Regs().Gpr[regs.REGISTER_A3] = 0

	name[0] = 0 //TODO: hack for the moment

	if name[0] != 0 {
		panic("syscall sysctl is not supported with name[0] != 0")
	}
}

func (syscallEmulation *SyscallEmulation) nanosleep_impl(context *Context) {
	var preq = context.Regs().Gpr[regs.REGISTER_A0]
	var sec = context.Process.Memory().ReadWordAt(preq)
	var nsec = context.Process.Memory().ReadWordAt(preq + 4)

	var total = int64(sec * native.CLOCKS_PER_SEC + nsec / 1e9 * native.CLOCKS_PER_SEC)

	var e = NewResumeEvent(context)
	e.TimeCriterion.When = native.Clock(context.Kernel.Experiment.CycleAccurateEventQueue().CurrentCycle + total)
	context.Kernel.SystemEvents = append(context.Kernel.SystemEvents, e)
	context.Suspend()
}

func (sysallEmulation *SyscallEmulation) mremap_impl(context *Context) {
	var oldAddr = context.Regs().Gpr[regs.REGISTER_A0]
	var oldSize = context.Regs().Gpr[regs.REGISTER_A1]
	var newSize = context.Regs().Gpr[regs.REGISTER_A2]

	var start = context.Process.Memory().Remap(oldAddr, oldSize, newSize)

	context.Regs().Gpr[regs.REGISTER_V0] = start

	sysallEmulation.Error = sysallEmulation.checkSyscallError(context)
}

func (syscallEmulation *SyscallEmulation) poll_impl(context *Context) {
	var pufds = context.Regs().Gpr[regs.REGISTER_A0]
	var nfds = int32(context.Regs().Gpr[regs.REGISTER_A1])
	var timeout = int32(context.Regs().Gpr[regs.REGISTER_A2])

	if nfds < 1 {
		panic("syscall poll: nfds < 1")
	}

	for i := int32(0); i < nfds; i++ {
		var fd = int32(context.Process.Memory().ReadWordAt(pufds))
		var events = int16(context.Process.Memory().ReadHalfWordAt(pufds + 4))

		if events != 1 {
			panic(fmt.Sprintf("syscall poll: ufds.events (%d) != POLLIN", events))
		}

		var e = NewPollEvent(context)
		e.TimeCriterion.When = native.Clock(context.Kernel.Experiment.CycleAccurateEventQueue().CurrentCycle) + int64(timeout * native.CLOCKS_PER_SEC / 1000)
		e.WaitForFileDescriptorCriterion.Buffer = context.Kernel.GetReadBuffer(fd)

		if e.WaitForFileDescriptorCriterion.Buffer == nil {
			panic("syscall poll: fd does not belong to a pipe read buffer")
		}

		e.WaitForFileDescriptorCriterion.Pufds = pufds
		context.Kernel.SystemEvents = append(context.Kernel.SystemEvents, e)

		pufds += 8
	}

	context.Suspend()
}

func (syscallEmulation *SyscallEmulation) rt_sigaction_impl(context *Context) {
	var signum = context.Regs().Gpr[regs.REGISTER_A0]
	var pact = context.Regs().Gpr[regs.REGISTER_A1]
	var poact = context.Regs().Gpr[regs.REGISTER_A2]

	if poact != 0 {
		context.Kernel.SignalActions[signum - 1].SaveTo(context.Process.Memory(), poact)
	}

	if pact != 0 {
		context.Kernel.SignalActions[signum - 1].LoadFrom(context.Process.Memory(), pact)
	}

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) rt_sigprocmask_impl(context *Context) {
	var how = context.Regs().Gpr[regs.REGISTER_A0]
	var pset = context.Regs().Gpr[regs.REGISTER_A1]
	var poset = context.Regs().Gpr[regs.REGISTER_A2]

	if poset != 0 {
		context.SignalMasks.Blocked.SaveTo(context.Process.Memory(), poset)
	}

	if pset != 0 {
		var set = NewSignalMask()
		set.LoadFrom(context.Process.Memory(), pset)

		switch how {
		case 1:
			for i := uint32(1); i <= MAX_SIGNAL; i++ {
				if set.Contains(i) {
					context.SignalMasks.Blocked.Set(i)
				}
			}
		case 2:
			for i := uint32(1); i <= MAX_SIGNAL; i++ {
				if set.Contains(i) {
					context.SignalMasks.Blocked.Clear(i)
				}
			}
		case 3:
			context.SignalMasks.Blocked = set.Clone()
		}
	}

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) rt_sigsuspend_impl(context *Context) {
	var pmask = context.Regs().Gpr[regs.REGISTER_A0]

	if pmask == 0 {
		panic("syscall sigsuspend: mask is nil")
	}

	context.SignalMasks.Backup = context.SignalMasks.Blocked.Clone()

	context.SignalMasks.Blocked.LoadFrom(context.Process.Memory(), pmask)
	context.Suspend()

	var e = NewSignalSuspendEvent(context)
	context.Kernel.SystemEvents = append(context.Kernel.SystemEvents, e)

	context.Regs().Gpr[regs.REGISTER_A3] = 0
	context.Regs().Gpr[regs.REGISTER_V0] = 0
}

func (syscallEmulation *SyscallEmulation) exit_group_impl(context *Context) {
	context.Finish()
}