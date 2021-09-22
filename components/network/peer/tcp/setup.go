package tcp

import (
	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/packet"
	"github.com/zylikedream/galaxy/components/network/peer"
)

func init() {
	p := &TcpListener{}
	peer.Register(p.Type(), func(pc packet.PacketCodec, mc message.MessageCodec) peer.Peer {
		return newTcpListener(pc, mc)
	})
}
