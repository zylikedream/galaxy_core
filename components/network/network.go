package network

const (
	PEER_TCP_ACCEPTOR = iota
)

const (
	PACKET_LTIV = iota
)

const (
	MESSAGE_JSON = iota
	MESSAGE_PROTOBUF
)

type Network struct {
}

func NewNetwork(tcpType int, packetType int, codec int) *Network {

}
