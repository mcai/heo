package noc

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type TraceFileLine struct {
	ThreadId int32
	Pc       int64
	Read     bool
	Ea       int64
}

type TraceFileBasedTrafficSource struct {
	Network              *BaseNetwork
	PacketInjectionRate  float64
	MaxPackets           int64
	TraceFileName        string
	TraceFileLines       []*TraceFileLine
	CurrentTraceFileLine int
}

func NewTraceFileBasedTrafficSource(network *BaseNetwork, packetInjectionRate float64, maxPackets int64, traceFileName string) *TraceFileBasedTrafficSource {
	var source = &TraceFileBasedTrafficSource{
		Network:             network,
		PacketInjectionRate: packetInjectionRate,
		MaxPackets:          maxPackets,
		TraceFileName:       traceFileName,
	}

	traceFile, err := os.Open(traceFileName)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(traceFile)
	for scanner.Scan() {
		var line = scanner.Text()
		var parts = strings.Split(line, ",")

		if parts[0] == "" {
			continue
		}

		threadId, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		if int(threadId) >= network.Config().NumNodes-2 {
			log.Printf("threadId is out of range, corresponding line: %s, threadId: %d, numNodes: %d\n", line, threadId, network.Config().NumNodes)
			continue
		}

		pc, err := strconv.ParseInt(parts[1], 16, 64)
		if err != nil {
			log.Fatal(err)
		}

		var read = parts[2] == "R"

		ea, err := strconv.ParseInt(parts[3], 16, 64)
		if err != nil {
			log.Fatal(err)
		}

		var traceFileLine = &TraceFileLine{
			ThreadId: int32(threadId),
			Pc:       pc,
			Read:     read,
			Ea:       ea,
		}

		source.TraceFileLines = append(source.TraceFileLines, traceFileLine)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := traceFile.Close(); err != nil {
		log.Fatal(err)
	}

	source.CurrentTraceFileLine = 0

	return source
}

func (source *TraceFileBasedTrafficSource) AdvanceOneCycle() {
	if !source.Network.AcceptPacket || source.MaxPackets != -1 && source.Network.NumPacketsReceived > source.MaxPackets {
		return
	}

	if rand.Float64() <= source.PacketInjectionRate {
		if source.CurrentTraceFileLine < len(source.TraceFileLines) {
			var traceFileLine = source.TraceFileLines[source.CurrentTraceFileLine]
			source.CurrentTraceFileLine += 1

			var src = int(traceFileLine.ThreadId)
			var dest = source.Network.Config().NumNodes - 1

			var packet = NewDataPacket(source.Network, src, dest, source.Network.Config().DataPacketSize, true, func() {})

			source.Network.driver.CycleAccurateEventQueue().Schedule(func() {
				source.Network.Receive(packet)
			}, 1)
		}
	}
}
