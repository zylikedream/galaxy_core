package tcp

import (
	"io"
	"net"
	"sync/atomic"

	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/peer"
)

type TcpSession struct {
	proc   peer.Processor
	conn   net.Conn
	sendCh chan interface{}
	exit   int32
}

func NewTcpSession(conn net.Conn) *TcpSession {
	return &TcpSession{
		conn: conn,
	}
}

func (t *TcpSession) Start() {
	go t.recvLoop()
	go t.sendLoop()
}

func (t *TcpSession) recvLoop() {
	pktCodec := t.proc.PktCodec
	msgCodec := t.proc.MsgCodec
	for {
		sizebuf, err := io.ReadAll(io.LimitReader(t.conn, int64(pktCodec.MsgLenLength())))
		if err != nil {
			netErr, ok := err.(*net.OpError)
			if ok && netErr.Err == net.ErrClosed { // 主动断开
				return
			} else {
				break
			}
		}
		// eof
		if len(sizebuf) == 0 {
			break
		}
		size := pktCodec.Uint(sizebuf)
		data, err := io.ReadAll(io.LimitReader(t.conn, int64(size)))
		if err != nil {
			break
		}
		if len(data) < int(size) {
			break
		}
		msg, err := pktCodec.Decode(data)
		if err != nil {
			break
		}
		msg.Msg, err = msgCodec.Decode(msg.ID, msg.Payload)
		if err != nil {
			break
		}
		msg.Sess = t
		// todo handle msg
	}
	// 被动断开，出错或者对方关闭
	t.passiveClose()

}

func (t *TcpSession) Send(msg interface{}) error {
	t.sendCh <- msg
	return nil
}

func (t *TcpSession) sendLoop() {
	pktCodec := t.proc.PktCodec
	msgCodec := t.proc.MsgCodec
	var err error
	for rawMsg := range t.sendCh {
		msg := &message.Message{
			Msg: rawMsg,
		}
		msg.ID, msg.Payload, err = msgCodec.Encode(rawMsg)
		if err != nil {
			break
		}
		data, err := pktCodec.Encode(msg)
		if err != nil {
			break
		}
		_, err = t.conn.Write(data)
		if err != nil {
			break
		}
	}
	// 关闭整个连接
	t.conn.Close()
}

func (t *TcpSession) Close() {
	if atomic.LoadInt32(&t.exit) == 1 {
		return
	}
	atomic.AddInt32(&t.exit, 1)
	tcpConn := t.conn.(*net.TCPConn)
	tcpConn.CloseRead()
	close(t.sendCh)
}

func (t *TcpSession) passiveClose() {
	if atomic.LoadInt32(&t.exit) == 1 {
		return
	}
	atomic.AddInt32(&t.exit, 1)
	close(t.sendCh)
}
