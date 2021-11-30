package module

import (
	"context"

	"github.com/zylikedream/galaxy/core/game/gclient/define"
	"github.com/zylikedream/galaxy/core/network/session"
)

func GetSessionFromCtx(ctx context.Context) session.Session {
	return ctx.Value(define.SessionCtxKey).(session.Session)
}
