package cpu

import "github.com/mcai/heo/cpu/mem"

type Utsname struct {
	Sysname    string
	Nodename   string
	Release    string
	Version    string
	Machine    string
	Domainname string
}

func NewUtsname() *Utsname {
	var utsname = &Utsname{
	}

	return utsname
}

func (utsname *Utsname) GetBytes(littleEndian bool) []byte {
	var sysname_buf = []byte(utsname.Sysname + "\x00")
	var nodename_buf = []byte(utsname.Nodename + "\x00")
	var release_buf = []byte(utsname.Release + "\x00")
	var version_buf = []byte(utsname.Version + "\x00")
	var machine_buf = []byte(utsname.Machine + "\x00")
	var domainname_buf = []byte(utsname.Domainname + "\x00")

	var _sysname_size = uint32(64 + 1)
	var size_of = _sysname_size * 6

	var buf = make([]byte, size_of)

	var memory = mem.NewSimpleMemory(littleEndian, buf)

	memory.WriteBlockAt(0, uint32(len(sysname_buf)), sysname_buf)
	memory.WriteBlockAt(_sysname_size, uint32(len(nodename_buf)), nodename_buf)
	memory.WriteBlockAt(_sysname_size*2, uint32(len(release_buf)), release_buf)
	memory.WriteBlockAt(_sysname_size*3, uint32(len(version_buf)), version_buf)
	memory.WriteBlockAt(_sysname_size*4, uint32(len(machine_buf)), machine_buf)
	memory.WriteBlockAt(_sysname_size*5, uint32(len(domainname_buf)), domainname_buf)

	return buf
}
