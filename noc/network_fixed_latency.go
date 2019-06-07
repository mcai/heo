package noc

type FixedLatencyNetwork struct {
	driver NetworkDriver
	config *NoCConfig

	currentPacketId int64
}

func NewFixedLatencyNetwork(driver NetworkDriver, config *NoCConfig) *FixedLatencyNetwork {
	var fixedLatencyNetwork = &FixedLatencyNetwork{
		driver: driver,
		config: config,
	}

	return fixedLatencyNetwork
}

func (fixedLatencyNetwork *FixedLatencyNetwork) CurrentPacketId() int64  {
	return fixedLatencyNetwork.currentPacketId
}

func (fixedLatencyNetwork *FixedLatencyNetwork) SetCurrentPacketId(currentPacketId int64)  {
	fixedLatencyNetwork.currentPacketId = currentPacketId
}

func (fixedLatencyNetwork *FixedLatencyNetwork) Driver() NetworkDriver {
	return fixedLatencyNetwork.driver
}

func (fixedLatencyNetwork *FixedLatencyNetwork) Config() *NoCConfig {
	return fixedLatencyNetwork.config
}

func (fixedLatencyNetwork *FixedLatencyNetwork) Receive(packet Packet) bool {
	var fixedLatency = 20

	fixedLatencyNetwork.Driver().CycleAccurateEventQueue().Schedule(func() {
		packet.OnCompletedCallback()()
	}, fixedLatency)

	return true
}