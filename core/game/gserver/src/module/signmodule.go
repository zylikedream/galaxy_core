// 签到
// 玩家上线自动签到当天
// 玩家需要手动领取(可领取多天)
// 玩家为上线的天数可以补签
// 签到一定天数可以领取累积奖励
package module

import (
	"errors"
	"time"

	"github.com/zylikedream/galaxy/core/game/gserver/src/cookie"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/gxylog"
	"go.uber.org/zap"
)

type SignModule struct {
	BaseModule
	logger *gxylog.GalaxyLog
}

func (l *SignModule) Init(ctx *gscontext.Context) error {
	logger := ctx.GetLogger()
	l.logger = logger.With(zap.Namespace("sign"))
	return nil
}

func (l *SignModule) reqSignInfo(ctx *gscontext.Context, cook *cookie.Cookie, req *proto.ReqSignInfo, rsp *proto.RspSignInfo) error {
	rsign := cook.Role.Sign
	rsp.SignDay = rsign.SignDay
	rsp.SignTime = rsign.SignTime
	return nil
}

func (l *SignModule) reqSignSign(ctx *gscontext.Context, cook *cookie.Cookie, req *proto.ReqSignSign, rsp *proto.RspSignSign) error {
	// check
	sign := cook.Role.Sign
	if sign.SignTime > 0 {
		return errors.New("already signed")
	}

	sign.SignTime = time.Now().Unix()
	rewards := ctx.GetGameConfig().TbSign.Get(int32(sign.SignDay)).Rewards
	if err := cook.Role.Bag.AddItemRc(ctx, rewards); err != nil {
		return err
	}

	// trigger
	return nil
}

func init() {
	Register(&SignModule{})
}
