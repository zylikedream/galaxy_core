package main

import (
	"context"
	"time"

	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/gxyrpc"
	"github.com/zylikedream/galaxy/core/gxyrpc/example/proto"
)

func main() {
	cli, err := gxyrpc.NewGrpcClient("config/config.etcd.toml")
	if err != nil {
		panic(err)
	}
	req := &proto.MulRequest{
		A: 94545,
		B: 7824,
	}
	reply := &proto.MulReply{}
	for {
		if err := cli.Call(context.Background(), "Arith", "Mul", req, reply); err != nil {
			panic(err)
		}
		gxylog.Debugf("%d * %d = %d", reply.A, reply.B, reply.Result)
		time.Sleep(time.Second * 3)
	}
}
