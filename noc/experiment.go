package noc

import (
	"time"
	"github.com/mcai/heo/simutil"
	"fmt"
)

type NoCExperiment struct {
	cycleAccurateEventQueue *simutil.CycleAccurateEventQueue

	Network *Network

	BeginTime, EndTime time.Time

	Stats   simutil.Stats
	statMap map[string]interface{}
}

func NewNoCExperiment(config *NoCConfig) *NoCExperiment {
	var experiment = &NoCExperiment{
		cycleAccurateEventQueue: simutil.NewCycleAccurateEventQueue(),
	}

	experiment.Network = NewNetwork(experiment, config)

	switch dataPacketTraffic := config.DataPacketTraffic; dataPacketTraffic {
	case TrafficUniform:
		experiment.Network.AddTrafficGenerator(
			NewUniformTrafficGenerator(experiment.Network, config.DataPacketInjectionRate, config.MaxPackets, func(src int, dest int) Packet {
				return NewDataPacket(experiment.Network, src, dest, config.DataPacketSize, true, func() {})
			}),
		)
	case TrafficTranspose1:
		experiment.Network.AddTrafficGenerator(
			NewTranspose1TrafficGenerator(experiment.Network, config.DataPacketInjectionRate, config.MaxPackets, func(src int, dest int) Packet {
				return NewDataPacket(experiment.Network, src, dest, config.DataPacketSize, true, func() {})
			}),
		)
	case TrafficTranspose2:
		experiment.Network.AddTrafficGenerator(
			NewTranspose2TrafficGenerator(experiment.Network, config.DataPacketInjectionRate, config.MaxPackets, func(src int, dest int) Packet {
				return NewDataPacket(experiment.Network, src, dest, config.DataPacketSize, true, func() {})
			}),
		)
	case TrafficTrace:
		experiment.Network.AddTrafficGenerator(
			NewTraceTrafficGenerator(experiment.Network, config.DataPacketInjectionRate, config.MaxPackets, config.TraceFileName),
		)
	default:
		panic(fmt.Sprintf("data packet traffic %s is not supported", dataPacketTraffic))
	}

	return experiment
}

func (experiment *NoCExperiment) CycleAccurateEventQueue() *simutil.CycleAccurateEventQueue {
	return experiment.cycleAccurateEventQueue
}

func (experiment *NoCExperiment) Run() {
	experiment.BeginTime = time.Now()

	for (experiment.Network.Config.MaxCycles == -1 || experiment.CycleAccurateEventQueue().CurrentCycle < experiment.Network.Config.MaxCycles) && (experiment.Network.Config.MaxPackets == -1 || experiment.Network.NumPacketsReceived < experiment.Network.Config.MaxPackets) {
		experiment.CycleAccurateEventQueue().AdvanceOneCycle()
	}

	if experiment.Network.Config.DrainPackets {
		experiment.Network.AcceptPacket = false

		for experiment.Network.NumPacketsReceived != experiment.Network.NumPacketsTransmitted {
			experiment.CycleAccurateEventQueue().AdvanceOneCycle()
		}
	}

	experiment.EndTime = time.Now()

	experiment.Network.Config.Dump(experiment.Network.Config.OutputDirectory)

	experiment.DumpStats()
}

func (experiment *NoCExperiment) SimulationTime() time.Duration {
	return experiment.EndTime.Sub(experiment.BeginTime)
}

func (experiment *NoCExperiment) CyclesPerSecond() float64 {
	return float64(experiment.CycleAccurateEventQueue().CurrentCycle) / experiment.EndTime.Sub(experiment.BeginTime).Seconds()
}
