package main

import (
	"context"
	"time"

	"github.com/zylikedream/galaxy/core/game/gclient/define"
	"github.com/zylikedream/galaxy/core/game/gclient/module"
	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/network"
	"github.com/zylikedream/galaxy/core/network/message"
	"github.com/zylikedream/galaxy/core/network/peer"
	"github.com/zylikedream/galaxy/core/network/session"
	"go.uber.org/zap"
)

type Client struct {
	p peer.Peer
	session.BaseEventHandler
	sess session.Session
}

func NewClient() *Client {
	cli := &Client{}
	p, err := network.NewNetwork("config/network.toml")
	if err != nil {
		panic(err)
	}
	cli.p = p
	logger := glog.NewLogger("client", "config/log.toml")
	glog.SetDefaultLogger(logger)
	return cli
}

func (c *Client) Work() {
	c.send(&proto.ReqHandshake{
		LoginKey: "golang client",
	})
	time.Sleep(time.Second)
	c.send(&proto.ReqAccountLogin{
		Account:    "zhangyi",
		ClientInfo: proto.PClientInfo{},
	})
}

func (c *Client) send(msg interface{}) error {
	return c.sess.Send(msg)
}

func (c *Client) Run() error {
	if err := c.p.Start(c); err != nil {
		return err
	}
	return nil
}

func (c *Client) OnOpen(ctx context.Context, sess session.Session) error {
	gctx := ctx.(*gcontext.Context)
	c.sess = sess
	gctx.SetValue(define.SessionCtxKey, sess)
	go c.Work()
	return nil
}

func (c *Client) OnClose(context.Context, session.Session) {
}

func (c *Client) OnError(context.Context, session.Session, error) {
}

func (c *Client) OnMessage(ctx context.Context, sess session.Session, m *message.Message) error {
	if err := module.HandleMessage(ctx, m.Msg); err != nil {
		glog.Error("handle message error", zap.Error(err))
	}
	return nil
}

func main() {
	c := NewClient()
	if err := c.Run(); err != nil {
		glog.Error("client run err", zap.Error(err))
	}
}
