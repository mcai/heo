package noc

type SelectionAlgorithm interface {
	Select(packet Packet, ivc int, directions []Direction) Direction
}
