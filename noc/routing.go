package noc

type RoutingAlgorithm interface {
	NextHop(packet Packet, parent int) []Direction
}