package cpu

import (
	"syscall"
	"github.com/mcai/heo/cpu/regs"
	"github.com/mcai/heo/cpu/mem"
)

func (syscallEmulation *SyscallEmulation) fstat64_impl(context *Context) {
	var fd = context.Process.TranslateFileDescriptor(int32(context.Regs().Gpr[regs.REGISTER_A0]))
	var bufAddr = context.Regs().Gpr[regs.REGISTER_A1]

	var fstat syscall.Stat_t

	syscall.Fstat(int(fd), &fstat)

	context.Regs().Gpr[regs.REGISTER_V0] = 0

	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)

	if !syscallEmulation.Error {
		var sizeOfDataToWrite = uint32(104)
		var dataToWrite = make([]byte, sizeOfDataToWrite)

		var memory = mem.NewSimpleMemory(context.Process.LittleEndian, dataToWrite)

		//TODO: correct?
		memory.WriteWordAt(0, uint32(fstat.Dev))
		memory.WriteWordAt(16, uint32(fstat.Ino))
		memory.WriteWordAt(24, uint32(fstat.Mode))
		memory.WriteWordAt(28, uint32(fstat.Nlink))
		memory.WriteWordAt(32, uint32(fstat.Uid))
		memory.WriteWordAt(36, uint32(fstat.Gid))
		memory.WriteWordAt(40, uint32(fstat.Rdev))
		memory.WriteWordAt(56, uint32(fstat.Size))

		memory.WriteWordAt(64, uint32(fstat.Atimespec.Nano()))
		memory.WriteWordAt(72, uint32(fstat.Mtimespec.Nano()))
		memory.WriteWordAt(80, uint32(fstat.Ctimespec.Nano()))

		memory.WriteWordAt(88, uint32(fstat.Blksize))
		memory.WriteWordAt(96, uint32(fstat.Blocks))

		context.Process.Memory().WriteBlockAt(bufAddr, sizeOfDataToWrite, dataToWrite)
	}
}
