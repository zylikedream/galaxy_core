/*
 * @Author: your name
 * @Date: 2021-11-04 17:39:40
 * @LastEditTime: 2021-11-04 17:44:32
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /components/network/example/echo/client/config/main.go
 */
package main

import (
	"github.com/zylikedream/galaxy/components/glog"
	"github.com/zylikedream/galaxy/components/network"
	"github.com/zylikedream/galaxy/components/network/example/echo/proto"
	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/session"
	"go.uber.org/zap"
)

func main() {
	EchoClient()
}

var gsess session.Session

type EchoEventHandler struct {
	session.BaseEventHandler
}

func (e *EchoEventHandler) OnOpen(sess session.Session) error {
	glog.Infof("session open, addr=%s", sess.Conn().RemoteAddr())
	gsess = sess
	return nil
}

func (e *EchoEventHandler) OnClose(sess session.Session) {
	glog.Infof("session close, addr=%s", sess.Conn().RemoteAddr())
}

func (e *EchoEventHandler) OnMessage(sess session.Session, msg *message.Message) error {
	switch m := msg.Msg.(type) {
	case *proto.EchoReq:
		glog.Infof("recv message:%v", msg)
		sess.Send(&proto.EchoAck{
			Code: 0,
			Msg:  m.Msg,
		})
	}
	return nil
}

func EchoClient() {
	p, err := network.NewNetwork("config/config.toml")
	if err != nil {
		glog.Error("network", zap.Namespace("new failed"), zap.Error(err))
		return
	}
	p.Start(&EchoEventHandler{})
}
