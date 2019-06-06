package uncore

import (
	"fmt"
	"github.com/mcai/heo/noc"
	"github.com/mcai/heo/simutil"
	"math"
	"reflect"
)

type UncoreDriver interface {
	CycleAccurateEventQueue() *simutil.CycleAccurateEventQueue
	BlockingEventDispatcher() *simutil.BlockingEventDispatcher
}

type MemoryHierarchy interface {
	Driver() UncoreDriver
	Config() *UncoreConfig

	CurrentMemoryHierarchyAccessId() int32
	SetCurrentMemoryHierarchyAccessId(currentMemoryHierarchyAccessId int32)

	CurrentCacheCoherenceFlowId() int32
	SetCurrentCacheCoherenceFlowId(currentCacheCoherenceFlowId int32)

	PendingFlows() []CacheCoherenceFlow
	SetPendingFlows(pendingFlows []CacheCoherenceFlow)

	MemoryController() *MemoryController
	L2Controller() *DirectoryController
	L1IControllers() []*L1IController
	L1DControllers() []*L1DController

	ITlbs() []*TranslationLookasideBuffer
	DTlbs() []*TranslationLookasideBuffer

	Network() *noc.Network

	Transfer(from MemoryDevice, to MemoryDevice, size uint32, onCompletedCallback func())
	TransferMessage(from Controller, to Controller, size uint32, message CoherenceMessage)

	DumpPendingFlowTree()

	ResetStats()
}

type BaseMemoryHierarchy struct {
	driver UncoreDriver
	config *UncoreConfig

	currentMemoryHierarchyAccessId int32
	currentCacheCoherenceFlowId    int32

	pendingFlows []CacheCoherenceFlow

	memoryController *MemoryController
	l2Controller     *DirectoryController
	l1IControllers   []*L1IController
	l1DControllers   []*L1DController

	iTlbs []*TranslationLookasideBuffer
	dTlbs []*TranslationLookasideBuffer

	p2pReorderBuffers map[Controller]map[Controller]*P2PReorderBuffer

	network          *noc.Network
	DevicesToNodeIds map[interface{}]uint32
}

func NewBaseMemoryHierarchy(driver UncoreDriver, config *UncoreConfig, nocConfig *noc.NoCConfig) *BaseMemoryHierarchy {
	var memoryHierarchy = &BaseMemoryHierarchy{
		driver:           driver,
		config:           config,
		DevicesToNodeIds: make(map[interface{}]uint32),
	}

	memoryHierarchy.memoryController = NewMemoryController(memoryHierarchy)

	memoryHierarchy.l2Controller = NewDirectoryController(memoryHierarchy, "l2")
	memoryHierarchy.l2Controller.SetNext(memoryHierarchy.memoryController)

	for i := int32(0); i < config.NumCores; i++ {
		var l1IController = NewL1IController(memoryHierarchy, fmt.Sprintf("c%d/icache", i))
		l1IController.SetNext(memoryHierarchy.l2Controller)
		memoryHierarchy.l1IControllers = append(memoryHierarchy.l1IControllers, l1IController)

		var l1DController = NewL1DController(memoryHierarchy, fmt.Sprintf("c%d/dcache", i))
		l1DController.SetNext(memoryHierarchy.l2Controller)
		memoryHierarchy.l1DControllers = append(memoryHierarchy.l1DControllers, l1DController)

		for j := int32(0); j < config.NumThreadsPerCore; j++ {
			memoryHierarchy.iTlbs = append(
				memoryHierarchy.iTlbs,
				NewTranslationLookasideBuffer(
					memoryHierarchy,
					fmt.Sprintf("c%dt%d/itlb", i, j),
				),
			)

			memoryHierarchy.dTlbs = append(
				memoryHierarchy.dTlbs,
				NewTranslationLookasideBuffer(
					memoryHierarchy,
					fmt.Sprintf("c%dt%d/dtlb", i, j),
				),
			)
		}
	}

	memoryHierarchy.p2pReorderBuffers = make(map[Controller]map[Controller]*P2PReorderBuffer)

	var numNodes = uint32(0)

	for i, l1IController := range memoryHierarchy.L1IControllers() {
		memoryHierarchy.DevicesToNodeIds[l1IController] = numNodes

		var l1DController = memoryHierarchy.L1DControllers()[i]

		memoryHierarchy.DevicesToNodeIds[l1DController] = numNodes

		numNodes++
	}

	memoryHierarchy.DevicesToNodeIds[memoryHierarchy.L2Controller()] = numNodes

	numNodes++

	memoryHierarchy.DevicesToNodeIds[memoryHierarchy.MemoryController()] = numNodes

	numNodes++

	var width = uint32(math.Sqrt(float64(numNodes)))

	if width*width != numNodes {
		numNodes = (width + 1) * (width + 1)
	}

	nocConfig.NumNodes = int(numNodes)
	nocConfig.MaxInputBufferSize = int(memoryHierarchy.l2Controller.Cache.LineSize() + 8)

	memoryHierarchy.network = noc.NewNetwork(driver.(noc.NetworkDriver), nocConfig)

	return memoryHierarchy
}

func (memoryHierarchy *BaseMemoryHierarchy) Driver() UncoreDriver {
	return memoryHierarchy.driver
}

func (memoryHierarchy *BaseMemoryHierarchy) Config() *UncoreConfig {
	return memoryHierarchy.config
}

func (memoryHierarchy *BaseMemoryHierarchy) CurrentMemoryHierarchyAccessId() int32 {
	return memoryHierarchy.currentMemoryHierarchyAccessId
}

