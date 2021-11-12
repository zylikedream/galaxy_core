package proto

import (
	"context"

	"github.com/zylikedream/galaxy/core/glog"
)

type MulRequest struct {
	A int
	B int
}

type MulReply struct {
	A      int
	B      int
	Result int
}

type Arith struct {
}

func (a *Arith) Mul(ctx context.Context, req *MulRequest, reply *MulReply) error {
	reply.A = req.A
	reply.B = req.B
	reply.Result = req.A * req.B
	glog.Debugf("call %d*%d=%d", reply.A, reply.B, reply.Result)
	return nil
}

type ArithService struct {
}

func (a *ArithService) Service() interface{} {
	return new(Arith)
}

func (a *ArithService) Name() string {
	return "Arith"
}

func (a *ArithService) Meta() string {
	return ""
}
