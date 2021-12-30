/*
 * @Author: your name
 * @Date: 2021-11-04 17:39:40
 * @LastEditTime: 2021-11-05 15:58:53
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /components/gxynet/example/echo/client/config/main.go
 */
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/gxynet"
	"github.com/zylikedream/galaxy/core/gxynet/example/echo/proto"
	"github.com/zylikedream/galaxy/core/gxynet/message"
	"github.com/zylikedream/galaxy/core/gxynet/session"
	"go.uber.org/zap"
)

var wg sync.WaitGroup

func main() {
	EchoClient()
}

type EchoEventHandler struct {
	session.BaseEventHandler
}

func (e *EchoEventHandler) OnOpen(ctx context.Context, sess session.Session) error {
	gxylog.Infof("session open, addr=%s", sess.Conn().RemoteAddr())
	go run(sess)
	return nil
}

func (e *EchoEventHandler) OnClose(ctx context.Context, sess session.Session) {
	gxylog.Infof("session close, addr=%s", sess.Conn().RemoteAddr())
}

func (e *EchoEventHandler) OnMessage(ctx context.Context, sess session.Session, msg *message.Message) error {
	switch m := msg.Msg.(type) {
	case *proto.EchoAck:
		gxylog.Infof("recv message:%v", m)
		sess.Send(&proto.EchoAck{
			Code: 0,
			Msg:  m.Msg,
		})
	}
	return nil
}

func run(sess session.Session) {
	var i int
	for {
		msg := &proto.EchoReq{
			Msg: fmt.Sprintf("hello %d", i),
		}
		if err := sess.Send(msg); err != nil {
			gxylog.Error("send error", zap.Error(err))
			break
		}
		i++
		time.Sleep(time.Second * 5)
	}
}

func EchoClient() {
	p, err := gxynet.NewNetwork(gxyconfig.New("config/gxynet.toml"))
	if err != nil {
		gxylog.Error("gxynet", zap.Namespace("new failed"), zap.Error(err))
		return
	}
	if err := p.Start(context.Background(), &EchoEventHandler{}); err != nil {
		gxylog.Error("gxynet", zap.Namespace("start failed"), zap.Error(err))
		return
	}
}
