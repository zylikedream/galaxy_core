package session

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/core/network/logger"
	"github.com/zylikedream/galaxy/core/network/message"
	"go.uber.org/zap"
)

type TcpSession struct {
	SessionBundle
	conn   net.Conn
	sendCh chan interface{}
	ctx    context.Context
	exit   int32
}

func NewTcpSession(conn net.Conn, bundle SessionBundle) *TcpSession {
	return &TcpSession{
		ctx:           context.Background(),
		conn:          conn,
		SessionBundle: bundle,
		sendCh:        make(chan interface{}, 64),
	}
}

func (t *TcpSession) Start() {
	go t.recvLoop()
	go t.sendLoop()
	if err := t.Handler.OnOpen(t.ctx, t); err != nil {
		t.Close(errors.Wrap(err, "on open error"))
		return
	}
}

func (t *TcpSession) recvLoop() {
	var err error
	var msg *message.Message

	buf := &bytes.Buffer{}
	var pkgLen uint64
	data := make([]byte, 1024)
	var n int
	for {
		n, err = t.conn.Read(data)
		if err != nil {
			netErr, ok := err.(*net.OpError)
			if ok && netErr.Err == net.ErrClosed { // 调用close主动断开 已经执行过断开逻辑了 直接返回
				return
			}
			// 出错了
			break
		}
		buf.Write(data[:n])
		pkgLen, msg, err = t.Proc.Decode(buf.Bytes())
		if err != nil {
			break
		}
		if msg != nil {
			buf.Next(int(pkgLen))
			if err = t.Handler.OnMessage(t.ctx, t, msg); err != nil {
				break
			}
		}
	}
	// 被动断开。出错或者对方关闭
	t.Close(errors.Wrap(err, "recv error"))
}

func (t *TcpSession) IsClosed() bool {
	return atomic.LoadInt32(&t.exit) == 1
}

func (t *TcpSession) Send(msg interface{}) error {
	if t.IsClosed() {
		return fmt.Errorf("session closed")
	}
	t.sendCh <- msg
	return nil
}

func (t *TcpSession) sendMsg(msg interface{}) error {
	data, err := t.Proc.Encode(msg)
	if err != nil {
		return err
	}
	if _, err = t.conn.Write(data); err != nil {
		return err
	}
	return nil
}

func (t *TcpSession) sendLoop() {
	var err error
	for rawMsg := range t.sendCh {
		if err = t.sendMsg(rawMsg); err != nil {
			t.Close(errors.Wrap(err, "send error"))
			continue
		}
	}
	// 关闭整个连接
	t.conn.Close()
	// 这儿才是真正的关闭流程结束
	t.Handler.OnClose(t.ctx, t)
}

func (t *TcpSession) Close(err error) {
	if atomic.LoadInt32(&t.exit) == 1 {
		return
	}
	atomic.AddInt32(&t.exit, 1)
	if err != nil { // 发生错误肯定会调用close
		t.Handler.OnError(t.ctx, t, err)
		logger.Nlog.Error("tcpsession close", zap.Error(err))
	}
	tcpConn := t.conn.(*net.TCPConn)
	_ = tcpConn.CloseRead()
	close(t.sendCh)
}

func (t *TcpSession) Conn() net.Conn {
	return t.conn
}
