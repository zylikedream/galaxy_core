package module

import (
	"context"
	"time"

	"github.com/zylikedream/galaxy/core/game/gserver/src/entity"
	"github.com/zylikedream/galaxy/core/game/proto"
)

type LoginModule struct {
	BaseModule
}

func (l *LoginModule) Handshake(ctx context.Context, req *proto.ReqHandshake, rsp *proto.RspHandshake) error {
	rsp.Timestamp = uint64(time.Now().Unix())
	return nil
}

func (l *LoginModule) AccountLogin(ctx context.Context, req *proto.ReqAccountLogin, rsp *proto.RspAccountLogin) error {
	role := entity.NewRoleEntity(0)
	rsp.Create = false
	return nil
}

func init() {
	Register(&LoginModule{})
}
