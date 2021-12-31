package main

import (
	"context"

	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/gxynet"
	"github.com/zylikedream/galaxy/core/gxynet/endpoint"
	"github.com/zylikedream/galaxy/core/gxynet/example/echo/proto"
	"github.com/zylikedream/galaxy/core/gxynet/message"
	"go.uber.org/zap"
)

func main() {
	EchoServer()
}

type EchoEventHandler struct {
	endpoint.BaseEventHandler
}

func (e *EchoEventHandler) OnOpen(ctx context.Context, ep endpoint.Endpoint) error {
	gxylog.Infof("conn open, addr=%s", ep.Conn().RemoteAddr())
	return nil
}

func (e *EchoEventHandler) OnClose(ctx context.Context, ep endpoint.Endpoint) {
	gxylog.Infof("conn close, addr=%s", ep.Conn().RemoteAddr())
}

func (e *EchoEventHandler) OnMessage(ctx context.Context, ep endpoint.Endpoint, msg *message.Message) error {
	switch m := msg.Msg.(type) {
	case *proto.EchoReq:
		gxylog.Infof("recv echo req:%v", m)
		ep.Send(&proto.EchoAck{
			Code: 0,
			Msg:  m.Msg,
		})
	}
	return nil
}

func EchoServer() {
	p, err := gxynet.NewNetwork(gxyconfig.New("config/config.toml"))
	if err != nil {
		gxylog.Error("gxynet", zap.Namespace("new failed"), zap.Error(err))
		return
	}
	if err := p.Start(context.Background(), &EchoEventHandler{}); err != nil {
		gxylog.Error("gxynet", zap.Namespace("start failed"), zap.Error(err))
		return
	}
}
