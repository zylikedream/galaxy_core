package module

import (
	"context"

	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/glog"
	"go.uber.org/zap"
)

type LoginModule struct {
	BaseModule
}

func (l *LoginModule) Handshake(ctx context.Context, rsp *proto.RspHandshake) error {
	glog.Debug("recv handshake ", zap.Any("data", rsp))
	return nil
}

func (l *LoginModule) AccountLogin(ctx context.Context, rsp *proto.RspAccountLogin) error {
	glog.Debug("recv accountLogin ", zap.Any("data", rsp))
	return nil
}

func init() {
	Register(&LoginModule{})
}
