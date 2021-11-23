package module

import (
	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/gcontext"
)

type LoginModule struct {
	BaseModule
	filters []ModuleFilter
}

func (l *LoginModule) Handshake(ctx gcontext.Context, req *proto.ReqHandshake, rsp *proto.RspHandshake) error {
	return nil
}
