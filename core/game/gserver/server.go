package main

import (
	"github.com/zylikedream/galaxy/core/game/gserver/logic"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/network"
	"github.com/zylikedream/galaxy/core/network/peer"
	"github.com/zylikedream/galaxy/core/network/session"
)

type Server struct {
	p peer.Peer
}

func NewServer() *Server {
	svr := &Server{}
	p, err := network.NewNetwork("network.toml")
	if err != nil {
		panic(err)
	}
	svr.p = p
	logger := glog.NewLogger("server", "log.toml")
	glog.SetDefaultLogger(logger)
	return svr
}

func (s *Server) Run() error {
	if err := s.p.Start(&logic.LogicHandle{}); err != nil {
		return err
	}
	return nil
}
