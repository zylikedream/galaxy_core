package peer

import (
	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/packet"
)

type newPeerFunc = func(packet.PacketCodec, message.MessageCodec) Peer

var peerMap map[int]newPeerFunc = make(map[int]newPeerFunc)

// processor必须包含次结构
type Processor struct {
	PktCodec packet.PacketCodec
	MsgCodec message.MessageCodec
}

const (
	PEER_TCP_ACCEPTOR = iota
)

type Peer interface {
	Init() error
	Start() error
	Stop()
	Type() int
}

func Register(peerType int, nfun newPeerFunc) {
	peerMap[peerType] = nfun
}

func NewPeer(peerType int, pktCodec packet.PacketCodec, msgCodec message.MessageCodec) Peer {
	Func := peerMap[peerType]
	return Func(pktCodec, msgCodec)
}
