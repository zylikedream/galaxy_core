package logic

import (
	"github.com/zylikedream/galaxy/core/network/message"
	"github.com/zylikedream/galaxy/core/network/session"
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

func (l *LogicHandle) OnMessage(session.Session, *message.Message) error {
	return nil
}
