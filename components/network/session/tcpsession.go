/*
 * @Author: your name
 * @Date: 2021-10-19 17:41:17
 * @LastEditTime: 2021-11-04 17:20:53
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /components/network/session/tcpsession.go
 */
package session

import (
	"io"
	"net"
	"sync/atomic"

	"github.com/zylikedream/galaxy/components/network/logger"
	"github.com/zylikedream/galaxy/components/network/message"
	"go.uber.org/zap"
)

type TcpSession struct {
	SessionBundle
	conn   net.Conn
	sendCh chan interface{}
	exit   int32
}

func NewTcpSession(conn net.Conn, bundle SessionBundle) *TcpSession {
	return &TcpSession{
		conn:          conn,
		SessionBundle: bundle,
		sendCh:        make(chan interface{}, 64),
	}
}

func (t *TcpSession) Start() {
	if err := t.Handler.OnOpen(t); err != nil {
		logger.Nlog.Error("tcpsession", zap.String("msg", "session On Open"), zap.Error(err))
		t.Close()
		return
	}
	go t.recvLoop()
	go t.sendLoop()
}

func (t *TcpSession) recvLoop() {
	var err error
	var msg *message.Message
	for {
		msg, err = t.Proc.ReadAndDecode(t.conn)
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
		if err = t.Handler.OnMessage(t, msg); err != nil {
			break
		}
	}
	// 被动断开。出错或者对方关闭
	logger.Nlog.Error("tcpsession", zap.String("msg", "recv loop break"), zap.Error(err))
	t.Handler.OnError(t, err)
	t.passiveClose()
}

func (t *TcpSession) Send(msg interface{}) error {
	t.sendCh <- msg
	return nil
}

func (t *TcpSession) sendLoop() {
	var err error
	var data []byte
	for rawMsg := range t.sendCh {
		data, err = t.Proc.Encode(rawMsg)
		if err != nil {
			break
		}
		_, err = t.conn.Write(data)
		if err != nil {
			break
		}
	}
	if err != nil {
		t.Handler.OnError(t, err)
		logger.Nlog.Error("tcpsession", zap.String("msg", "send loop break"), zap.Error(err))
	}
	// 关闭整个连接
	t.conn.Close()
	t.Handler.OnClose(t)
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

func (t *TcpSession) Conn() net.Conn {
	return t.conn
}
