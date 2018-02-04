package uncore

type MemoryController struct {
	*BaseController
	NumReads  int32
	NumWrites int32
}

func NewMemoryController(memoryHierarchy MemoryHierarchy) *MemoryController {
	var memoryController = &MemoryController{
	}

	memoryController.BaseController = NewBaseController(
		memoryHierarchy,
		"mem",
		MemoryDeviceType_MEMORY_CONTROLLER,
	)

	return memoryController
}

func (memoryController *MemoryController) LineSize() uint32 {
	return memoryController.MemoryHierarchy().Config().MemoryControllerLineSize
}

func (memoryController *MemoryController) Latency() uint32 {
	return memoryController.MemoryHierarchy().Config().MemoryControllerLatency
}

func (memoryController *MemoryController) access(address uint32, onCompletedCallback func()) {
	memoryController.MemoryHierarchy().Driver().CycleAccurateEventQueue().Schedule(
		onCompletedCallback,
		int(memoryController.Latency()),
	)
}

func (memoryController *MemoryController) ReceiveMemReadRequest(source MemoryDevice, tag uint32, onCompletedCallback func()) {
	memoryController.NumReads++

	memoryController.access(
		tag,
		func() {
			memoryController.Transfer(
				source,
				source.(*DirectoryController).Cache.LineSize()+8,
				onCompletedCallback,
			)
		},
	)
}

func (memoryController *MemoryController) ReceiveMemWriteRequest(source MemoryDevice, tag uint32, onCompletedCallback func()) {
	memoryController.NumWrites++

	memoryController.access(
		tag,
		func() {
			memoryController.Transfer(
				source,
				8,
				onCompletedCallback,
			)
		},
	)
}
