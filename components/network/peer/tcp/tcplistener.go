package tcp

import (
	"net"
	"time"

	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/packet"
	"github.com/zylikedream/galaxy/components/network/peer"
)

type TcpListener struct {
	peer.Processor
	addr     string
	listener net.Listener
}

func newTcpListener(pktCodec packet.PacketCodec, msgCodec message.MessageCodec) *TcpListener {
	return &TcpListener{
		Processor: peer.Processor{
			PktCodec: pktCodec,
			MsgCodec: msgCodec,
		},
	}

}

func (t *TcpListener) Init() error {
	return nil
}

func (t *TcpListener) Start() error {
	var err error
	t.listener, err = net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	go t.accept()
	return nil
}

func (t *TcpListener) accept() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				time.Sleep(time.Millisecond)
				continue
			}
			break
		}
		sess := NewTcpSession(conn)
		go sess.Start()
	}
}

func (t *TcpListener) Type() int {
	return peer.PEER_TCP_ACCEPTOR
}

func (t *TcpListener) Stop() {

}
