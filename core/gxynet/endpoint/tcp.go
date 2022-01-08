package endpoint

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/core/gxynet/logger"
	"github.com/zylikedream/galaxy/core/gxynet/message"
	"go.uber.org/zap"
)

type TcpEndpoint struct {
	CoreBundle
	conn   net.Conn
	sendCh chan interface{}
	exit   int32
	data   interface{}
}

func NewTcpEndPoint(conn net.Conn, bundle CoreBundle) *TcpEndpoint {
	return &TcpEndpoint{
		conn:       conn,
		CoreBundle: bundle,
		sendCh:     make(chan interface{}, 64),
	}
}

func (t *TcpEndpoint) Start(ctx context.Context) error {
	err := t.Handler.OnOpen(ctx, t)
	if err != nil {
		return errors.Wrap(err, "on open error")
	}
	go t.sendLoop(ctx)
	t.recvLoop(ctx)
	return nil
}

func (t *TcpEndpoint) recvLoop(ctx context.Context) {
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
			if errors.Is(err, io.EOF) { // 对方主动断开
				err = nil
				logger.Nlog.Debug("remote closed")
				break
			}
			// 出错了
			break
		}
		buf.Write(data[:n])
		pkgLen, msg, err = t.Decode(buf.Bytes())
		if err != nil {
			break
		}
		if msg != nil {
			buf.Next(int(pkgLen))
			if err = t.Handler.OnMessage(ctx, t, msg); err != nil {
				break
			}
		}
	}
	// 被动断开。出错或者对方关闭
	t.Close(ctx, errors.Wrap(err, "recv error"))
}

func (t *TcpEndpoint) IsClosed() bool {
	return atomic.LoadInt32(&t.exit) == 1
}

func (t *TcpEndpoint) Send(msg interface{}) error {
	if t.IsClosed() {
		return fmt.Errorf("conn closed")
	}
	t.sendCh <- msg
	return nil
}

func (t *TcpEndpoint) sendMsg(msg interface{}) error {
	data, err := t.Encode(msg)
	if err != nil {
		return err
	}
	if _, err = t.conn.Write(data); err != nil {
		return err
	}
	return nil
}

func (t *TcpEndpoint) sendLoop(ctx context.Context) {
	var err error
	for rawMsg := range t.sendCh {
		if err = t.sendMsg(rawMsg); err != nil {
			t.Close(ctx, errors.Wrap(err, "send error"))
			continue
		}
	}
	// 关闭整个连接
	t.conn.Close()
	// 这儿才是真正的关闭流程结束
	t.Handler.OnClose(ctx, t)
}

func (t *TcpEndpoint) Close(ctx context.Context, err error) {
	if atomic.LoadInt32(&t.exit) == 1 {
		return
	}
	atomic.AddInt32(&t.exit, 1)
	if err != nil { // 发生错误肯定会调用close
		t.Handler.OnError(ctx, t, err)
		logger.Nlog.Error("tcpendpoint close", zap.Error(err))
	}
	tcpConn := t.conn.(*net.TCPConn)
	_ = tcpConn.CloseRead()
	close(t.sendCh)
}

func (t *TcpEndpoint) Conn() net.Conn {
	return t.conn
}

func (t *TcpEndpoint) GetData() interface{} {
	return t.data
}

func (t *TcpEndpoint) SetData(d interface{}) {
	t.data = d
}
