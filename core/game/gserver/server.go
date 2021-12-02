package main

import (
	"context"

	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/game/gserver/src/logic"
	_ "github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/gmongo"
	"github.com/zylikedream/galaxy/core/network"
	"github.com/zylikedream/galaxy/core/network/peer"
	"go.uber.org/zap"
)

type Server struct {
	p      peer.Peer
	logger *glog.GalaxyLog
	mgoCli *gmongo.MongoClient
}

func NewServer(ctx *gscontext.Context) *Server {
	svr := &Server{}
	if err := svr.Init(ctx); err != nil {
		panic(err)
	}
	return svr
}

func (s *Server) Init(ctx *gscontext.Context) error {
	p, err := network.NewNetwork("config/network.toml")
	if err != nil {
		return err
	}
	s.p = p

	s.logger = glog.NewLogger("server", "config/log.toml")
	glog.SetDefaultLogger(s.logger)

	cli, err := gmongo.NewMongoClient(ctx, "config/mongo.toml")
	if err != nil {
		return err
	}
	s.mgoCli = cli
	ctx.SetLogger(s.logger)
	ctx.SetMongo(s.mgoCli)
	ctx.SetPeer(s.p)
	return nil

}

func (s *Server) Run(ctx *gscontext.Context) error {
	if err := s.p.Start(ctx, &logic.LogicHandle{}); err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := gscontext.NewContext(context.Background())
	s := NewServer(ctx)
	if err := s.Run(ctx); err != nil {
		glog.Error("server run err", zap.Error(err))
	}
}
