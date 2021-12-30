package logic

import (
	"context"

	"github.com/zylikedream/galaxy/core/game/gserver/src/cookie"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/game/gserver/src/module"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/network/message"
	"github.com/zylikedream/galaxy/core/network/session"
	"go.uber.org/zap"
)

type LogicHandle struct {
	session.BaseEventHandler
}

func (l *LogicHandle) OnOpen(ctx context.Context, sess session.Session) error {
	sess.SetData(cookie.NewCookie())
	gsctx := ctx.(*gscontext.Context)
	gsctx.SetSession(sess)
	return nil
}

func (l *LogicHandle) OnClose(context.Context, session.Session) {
}

func (l *LogicHandle) OnError(context.Context, session.Session, error) {
}

func (l *LogicHandle) OnMessage(ctx context.Context, sess session.Session, m *message.Message) error {
	gsctx := ctx.(*gscontext.Context)
	cook := sess.GetData().(*cookie.Cookie)
	if err := module.HandleMessage(gsctx, cook, m.Msg); err != nil {
		gxylog.Error("handle message error", zap.Error(err))
	}
	return nil
}
