package native

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	CLOCKS_PER_SEC = 1000000
	CPU_FREQUENCY  = 300000
)

func Getuid() int32 {
	return int32(syscall.Getuid())
}

func Geteuid() int32 {
	return int32(syscall.Geteuid())
}

func Getgid() int32 {
	return int32(syscall.Getgid())
}

func Getegid() int32 {
	return int32(syscall.Getegid())
}

func Read(fd int32, buf []byte) int32 {
	var count, _ = syscall.Read(int(fd), buf)
	return int32(count)
}

func Write(fd int32, buf []byte) int32 {
	var count, _ = syscall.Write(int(fd), buf)
	return int32(count)
}

func Open(path string, mode int32, perm uint32) int32 {
	var fd, _ = syscall.Open(path, int(mode), perm)
	return int32(fd)
}

func Close(fd int32) int32 {
	if err := syscall.Close(int(fd)); err != nil {
		panic(fmt.Sprintf("Cannot close fd (%s)", err))
	}

	panic("Unimplemented")

	return 0 //TODO
}

func Clock(numCycles int64) int64 {
	return CLOCKS_PER_SEC * numCycles / CPU_FREQUENCY
}

func Seek(fd int32, offset int64, whence int32) int64 {
	var off, _ = syscall.Seek(int(fd), offset, int(whence))
	return off
}

func Ioctl(fd int32, request int32, buf []byte) int64 {
	result, _, _ := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(request), uintptr(unsafe.Pointer(&buf[0])))
	return int64(result)
}
