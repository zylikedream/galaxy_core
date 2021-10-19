package tcp

import (
	"io"
	"net"
	"sync/atomic"

	"github.com/zylikedream/galaxy/components/network/peer/processor"
)

type TcpSession struct {
	proc   *processor.Processor
	conn   net.Conn
	sendCh chan interface{}
	exit   int32
}

func NewTcpSession(conn net.Conn, proc *processor.Processor) *TcpSession {
	return &TcpSession{
		conn: conn,
		proc: proc,
	}
}

func (t *TcpSession) Start() {
	go t.recvLoop()
	go t.sendLoop()
}

func (t *TcpSession) recvLoop() {
	for {
		msg, err := t.proc.ReadAndDecode(t.conn)
		if err != nil {
			netErr, ok := err.(*net.OpError)
			if ok && netErr.Err == net.ErrClosed { // 主动断开 不执行断开逻辑，已经断开
				return
			}
			if err == io.EOF { // 对方已关闭
				break
			}
			// 出错了
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
	for rawMsg := range t.sendCh {
		data, err := t.proc.Encode(rawMsg)
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
	_ = tcpConn.CloseRead()
	close(t.sendCh)
}

func (t *TcpSession) passiveClose() {
	if atomic.LoadInt32(&t.exit) == 1 {
		return
	}
	atomic.AddInt32(&t.exit, 1)
	close(t.sendCh)
}
