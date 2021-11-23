package logic

import (
	"context"

	"github.com/zylikedream/galaxy/core/game/gserver/module"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/network/message"
	"github.com/zylikedream/galaxy/core/network/session"
	"go.uber.org/zap"
)

type LogicHandle struct {
	session.BaseEventHandler
}

func (l *LogicHandle) OnOpen(context.Context, session.Session) error {
	return nil
}

func (l *LogicHandle) OnClose(context.Context, session.Session) {
}

func (l *LogicHandle) OnError(context.Context, session.Session, error) {
}

func (l *LogicHandle) OnMessage(ctx context.Context, sess session.Session, m *message.Message) error {
	if err := module.HandleMessage(ctx, m.Msg); err != nil {
		glog.Error("handle message error", zap.Error(err))
	}
	return nil
}
