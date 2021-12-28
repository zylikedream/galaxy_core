package module

import (
	"time"

	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/core/game/gserver/src/cookie"
	"github.com/zylikedream/galaxy/core/game/gserver/src/entity"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/glog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type LoginModule struct {
	BaseModule
	logger *glog.GalaxyLog
}

func (l *LoginModule) Init(ctx *gscontext.Context) error {
	logger := ctx.GetLogger()
	l.logger = logger.With(zap.Namespace("login"))
	return nil
}

func (l *LoginModule) Handshake(ctx *gscontext.Context, cook *cookie.Cookie, req *proto.ReqHandshake, rsp *proto.RspHandshake) error {
	rsp.Timestamp = uint64(time.Now().Unix())
	return nil
}

func (l *LoginModule) AccountLogin(ctx *gscontext.Context, cook *cookie.Cookie, req *proto.ReqAccountLogin, rsp *proto.RspAccountLogin) error {
	role := entity.NewRoleEntity()
	var newRole bool
	if err := role.LoadByAccount(ctx, req.Account); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			newRole = true
			// 新号
		} else {
			return errors.Wrap(err, "load role failed")
		}
	}
	if newRole {
		if err := role.Create(ctx, req.Account); err != nil {
			return errors.Wrap(err, "create role faield")
		}
	}
	cook.Role = role
	rsp.Create = newRole
	return nil
}

func init() {
	Register(&LoginModule{})
}
