package uncore

type Controller interface {
	MemoryDevice
	Next() MemoryDevice
	SetNext(next MemoryDevice)
	HitLatency() uint32
	ReceiveMessage(message CoherenceMessage)
	TransferMessage(to Controller, size uint32, message CoherenceMessage)
}

type BaseController struct {
	*BaseMemoryDevice
	next MemoryDevice
}

func NewBaseController(memoryHierarchy MemoryHierarchy, name string, deviceType MemoryDeviceType) *BaseController {
	var controller = &BaseController{
		BaseMemoryDevice: NewBaseMemoryDevice(memoryHierarchy, name, deviceType),
	}

	return controller
}

func (controller *BaseController) HitLatency() uint32 {
	panic("Impossible")
}

func (controller *BaseController) ReceiveMessage(message CoherenceMessage) {
	panic("Impossible")
}

func (controller *BaseController) TransferMessage(to Controller, size uint32, message CoherenceMessage) {
	controller.MemoryHierarchy().TransferMessage(controller, to, size, message)
}

func (controller *BaseController) Next() MemoryDevice {
	return controller.next
}

func (controller *BaseController) SetNext(next MemoryDevice) {
	controller.next = next
}

type BaseCacheController struct {
	*BaseController

	NumDownwardReadHits    int64
	NumDownwardReadMisses  int64
	NumDownwardWriteHits   int64
	NumDownwardWriteMisses int64
	NumEvictions           int64
}

func NewBaseCacheController(memoryHierarchy MemoryHierarchy, name string, deviceType MemoryDeviceType) *BaseCacheController {
	var controller = &BaseCacheController{
		BaseController: NewBaseController(memoryHierarchy, name, deviceType),
	}

	return controller
}

func (controller *BaseCacheController) UpdateStats(write bool, hitInCache bool) {
	if write {
		if hitInCache {
			controller.NumDownwardWriteHits++
		} else {
			controller.NumDownwardWriteMisses++
		}
	} else {
		if hitInCache {
			controller.NumDownwardReadHits++
		} else {
			controller.NumDownwardReadMisses++
		}
	}
}

func (controller *BaseCacheController) NumDownwardHits() int64 {
	return controller.NumDownwardReadHits + controller.NumDownwardWriteHits
}

func (controller *BaseCacheController) NumDownwardMisses() int64 {
	return controller.NumDownwardReadMisses + controller.NumDownwardWriteMisses
}

func (controller *BaseCacheController) NumDownwardAccesses() int64 {
	return controller.NumDownwardHits() + controller.NumDownwardMisses()
}

func (controller *BaseCacheController) HitRatio() float64 {
	if controller.NumDownwardAccesses() == 0 {
		return 0
	} else {
		return float64(controller.NumDownwardHits()) / float64(controller.NumDownwardAccesses())
	}
}
