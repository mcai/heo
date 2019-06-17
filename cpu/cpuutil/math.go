package cpuutil

func GetBits32(x uint32, high uint32, low uint32) uint32 {
	return (x >> low) & ((uint32(1) << (high - low + 1)) - 1)
}

func GetBit32(x uint32, bit uint32) uint32 {
	return x & (uint32(1) << bit)
}

func GetBits64(x uint64, high uint64, low uint64) uint64 {
	return (x >> low) & ((uint64(1) << (high - low + 1)) - 1)
}

func SetBit32(x uint32, b uint32) uint32 {
	return x | (uint32(1) << b)
}

func ClearBit32(x uint32, b uint32) uint32 {
	return x & ^(uint32(1) << b)
}

func SetBitValue32(x uint32, b uint32, value bool) uint32 {
	if value {
		return SetBit32(x, b)
	} else {
		return ClearBit32(x, b)
	}
}

func SignExtend32(x uint32, b uint32) uint32 {
	if (x&(uint32(1)<<(b-1))) != 0 {
		return x | ^((uint32(1) << b) - 1)
	} else {
		return x
	}
}

func RoundUp(n uint32, alignment uint32) uint32 {
	return (n + alignment - 1) & ^(alignment - 1)
}