func (memoryHierarchy *BaseMemoryHierarchy) SetCurrentMemoryHierarchyAccessId(currentMemoryHierarchyAccessId int32) {
	memoryHierarchy.currentMemoryHierarchyAccessId = currentMemoryHierarchyAccessId
}

func (memoryHierarchy *BaseMemoryHierarchy) CurrentCacheCoherenceFlowId() int32 {
	return memoryHierarchy.currentCacheCoherenceFlowId
}

func (memoryHierarchy *BaseMemoryHierarchy) SetCurrentCacheCoherenceFlowId(currentCacheCoherenceFlowId int32) {
	memoryHierarchy.currentCacheCoherenceFlowId = currentCacheCoherenceFlowId
}

func (memoryHierarchy *BaseMemoryHierarchy) PendingFlows() []CacheCoherenceFlow {
	return memoryHierarchy.pendingFlows
}

func (memoryHierarchy *BaseMemoryHierarchy) SetPendingFlows(pendingFlows []CacheCoherenceFlow) {
	memoryHierarchy.pendingFlows = pendingFlows
}

func (memoryHierarchy *BaseMemoryHierarchy) MemoryController() *MemoryController {
	return memoryHierarchy.memoryController
}

func (memoryHierarchy *BaseMemoryHierarchy) L2Controller() *DirectoryController {
	return memoryHierarchy.l2Controller
}

func (memoryHierarchy *BaseMemoryHierarchy) L1IControllers() []*L1IController {
	return memoryHierarchy.l1IControllers
}

func (memoryHierarchy *BaseMemoryHierarchy) L1DControllers() []*L1DController {
	return memoryHierarchy.l1DControllers
}

func (memoryHierarchy *BaseMemoryHierarchy) ITlbs() []*TranslationLookasideBuffer {
	return memoryHierarchy.iTlbs
}

func (memoryHierarchy *BaseMemoryHierarchy) DTlbs() []*TranslationLookasideBuffer {
	return memoryHierarchy.dTlbs
}

func (memoryHierarchy *BaseMemoryHierarchy) Network() *noc.Network {
	return memoryHierarchy.network
}

func (memoryHierarchy *BaseMemoryHierarchy) Transfer(from MemoryDevice, to MemoryDevice, size uint32, onCompletedCallback func()) {
	var src = memoryHierarchy.DevicesToNodeIds[from]
	var dest = memoryHierarchy.DevicesToNodeIds[to]

	if src != dest {
		var packet = noc.NewDataPacket(memoryHierarchy.network, int(src), int(dest), int(size), true, onCompletedCallback)

		memoryHierarchy.Driver().CycleAccurateEventQueue().Schedule(func() {
			memoryHierarchy.network.Receive(packet)
		}, 1)
	} else {
		onCompletedCallback()
	}
}

func (memoryHierarchy *BaseMemoryHierarchy) TransferMessage(from Controller, to Controller, size uint32, message CoherenceMessage) {
	if _, ok := memoryHierarchy.p2pReorderBuffers[from]; !ok {
		memoryHierarchy.p2pReorderBuffers[from] = make(map[Controller]*P2PReorderBuffer)
	}

	if _, ok := memoryHierarchy.p2pReorderBuffers[from][to]; !ok {
		memoryHierarchy.p2pReorderBuffers[from][to] = NewP2PReorderBuffer(from, to)
	}

	var p2pReorderBuffer = memoryHierarchy.p2pReorderBuffers[from][to]

	p2pReorderBuffer.Messages = append(p2pReorderBuffer.Messages, message)

	memoryHierarchy.Transfer(from, to, size, func() {
		p2pReorderBuffer.OnDestArrived(message)
	})
}

func (memoryHierarchy *BaseMemoryHierarchy) DumpPendingFlowTree() {
	for _, pendingFlow := range memoryHierarchy.pendingFlows {
		simutil.PrintNode(
			pendingFlow,
			func(node interface{}) interface{} {
				var cacheCoherenceFlow = node.(CacheCoherenceFlow)

				if cacheCoherenceFlow.Completed() {
					return fmt.Sprintf("%s -> created at %d, completed at %d", reflect.TypeOf(cacheCoherenceFlow), cacheCoherenceFlow.BeginCycle(), cacheCoherenceFlow.EndCycle())
				} else {
					return fmt.Sprintf("%s -> created at %d", reflect.TypeOf(cacheCoherenceFlow), cacheCoherenceFlow.BeginCycle())
				}
			},
			func(node interface{}) []interface{} {
				var cacheCoherenceFlow = node.(CacheCoherenceFlow)

				var children []interface{}

				for _, childFlow := range cacheCoherenceFlow.ChildFlows() {
					children = append(children, childFlow)
				}

				return children
			},
		)
		fmt.Println()
	}
}

func (memoryHierarchy *BaseMemoryHierarchy) ResetStats() {
}

type P2PReorderBuffer struct {
	Messages               []CoherenceMessage
	From                   Controller
	To                     Controller
	LastCompletedMessageId int32
}

func NewP2PReorderBuffer(from Controller, to Controller) *P2PReorderBuffer {
	var buffer = &P2PReorderBuffer{
		From:                   from,
		To:                     to,
		LastCompletedMessageId: -1,
	}

	return buffer
}

func (buffer *P2PReorderBuffer) OnDestArrived(message CoherenceMessage) {
	message.SetDestArrived(true)

	for len(buffer.Messages) > 0 {
		var message = buffer.Messages[0]

		if !message.DestArrived() {
			break
		}

		buffer.Messages = buffer.Messages[1:]

		buffer.To.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
			func() {
				message.Complete()

				buffer.LastCompletedMessageId = message.Id()

				buffer.To.ReceiveMessage(message)
			},
			int(buffer.To.HitLatency()),
		)
	}
}
