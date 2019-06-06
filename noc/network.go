package noc

import (
	"github.com/mcai/heo/simutil"
)

type NetworkType string

const (
	NetworkType_BASE = NetworkType("BASE")
	NetworkType_FIXED_LATENCY = NetworkType("FIXED_LATENCY")
)

type NetworkDriver interface {
	CycleAccurateEventQueue() *simutil.CycleAccurateEventQueue
}

type Network interface {
	CurrentPacketId() int64
	SetCurrentPacketId(currentPacketId int64)

	Driver() NetworkDriver
	Config() *NoCConfig
	Receive(packet Packet) bool
}
