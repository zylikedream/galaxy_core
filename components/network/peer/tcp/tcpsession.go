package tcp

import (
	"io"
	"net"

	"github.com/zylikedream/galaxy/components/network/peer"
)

type TcpSession struct {
	p      peer.Peer
	conn   net.Conn
	sendCh chan interface{}
	exitCh chan struct{}
	exit   int32
}

func NewTcpSession(conn net.Conn) *TcpSession {
	return &TcpSession{
		conn: conn,
	}
}

func (t *TcpSession) Start() {
	msgCodec := t.p.MsgCodec()
	for {
		sizebuf, err := io.ReadAll(io.LimitReader(t.conn, int64(msgCodec.MsgLenLength())))
		if err != nil {
			if netErr, ok := err.(*net.OpError); ok { // 主动断开
				if netErr.Err == net.ErrClosed {
					break
				}
			} else {
				break
			}
		}
		// eof
		if len(sizebuf) == 0 {
			break
		}
		size := msgCodec.Uint(sizebuf)
		body, err := io.ReadAll(io.LimitReader(t.conn, int64(size)))
		if err != nil {
			break
		}
		if len(body) < int(size) {
			break
		}
		packet, err := msgCodec.Decode(body)
		if err != nil {
			break
		}
		if err := packet.Decode(); err != nil {
			return
		}
	}
}

func (t *TcpSession) recvLoop() {

}

func (t *TcpSession) sendLoop() {

}
