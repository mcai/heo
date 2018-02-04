package uncore

type MemoryDeviceType string

const (
	MemoryDeviceType_L1I_CONTROLLER    = MemoryDeviceType("l1i")
	MemoryDeviceType_L1D_CONTROLLER    = MemoryDeviceType("l1d")
	MemoryDeviceType_L2_CONTROLLER     = MemoryDeviceType("l2")
	MemoryDeviceType_MEMORY_CONTROLLER = MemoryDeviceType("mem")
)

type MemoryDevice interface {
	MemoryHierarchy() MemoryHierarchy
	Name() string
	DeviceType() MemoryDeviceType
	Transfer(to MemoryDevice, size uint32, onCompletedCallback func())
}

type BaseMemoryDevice struct {
	memoryHierarchy MemoryHierarchy
	name            string
	deviceType      MemoryDeviceType
}

func NewBaseMemoryDevice(memoryHierarchy MemoryHierarchy, name string, deviceType MemoryDeviceType) *BaseMemoryDevice {
	var memoryDevice = &BaseMemoryDevice{
		memoryHierarchy: memoryHierarchy,
		name:            name,
		deviceType:      deviceType,
	}

	return memoryDevice
}

func (memoryDevice *BaseMemoryDevice) Transfer(to MemoryDevice, size uint32, onCompletedCallback func()) {
	memoryDevice.memoryHierarchy.Transfer(memoryDevice, to, size, onCompletedCallback)
}

func (memoryDevice *BaseMemoryDevice) MemoryHierarchy() MemoryHierarchy {
	return memoryDevice.memoryHierarchy
}

func (memoryDevice *BaseMemoryDevice) Name() string {
	return memoryDevice.name
}

func (memoryDevice *BaseMemoryDevice) DeviceType() MemoryDeviceType {
	return memoryDevice.deviceType
}
