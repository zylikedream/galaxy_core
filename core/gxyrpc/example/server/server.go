package main

import (
	"github.com/zylikedream/galaxy/core/gxyrpc/example/proto"
	"github.com/zylikedream/galaxy/core/gxyrpc/server"
)

func main() {
	s, err := server.NewGrpcServer("config/config.etcd.toml")
	if err != nil {
		panic(err)
	}
	if err := s.ReigsterService(new(proto.Arith)); err != nil {
		panic(err)
	}
	if err := s.Start(); err != nil {
		panic(err)
	}
}
