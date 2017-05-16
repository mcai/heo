package noc

type InputVirtualChannel struct {
	InputPort            *InputPort
	Num                  int
	InputBuffer          *InputBuffer
	Route                Direction
	OutputVirtualChannel *OutputVirtualChannel
}

func NewInputVirtualChannel(inputPort *InputPort, num int) *InputVirtualChannel {
	var inputVirtualChannel = &InputVirtualChannel{
		InputPort:inputPort,
		Num:num,
		Route:DIRECTION_UNKNOWN,
	}

	inputVirtualChannel.InputBuffer = NewInputBuffer(inputVirtualChannel)

	return inputVirtualChannel
}
