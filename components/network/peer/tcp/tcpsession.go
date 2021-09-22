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
	exitCh chan struct{}
	exit   int32
}

func NewTcpSession(conn net.Conn) *TcpSession {
	return &TcpSession{
		conn: conn,
	}
}

func (t *TcpSession) Start() {
	go t.recvLoop()
}

func (t *TcpSession) recvLoop() {
	pktCodec := t.proc.PktCodec
	msgCodec := t.proc.MsgCodec
	for {
		sizebuf, err := io.ReadAll(io.LimitReader(t.conn, int64(pktCodec.MsgLenLength())))
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
		// todo handle msg
	}
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
}

func (t *TcpSession) Close() {
	if atomic.LoadInt32(&t.exit) == 1 {
		return
	}
	atomic.AddInt32(&t.exit, 1)
	t.exitCh <- struct{}{}
}

func (t *TcpSession) waitExit() {
	<-t.exitCh
	close(t.sendCh)
}
