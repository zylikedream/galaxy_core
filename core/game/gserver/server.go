package main

import (
	"context"

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

func NewServer() *Server {
	svr := &Server{}
	p, err := network.NewNetwork("config/network.toml")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	svr.p = p
	svr.logger = glog.NewLogger("server", "config/log.toml")
	glog.SetDefaultLogger(svr.logger)

	cli, err := gmongo.NewMongoClient(ctx, "config/mongo.toml")
	if err != nil {
		panic(err)
	}
	svr.mgoCli = cli
	return svr
}

func (s *Server) Run() error {
	if err := s.p.Start(&logic.LogicHandle{}); err != nil {
		return err
	}
	return nil
}

func main() {
	s := NewServer()
	if err := s.Run(); err != nil {
		glog.Error("server run err", zap.Error(err))
	}
}
