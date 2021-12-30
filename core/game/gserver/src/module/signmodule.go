// 签到
// 玩家上线自动签到当天
// 玩家需要手动领取(可领取多天)
// 玩家为上线的天数可以补签
// 签到一定天数可以领取累积奖励
package module

import (
	"errors"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/zylikedream/galaxy/core/game/gserver/src/cookie"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/glog"
	"go.uber.org/zap"
)

type SignModule struct {
	BaseModule
	logger *glog.GalaxyLog
}

func (l *SignModule) Init(ctx *gscontext.Context) error {
	logger := ctx.GetLogger()
	l.logger = logger.With(zap.Namespace("sign"))
	return nil
}

func (l *SignModule) reqSignInfo(ctx *gscontext.Context, cook *cookie.Cookie, req *proto.ReqSignInfo, rsp *proto.RspSignInfo) error {
	sign := cook.Role.Sign
	rsp.SignDay = sign.SignDay
	rsp.SignTime = sign.SignTime
	return nil
}

func (l *SignModule) reqSignSign(ctx *gscontext.Context, cook *cookie.Cookie, req *proto.ReqSignSign, rsp *proto.RspSignSign) error {
	// check
	sign := cook.Role.Sign
	if sign.SignTime > 0 {
		return errors.New("already signed")
	}
	// do
	rewards := ctx.GetGameConfig().TbSign.Get(int32(sign.SignDay)).Rewards
	cook.Role.Bag.AddItem(ctx, rewards)

	sign.SignTime = time.Now().Unix()
	// trigger
	return nil
}

func init() {
	Register(&SignModule{})
}
