package cpu

import (
	"github.com/mcai/heo/cpu/mem"
	"github.com/mcai/heo/cpu/native"
	"github.com/mcai/heo/cpu/regs"
)

type SystemEventCriterion interface {
	NeedProcess(context *Context) bool
}

type TimeCriterion struct {
	When int64
}

func NewTimeCriterion() *TimeCriterion {
	var criterion = &TimeCriterion{
	}

	return criterion
}

func (criterion *TimeCriterion) NeedProcess(context *Context) bool {
	return criterion.When <= native.Clock(context.Kernel.Experiment.CycleAccurateEventQueue().CurrentCycle)
}

type SignalCriterion struct {
}

func NewSignalCriterion() *SignalCriterion {
	var criterion = &SignalCriterion{
	}

	return criterion
}

func (criterion *SignalCriterion) NeedProcess(context *Context) bool {
	for signal := uint32(1); signal <= MAX_SIGNAL; signal++ {
		if context.Kernel.MustProcessSignal(context, signal) {
			return true
		}
	}

	return false
}

type WaitForProcessIdCriterion struct {
	ProcessId int32
}

func NewWaitForProcessIdCriterion(processId int32) *WaitForProcessIdCriterion {
	var criterion = &WaitForProcessIdCriterion{
		ProcessId: processId,
	}

	return criterion
}

func (criterion *WaitForProcessIdCriterion) NeedProcess(context *Context) bool {
	if criterion.ProcessId == -1 {
		return context.Kernel.Experiment.Processor.NumZombies() > 0
	}

	var contextToWait = context.Kernel.GetContextFromProcessId(criterion.ProcessId)

	if contextToWait == nil || contextToWait.State == ContextState_FINISHED {
		return true
	}

	return false
}

type WaitForFileDescriptorCriterion struct {
	Buffer  *mem.CircularByteBuffer
	Address uint32
	Size    uint32
	Pufds   uint32
}

func NewWaitForFileDescriptorCriterion() *WaitForFileDescriptorCriterion {
	var criterion = &WaitForFileDescriptorCriterion{
	}

	return criterion
}

func (criterion *WaitForFileDescriptorCriterion) NeedProcess(context *Context) bool {
	return !criterion.Buffer.IsEmpty()
}

const (
	SystemEventType_READ           = 0
	SystemEventType_RESUME         = 1
	SystemEventType_WAIT           = 2
	SystemEventType_POLL           = 3
	SystemEventType_SIGNAL_SUSPEND = 4
)

type SystemEventType uint32

type SystemEvent interface {
	Context() *Context
	EventType() SystemEventType
	NeedProcess() bool
	Process()
}

type BaseSystemEvent struct {
	context   *Context
	eventType SystemEventType
}

func NewBaseSystemEvent(context *Context, eventType SystemEventType) *BaseSystemEvent {
	var event = &BaseSystemEvent{
		context:   context,
		eventType: eventType,
	}

	return event
}

func (event *BaseSystemEvent) Context() *Context {
	return event.context
}

func (event *BaseSystemEvent) EventType() SystemEventType {
	return event.eventType
}

type PollEvent struct {
	*BaseSystemEvent
	TimeCriterion                  *TimeCriterion
	WaitForFileDescriptorCriterion *WaitForFileDescriptorCriterion
}

func NewPollEvent(context *Context) *PollEvent {
	var event = &PollEvent{
		BaseSystemEvent:                NewBaseSystemEvent(context, SystemEventType_POLL),
		TimeCriterion:                  NewTimeCriterion(),
		WaitForFileDescriptorCriterion: NewWaitForFileDescriptorCriterion(),
	}

	return event
}

func (event *PollEvent) NeedProcess() bool {
	return event.TimeCriterion.NeedProcess(event.context) ||
		event.WaitForFileDescriptorCriterion.NeedProcess(event.context)
}

