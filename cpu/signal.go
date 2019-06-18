package cpu

import (
	"github.com/mcai/heo/cpu/cpuutil"
	"github.com/mcai/heo/cpu/mem"
)

const MAX_SIGNAL = 64

type SignalMask struct {
	signals []uint32
}

func NewSignalMask() *SignalMask {
	var signalMask = &SignalMask{
		signals: make([]uint32, MAX_SIGNAL/2),
	}

	return signalMask
}

func (signalMask *SignalMask) Clone() *SignalMask {
	var newSignalMask = NewSignalMask()
	copy(newSignalMask.signals, signalMask.signals)
	return newSignalMask
}

func (signalMask *SignalMask) Set(signal uint32) {
	if signal < 1 || signal > MAX_SIGNAL {
		return
	}

	signal--

	signalMask.signals[signal/32] = cpuutil.SetBit32(signalMask.signals[signal/32], signal%32)
}

func (signalMask *SignalMask) Clear(signal uint32) {
	if signal < 1 || signal > MAX_SIGNAL {
		return
	}

	signal--

	signalMask.signals[signal/32] = cpuutil.ClearBit32(signalMask.signals[signal/32], signal%32)
}

func (signalMask *SignalMask) Contains(signal uint32) bool {
	if signal < 1 || signal > MAX_SIGNAL {
		return false
	}

	signal--

	return cpuutil.GetBit32(signalMask.signals[signal/32], signal%32) != 0
}

func (signalMask *SignalMask) LoadFrom(memory *mem.PagedMemory, virtualAddress uint32) {
	for i := uint32(0); i < MAX_SIGNAL/32; i++ {
		signalMask.signals[i] = memory.ReadUInt32At(virtualAddress + i*4)
	}
}

func (signalMask *SignalMask) SaveTo(memory *mem.PagedMemory, virtualAddress uint32) {
	for i := uint32(0); i < MAX_SIGNAL/32; i++ {
		memory.WriteUInt32At(virtualAddress+i*4, signalMask.signals[i])
	}
}

type SignalMasks struct {
	Pending *SignalMask
	Blocked *SignalMask
	Backup  *SignalMask
}

func NewSignalMasks() *SignalMasks {
	var signalMasks = &SignalMasks{
		Pending: NewSignalMask(),
		Blocked: NewSignalMask(),
		Backup:  NewSignalMask(),
	}

	return signalMasks
}

const (
	SignalAction_HANDLER_OFFSET  = 4
	SignalAction_RESTORER_OFFSET = 136
	SignalAction_MASK_OFFSET     = 8
)

type SignalAction struct {
	Flags    uint32
	Handler  uint32
	Restorer uint32
	Mask     *SignalMask
}

func NewSignalAction() *SignalAction {
	var signalAction = &SignalAction{
		Mask: NewSignalMask(),
	}

	return signalAction
}

func (signalAction *SignalAction) LoadFrom(memory *mem.PagedMemory, virtualAddress uint32) {
	signalAction.Flags = memory.ReadUInt32At(virtualAddress)
	signalAction.Handler = memory.ReadUInt32At(virtualAddress + SignalAction_HANDLER_OFFSET)
	signalAction.Restorer = memory.ReadUInt32At(virtualAddress + SignalAction_RESTORER_OFFSET)
	signalAction.Mask.LoadFrom(memory, virtualAddress+SignalAction_MASK_OFFSET)
}

func (signalAction *SignalAction) SaveTo(memory *mem.PagedMemory, virtualAddress uint32) {
	memory.WriteUInt32At(virtualAddress, signalAction.Flags)
	memory.WriteUInt32At(virtualAddress+SignalAction_HANDLER_OFFSET, signalAction.Handler)
	memory.WriteUInt32At(virtualAddress+SignalAction_RESTORER_OFFSET, signalAction.Restorer)
	signalAction.Mask.SaveTo(memory, virtualAddress+SignalAction_MASK_OFFSET)
}
