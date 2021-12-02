package main

import (
	"context"
	"time"

	"github.com/zylikedream/galaxy/core/game/proto"
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

func (c *Client) Run(ctx context.Context) error {
	if err := c.p.Start(ctx, c); err != nil {
		return err
	}
	return nil
}

func (c *Client) OnOpen(ctx context.Context, sess session.Session) error {
	c.sess = sess
	go c.Work()
	return nil
}

func (c *Client) OnClose(context.Context, session.Session) {
}

func (c *Client) OnError(context.Context, session.Session, error) {
}

func (c *Client) OnMessage(ctx context.Context, sess session.Session, m *message.Message) error {
	switch v := m.Msg.(type) {
	case *proto.Ack:
		ack := v
		meta := message.MessageMetaByID(ack.MsgID)
		if ack.Code != proto.ACK_CODE_OK {
			glog.Error("ack failed", zap.String("msg", meta.TypeName()), zap.String("reason", ack.Reason))
			return nil
		}
		msg := meta.NewInstance()
		if err := c.p.GetMessageCodec().Decode(msg, ack.Data); err != nil {
			return err
		}
		glog.Debug("ack success:", zap.String("name", meta.TypeName()), zap.Any("msg", msg))
	default:
		glog.Debug("recv msg:", zap.Any("msg", m.Msg))
	}
	return nil
}

func main() {
	ctx := context.Background()
	c := NewClient()
	if err := c.Run(ctx); err != nil {
		glog.Error("client run err", zap.Error(err))
	}
}