func (event *PollEvent) Process() {
	if !event.WaitForFileDescriptorCriterion.Buffer.IsEmpty() {
		event.Context().Process.Memory().WriteHalfWordAt(event.WaitForFileDescriptorCriterion.Pufds+6, 1)
		event.Context().Regs().Gpr[regs.REGISTER_V0] = 1
	} else {
		event.Context().Regs().Gpr[regs.REGISTER_V0] = 0
	}

	event.Context().Regs().Gpr[regs.REGISTER_A3] = 0
	event.context.Resume()
}

type ReadEvent struct {
	*BaseSystemEvent
	WaitForFileDescriptorCriterion *WaitForFileDescriptorCriterion
}

func NewReadEvent(context *Context) *ReadEvent {
	var event = &ReadEvent{
		BaseSystemEvent:                NewBaseSystemEvent(context, SystemEventType_READ),
		WaitForFileDescriptorCriterion: NewWaitForFileDescriptorCriterion(),
	}

	return event
}

func (event *ReadEvent) NeedProcess() bool {
	return event.WaitForFileDescriptorCriterion.NeedProcess(event.context)
}

func (event *ReadEvent) Process() {
	event.Context().Resume()

	var buf = event.WaitForFileDescriptorCriterion.Buffer.Read(event.WaitForFileDescriptorCriterion.Size)

	var numRead = uint32(len(buf))

	event.Context().Regs().Gpr[regs.REGISTER_V0] = numRead
	event.Context().Regs().Gpr[regs.REGISTER_A3] = 0

	event.Context().Process.Memory().WriteBlockAt(event.WaitForFileDescriptorCriterion.Address, numRead, buf)
}

type ResumeEvent struct {
	*BaseSystemEvent
	TimeCriterion *TimeCriterion
}

func NewResumeEvent(context *Context) *ResumeEvent {
	var event = &ResumeEvent{
		BaseSystemEvent: NewBaseSystemEvent(context, SystemEventType_RESUME),
		TimeCriterion:   NewTimeCriterion(),
	}

	return event
}

func (event *ResumeEvent) NeedProcess() bool {
	return event.TimeCriterion.NeedProcess(event.context)
}

func (event *ResumeEvent) Process() {
	event.Context().Resume()
}

type SignalSuspendEvent struct {
	*BaseSystemEvent
	SignalCriterion *SignalCriterion
}

func NewSignalSuspendEvent(context *Context) *SignalSuspendEvent {
	var event = &SignalSuspendEvent{
		BaseSystemEvent: NewBaseSystemEvent(context, SystemEventType_SIGNAL_SUSPEND),
		SignalCriterion: NewSignalCriterion(),
	}

	return event
}

func (event *SignalSuspendEvent) NeedProcess() bool {
	return event.SignalCriterion.NeedProcess(event.context)
}

func (event *SignalSuspendEvent) Process() {
	event.Context().Resume()

	event.Context().Kernel.ProcessSignals()

	event.Context().SignalMasks.Blocked = event.Context().SignalMasks.Backup.Clone()
}

type WaitEvent struct {
	*BaseSystemEvent
	WaitForProcessIdCriterion *WaitForProcessIdCriterion
	SignalCriterion           *SignalCriterion
}

func NewWaitEvent(context *Context, processId int32) *WaitEvent {
	var event = &WaitEvent{
		BaseSystemEvent:           NewBaseSystemEvent(context, SystemEventType_WAIT),
		WaitForProcessIdCriterion: NewWaitForProcessIdCriterion(processId),
		SignalCriterion:           NewSignalCriterion(),
	}

	return event
}

func (event *WaitEvent) NeedProcess() bool {
	return event.WaitForProcessIdCriterion.NeedProcess(event.context) ||
		event.SignalCriterion.NeedProcess(event.context)
}

func (event *WaitEvent) Process() {
	event.Context().Resume()

	event.Context().Regs().Gpr[regs.REGISTER_V0] = uint32(event.WaitForProcessIdCriterion.ProcessId)
	event.Context().Regs().Gpr[regs.REGISTER_A3] = 0
}
