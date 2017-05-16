package noc

import (
	"fmt"
	"github.com/mcai/heo/simutil"
)

func (experiment *NoCExperiment) DumpStats() {
	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "SimulationTime",
		Value: fmt.Sprintf("%v", experiment.SimulationTime()),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "SimulationTimeInSeconds",
		Value: experiment.SimulationTime().Seconds(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "TotalCycles",
		Value: experiment.CycleAccurateEventQueue().CurrentCycle,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "CyclesPerSecond",
		Value: experiment.CyclesPerSecond(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "PacketsPerSecond",
		Value: float64(experiment.Network.NumPacketsTransmitted) / experiment.EndTime.Sub(experiment.BeginTime).Seconds(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "NumPacketsReceived",
		Value: experiment.Network.NumPacketsReceived,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "NumPacketsTransmitted",
		Value: experiment.Network.NumPacketsTransmitted,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "Throughput",
		Value: experiment.Network.Throughput(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "AveragePacketDelay",
		Value: experiment.Network.AveragePacketDelay(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "AveragePacketHops",
		Value: experiment.Network.AveragePacketHops(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "MaxPacketDelay",
		Value: experiment.Network.MaxPacketDelay,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "MaxPacketHops",
		Value: experiment.Network.MaxPacketHops,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "NumPayloadPacketsReceived",
		Value: experiment.Network.NumPayloadPacketsReceived,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "NumPayloadPacketsTransmitted",
		Value: experiment.Network.NumPayloadPacketsTransmitted,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "PayloadThroughput",
		Value: experiment.Network.PayloadThroughput(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "AveragePayloadPacketDelay",
		Value: experiment.Network.AveragePayloadPacketDelay(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "AveragePayloadPacketHops",
		Value: experiment.Network.AveragePayloadPacketHops(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "MaxPayloadPacketDelay",
		Value: experiment.Network.MaxPayloadPacketDelay,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key: "MaxPayloadPacketHops",
		Value: experiment.Network.MaxPayloadPacketHops,
	})

	for _, state := range VALID_FLIT_STATES {
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key: fmt.Sprintf("AverageFlitPerStateDelay[%s]", state),
			Value: experiment.Network.AverageFlitPerStateDelay(state),
		})
	}

	for _, state := range VALID_FLIT_STATES {
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key: fmt.Sprintf("MaxFlitPerStateDelay[%s]", state),
			Value: experiment.Network.MaxFlitPerStateDelay[state],
		})
	}

	simutil.WriteJsonFile(experiment.Stats, experiment.Network.Config.OutputDirectory, simutil.STATS_JSON_FILE_NAME)
}

func (experiment *NoCExperiment) LoadStats() {
	simutil.LoadJsonFile(experiment.Network.Config.OutputDirectory, simutil.STATS_JSON_FILE_NAME, &experiment.statMap)
}

func (experiment *NoCExperiment) GetStatMap() map[string]interface{} {
	if experiment.statMap == nil {
		experiment.statMap = make(map[string]interface{})

		if experiment.Stats == nil {
			experiment.LoadStats()
		}

		for _, stat := range experiment.Stats {
			experiment.statMap[stat.Key] = stat.Value
		}
	}

	return experiment.statMap
}
