package cpuutil

func Sext32(x uint32, b uint32) uint32 {
	if x & (uint32(1) << (b - 1)) != 0 {
		return x | ^((uint32(1) << b) - 1)
	} else {
		return x
	}
}

func Bits32(x uint32, hi uint32, lo uint32) uint32 {
	return (x >> lo) & ((uint32(1) << (hi - lo + 1)) - 1)
}

func Bits64(x uint64, hi uint64, lo uint64) uint64 {
	return (x >> lo) & ((uint64(1) << (hi - lo + 1)) - 1)
}

func GetBit32(x uint32, b uint32) uint32 {
	return x & (uint32(1) << b)
}

func GetBit64(x uint64, b uint64) uint64 {
	return x & (uint64(1) << b)
}

func SetBit32(x uint32, b uint32) uint32 {
	return x | (uint32(1) << b)
}

func SetBit64(x uint64, b uint64) uint64 {
	return x | (uint64(1) << b)
}

func ClearBit32(x uint32, b uint32) uint32 {
	return x & ^(uint32(1) << b)
}

func ClearBit64(x uint64, b uint64) uint64 {
	return x & ^(uint64(1) << b)
}

func SetBitValue32(x uint32, b uint32, v bool) uint32 {
	if v {
		return SetBit32(x, b)
	} else {
		return ClearBit32(x, b)
	}
}

func SetBitValue64(x uint64, b uint64, v bool) uint64 {
	if v {
		return SetBit64(x, b)
	} else {
		return ClearBit64(x, b)
	}
}

func RoundUp(n uint32, alignment uint32) uint32 {
	return (n + alignment - 1) & ^(alignment - 1)
}
