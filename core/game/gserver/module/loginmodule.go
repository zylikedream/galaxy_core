package module

import (
	"time"

	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/gcontext"
)

type LoginModule struct {
	BaseModule
}

func (l *LoginModule) Handshake(ctx gcontext.Context, req *proto.ReqHandshake, rsp *proto.RspHandshake) error {
	rsp.Timestamp = uint64(time.Now().Unix())
	return nil
}

func (l *LoginModule) AccountLogin(ctx gcontext.Context, req *proto.ReqAccountLogin, rsp *proto.RspAccountLogin) error {
	rsp.Create = false
	return nil
}

func init() {
	Register(&LoginModule{})
}
