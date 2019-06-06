package cpu

import (
	"fmt"
	"github.com/mcai/heo/noc"
	"github.com/mcai/heo/simutil"
)

func (experiment *CPUExperiment) dumpStats(prefix string) {
	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "SimulationTime",
		Value: fmt.Sprintf("%v", experiment.SimulationTime()),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "SimulationTimeInSeconds",
		Value: experiment.SimulationTime().Seconds(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "TotalCycles",
		Value: experiment.CycleAccurateEventQueue().CurrentCycle,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "NumDynamicInsts",
		Value: experiment.Processor.NumDynamicInsts(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "CyclesPerSecond",
		Value: experiment.CyclesPerSecond(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "InstructionsPerSecond",
		Value: experiment.InstructionsPerSecond(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "InstructionsPerCycle",
		Value: experiment.Processor.InstructionsPerCycle(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "CyclesPerInstructions",
		Value: experiment.Processor.CyclesPerInstructions(),
	})

	for _, core := range experiment.Processor.Cores {
		for _, thread := range core.Threads() {
			experiment.Stats = append(experiment.Stats, simutil.Stat{
				Key:   fmt.Sprintf("thread_%d.NumDynamicInsts", thread.Id()),
				Value: thread.NumDynamicInsts(),
			})

			experiment.Stats = append(experiment.Stats, simutil.Stat{
				Key:   fmt.Sprintf("thread_%d.InstructionsPerCycle", thread.Id()),
				Value: thread.InstructionsPerCycle(),
			})

			experiment.Stats = append(experiment.Stats, simutil.Stat{
				Key:   fmt.Sprintf("thread_%d.CyclesPerInstructions", thread.Id()),
				Value: thread.CyclesPerInstructions(),
			})

			if oooThread := thread.(*OoOThread); oooThread != nil {
				experiment.Stats = append(experiment.Stats, simutil.Stat{
					Key:   fmt.Sprintf("thread_%d.BranchPredictor.HitRatio", thread.Id()),
					Value: oooThread.BranchPredictor.HitRatio(),
				})
				experiment.Stats = append(experiment.Stats, simutil.Stat{
					Key:   fmt.Sprintf("thread_%d.BranchPredictor.NumAccesses", thread.Id()),
					Value: oooThread.BranchPredictor.NumAccesses(),
				})
				experiment.Stats = append(experiment.Stats, simutil.Stat{
					Key:   fmt.Sprintf("thread_%d.BranchPredictor.NumHits", thread.Id()),
					Value: oooThread.BranchPredictor.NumHits(),
				})
				experiment.Stats = append(experiment.Stats, simutil.Stat{
					Key:   fmt.Sprintf("thread_%d.BranchPredictor.NumMisses", thread.Id()),
					Value: oooThread.BranchPredictor.NumMisses(),
				})
			}
		}
	}

	for i, itlb := range experiment.MemoryHierarchy.ITlbs() {
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("itlb_%d.HitRatio", i),
			Value: itlb.HitRatio(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("itlb_%d.NumAccesses", i),
			Value: itlb.NumAccesses(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("itlb_%d.NumHits", i),
			Value: itlb.NumHits,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("itlb_%d.NumMisses", i),
			Value: itlb.NumMisses,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("itlb_%d.NumEvictions", i),
			Value: itlb.NumEvictions,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("itlb_%d.OccupancyRatio", i),
			Value: itlb.OccupancyRatio(),
		})
	}

	for i, dtlb := range experiment.MemoryHierarchy.DTlbs() {
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dtlb_%d.HitRatio", i),
			Value: dtlb.HitRatio(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dtlb_%d.NumAccesses", i),
			Value: dtlb.NumAccesses(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dtlb_%d.NumHits", i),
			Value: dtlb.NumHits,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dtlb_%d.NumMisses", i),
			Value: dtlb.NumMisses,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dtlb_%d.NumEvictions", i),
			Value: dtlb.NumEvictions,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dtlb_%d.OccupancyRatio", i),
			Value: dtlb.OccupancyRatio(),
		})
	}

	for i, cacheController := range experiment.MemoryHierarchy.L1IControllers() {
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.HitRatio", i),
			Value: cacheController.HitRatio(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.NumDownwardAccesses", i),
			Value: cacheController.NumDownwardAccesses(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.NumDownwardHits", i),
			Value: cacheController.NumDownwardHits(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.NumDownwardMisses", i),
			Value: cacheController.NumDownwardMisses(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.NumDownwardReadHits", i),
			Value: cacheController.NumDownwardReadHits,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.NumDownwardReadMisses", i),
			Value: cacheController.NumDownwardReadMisses,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.NumDownwardWriteHits", i),
			Value: cacheController.NumDownwardWriteHits,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.NumDownwardWriteMisses", i),
			Value: cacheController.NumDownwardWriteMisses,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.NumEvictions", i),
			Value: cacheController.NumEvictions,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("icache_%d.OccupancyRatio", i),
			Value: cacheController.Cache.OccupancyRatio(),
		})
	}

	for i, cacheController := range experiment.MemoryHierarchy.L1DControllers() {
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.HitRatio", i),
			Value: cacheController.HitRatio(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.NumDownwardAccesses", i),
			Value: cacheController.NumDownwardAccesses(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.NumDownwardHits", i),
			Value: cacheController.NumDownwardHits(),
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.NumDownwardMisses", i),
			Value: cacheController.NumDownwardMisses(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.NumDownwardReadHits", i),
			Value: cacheController.NumDownwardReadHits,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.NumDownwardReadMisses", i),
			Value: cacheController.NumDownwardReadMisses,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.NumDownwardWriteHits", i),
			Value: cacheController.NumDownwardWriteHits,
		})
		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.NumDownwardWriteMisses", i),
			Value: cacheController.NumDownwardWriteMisses,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.NumEvictions", i),
			Value: cacheController.NumEvictions,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   fmt.Sprintf("dcache_%d.OccupancyRatio", i),
			Value: cacheController.Cache.OccupancyRatio(),
		})
	}

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.HitRatio",
		Value: experiment.MemoryHierarchy.L2Controller().HitRatio(),
	})
	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.NumDownwardAccesses",
		Value: experiment.MemoryHierarchy.L2Controller().NumDownwardAccesses(),
	})
	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.NumDownwardHits",
		Value: experiment.MemoryHierarchy.L2Controller().NumDownwardHits(),
	})
	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.NumDownwardMisses",
		Value: experiment.MemoryHierarchy.L2Controller().NumDownwardMisses(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.NumDownwardReadHits",
		Value: experiment.MemoryHierarchy.L2Controller().NumDownwardReadHits,
	})
	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.NumDownwardReadMisses",
		Value: experiment.MemoryHierarchy.L2Controller().NumDownwardReadMisses,
	})
	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.NumDownwardWriteHits",
		Value: experiment.MemoryHierarchy.L2Controller().NumDownwardWriteHits,
	})
	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.NumDownwardWriteMisses",
		Value: experiment.MemoryHierarchy.L2Controller().NumDownwardWriteMisses,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.NumEvictions",
		Value: experiment.MemoryHierarchy.L2Controller().NumEvictions,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2cache.OccupancyRatio",
		Value: experiment.MemoryHierarchy.L2Controller().Cache.OccupancyRatio(),
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "mem.NumReads",
		Value: experiment.MemoryHierarchy.MemoryController().NumReads,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "mem.NumWrites",
		Value: experiment.MemoryHierarchy.MemoryController().NumWrites,
	})

	if experiment.NocConfig.NetworkType == noc.NetworkType_BASE {
		var baseNetwork = experiment.MemoryHierarchy.Network().(*noc.BaseNetwork)

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.PacketsPerSecond",
			Value: float64(baseNetwork.NumPacketsTransmitted) / experiment.EndTime.Sub(experiment.BeginTime).Seconds(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.NumPacketsReceived",
			Value: baseNetwork.NumPacketsReceived,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.NumPacketsTransmitted",
			Value: baseNetwork.NumPacketsTransmitted,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.Throughput",
			Value: baseNetwork.Throughput(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.AveragePacketDelay",
			Value: baseNetwork.AveragePacketDelay(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.AveragePacketHops",
			Value: baseNetwork.AveragePacketHops(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.MaxPacketDelay",
			Value: baseNetwork.MaxPacketDelay,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.MaxPacketHops",
			Value: baseNetwork.MaxPacketHops,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.NumPayloadPacketsReceived",
			Value: baseNetwork.NumPayloadPacketsReceived,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.NumPayloadPacketsTransmitted",
			Value: baseNetwork.NumPayloadPacketsTransmitted,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.PayloadThroughput",
			Value: baseNetwork.PayloadThroughput(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.AveragePayloadPacketDelay",
			Value: baseNetwork.AveragePayloadPacketDelay(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.AveragePayloadPacketHops",
			Value: baseNetwork.AveragePayloadPacketHops(),
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.MaxPayloadPacketDelay",
			Value: baseNetwork.MaxPayloadPacketDelay,
		})

		experiment.Stats = append(experiment.Stats, simutil.Stat{
			Key:   "noc.MaxPayloadPacketHops",
			Value: baseNetwork.MaxPayloadPacketHops,
		})

		for _, state := range noc.VALID_FLIT_STATES {
			experiment.Stats = append(experiment.Stats, simutil.Stat{
				Key:   fmt.Sprintf("noc.AverageFlitPerStateDelay[%s]", state),
				Value: baseNetwork.AverageFlitPerStateDelay(state),
			})
		}

		for _, state := range noc.VALID_FLIT_STATES {
			experiment.Stats = append(experiment.Stats, simutil.Stat{
				Key:   fmt.Sprintf("noc.MaxFlitPerStateDelay[%s]", state),
				Value: baseNetwork.MaxFlitPerStateDelay[state],
			})
		}
	}

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumL2DemandHits",
		Value: experiment.L2PrefetchRequestProfiler.NumL2DemandHits,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumL2DemandMisses",
		Value: experiment.L2PrefetchRequestProfiler.NumL2DemandMisses,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumL2PrefetchHits",
		Value: experiment.L2PrefetchRequestProfiler.NumL2PrefetchHits,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumL2PrefetchMisses",
		Value: experiment.L2PrefetchRequestProfiler.NumL2PrefetchMisses,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumRedundantHitToTransientTagL2PrefetchRequests",
		Value: experiment.L2PrefetchRequestProfiler.NumRedundantHitToTransientTagL2PrefetchRequests,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumRedundantHitToCacheL2PrefetchRequests",
		Value: experiment.L2PrefetchRequestProfiler.NumRedundantHitToCacheL2PrefetchRequests,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumGoodL2PrefetchRequests",
		Value: experiment.L2PrefetchRequestProfiler.NumGoodL2PrefetchRequests,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumTimelyL2PrefetchRequests",
		Value: experiment.L2PrefetchRequestProfiler.NumTimelyL2PrefetchRequests,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumLateL2PrefetchRequests",
		Value: experiment.L2PrefetchRequestProfiler.NumLateL2PrefetchRequests,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumBadL2PrefetchRequests",
		Value: experiment.L2PrefetchRequestProfiler.NumBadL2PrefetchRequests,
	})

	experiment.Stats = append(experiment.Stats, simutil.Stat{
		Key:   "l2PrefetchRequestProfiler.NumEarlyL2PrefetchRequests",
		Value: experiment.L2PrefetchRequestProfiler.NumEarlyL2PrefetchRequests,
	})

	simutil.WriteJsonFile(experiment.Stats, experiment.CPUConfig.OutputDirectory, prefix+"_"+simutil.STATS_JSON_FILE_NAME)
}

func (experiment *CPUExperiment) ResetStats() {
	experiment.Stats = []simutil.Stat{}
	experiment.statMap = nil

	experiment.ISA.ResetStats()
	experiment.Kernel.ResetStats()
	experiment.Processor.ResetStats()
	experiment.MemoryHierarchy.ResetStats()
	experiment.OoO.ResetStats()

	experiment.L2PrefetchRequestProfiler.ResetStats()
}

func (experiment *CPUExperiment) LoadStats() {
	simutil.LoadJsonFile(experiment.CPUConfig.OutputDirectory, simutil.STATS_JSON_FILE_NAME, &experiment.statMap)
}

func (experiment *CPUExperiment) GetStatMap() map[string]interface{} {
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
