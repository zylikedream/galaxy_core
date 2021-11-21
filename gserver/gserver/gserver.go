package gserver

import (
	"github.com/zylikedream/galaxy/core/network"
	"google.golang.org/grpc/peer"
)

type gserverConfig struct {
	Services []string `toml:"services"`
}

type Server struct {
	pr *peer.Peer
}

func NewServer() *Server {
	pr := network.NewNetwork()
}
