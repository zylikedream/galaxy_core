package peer

import (
	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/packet"
)

type processor struct {
	pktCodec packet.PacketCodec
	msgCodec message.MessageCodec
}

type Peer interface {
	Init() error
	Start() error
	Stop()
	Name() string
}
