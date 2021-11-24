package module

import (
	"github.com/zylikedream/galaxy/core/game/gserver/define"
	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/network/session"
)

func GetSessionFromCtx(ctx gcontext.Context) session.Session {
	return ctx.Value(define.SessionCtxKey).(session.Session)
}
