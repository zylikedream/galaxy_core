package module

import (
	"time"

	"github.com/zylikedream/galaxy/core/game/gserver/src/cookie"
	"github.com/zylikedream/galaxy/core/game/gserver/src/entity"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/game/proto"
)

type LoginModule struct {
	BaseModule
}

func (l *LoginModule) Handshake(ctx *gscontext.Context, cook *cookie.Cookie, req *proto.ReqHandshake, rsp *proto.RspHandshake) error {
	rsp.Timestamp = uint64(time.Now().Unix())
	return nil
}

func (l *LoginModule) AccountLogin(ctx *gscontext.Context, cook *cookie.Cookie, req *proto.ReqAccountLogin, rsp *proto.RspAccountLogin) error {
	role := entity.NewRoleEntity(0)
	cook.Role = role
	rsp.Create = false
	return nil
}

func init() {
	Register(&LoginModule{})
}
