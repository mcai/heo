package cpu

import (
	"github.com/mcai/heo/cpu/mem"
	"github.com/mcai/heo/cpu/regs"
	"syscall"
)

func (syscallEmulation *SyscallEmulation) fstat64_impl(context *Context) {
	var fd = context.Process.TranslateFileDescriptor(int32(context.Regs().Gpr[regs.REGISTER_A0]))
	var bufAddr = context.Regs().Gpr[regs.REGISTER_A1]

	var fstat syscall.Stat_t

	if err := syscall.Fstat(int(fd), &fstat); err != nil {
		panic("Error")
	}

	context.Regs().Gpr[regs.REGISTER_V0] = 0

	syscallEmulation.Error = syscallEmulation.checkSyscallError(context)

	if !syscallEmulation.Error {
		var sizeOfDataToWrite = uint32(104)
		var dataToWrite = make([]byte, sizeOfDataToWrite)

		var memory = mem.NewSimpleMemory(context.Process.LittleEndian, dataToWrite)

		//TODO: correct?
		memory.WriteUInt32At(0, uint32(fstat.Dev))
		memory.WriteUInt32At(16, uint32(fstat.Ino))
		memory.WriteUInt32At(24, uint32(fstat.Mode))
		memory.WriteUInt32At(28, uint32(fstat.Nlink))
		memory.WriteUInt32At(32, uint32(fstat.Uid))
		memory.WriteUInt32At(36, uint32(fstat.Gid))
		memory.WriteUInt32At(40, uint32(fstat.Rdev))
		memory.WriteUInt32At(56, uint32(fstat.Size))

		memory.WriteUInt32At(64, uint32(fstat.Atimespec.Nano()))
		memory.WriteUInt32At(72, uint32(fstat.Mtimespec.Nano()))
		memory.WriteUInt32At(80, uint32(fstat.Ctimespec.Nano()))

		memory.WriteUInt32At(88, uint32(fstat.Blksize))
		memory.WriteUInt32At(96, uint32(fstat.Blocks))

		context.Process.Memory().WriteBlockAt(bufAddr, sizeOfDataToWrite, dataToWrite)
	}
}
