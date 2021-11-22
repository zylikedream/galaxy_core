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

func (l *LogicHandle) OnOpen(session.Session) error {
	return nil
}

func (l *LogicHandle) OnClose(session.Session) {
}

func (l *LogicHandle) OnError(session.Session, error) {
}

func (l *LogicHandle) OnMessage(sess session.Session, m *message.Message) error {
	if err := module.HandleMessage(context.Background(), m.Msg); err != nil {
		glog.Error("handle message error", zap.Error(err))
	}
	return nil
}
