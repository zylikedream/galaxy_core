/*
 * @Author: your name
 * @Date: 2021-11-04 17:39:40
 * @LastEditTime: 2021-11-05 15:58:53
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /components/network/example/echo/client/config/main.go
 */
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/network"
	"github.com/zylikedream/galaxy/core/network/example/echo/proto"
	"github.com/zylikedream/galaxy/core/network/message"
	"github.com/zylikedream/galaxy/core/network/session"
	"go.uber.org/zap"
)

var wg sync.WaitGroup

func main() {
	EchoClient()
}

type EchoEventHandler struct {
	session.BaseEventHandler
}

func (e *EchoEventHandler) OnOpen(ctx *gcontext.Context, sess session.Session) error {
	glog.Infof("session open, addr=%s", sess.Conn().RemoteAddr())
	go run(sess)
	return nil
}

func (e *EchoEventHandler) OnClose(ctx *gcontext.Context, sess session.Session) {
	glog.Infof("session close, addr=%s", sess.Conn().RemoteAddr())
}

func (e *EchoEventHandler) OnMessage(ctx *gcontext.Context, sess session.Session, msg *message.Message) error {
	switch m := msg.Msg.(type) {
	case *proto.EchoAck:
		glog.Infof("recv message:%v", m)
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
			glog.Error("send error", zap.Error(err))
			break
		}
		i++
		time.Sleep(time.Second * 5)
	}
}

func EchoClient() {
	p, err := network.NewNetwork("config/config.toml")
	if err != nil {
		glog.Error("network", zap.Namespace("new failed"), zap.Error(err))
		return
	}
	if err := p.Start(&EchoEventHandler{}); err != nil {
		glog.Error("network", zap.Namespace("start failed"), zap.Error(err))
		return
	}
}
