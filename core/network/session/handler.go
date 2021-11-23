package session

import (
	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/network/message"
)

type EventHandler interface {
	OnOpen(gcontext.Context, Session) error
	OnClose(gcontext.Context, Session)
	OnError(gcontext.Context, Session, error)
	OnMessage(gcontext.Context, Session, *message.Message) error
}

type BaseEventHandler struct {
}

func (e *BaseEventHandler) OnOpen(gcontext.Context, Session) error {
	return nil
}

func (e *BaseEventHandler) OnClose(gcontext.Context, Session) {
}

func (e *BaseEventHandler) OnError(gcontext.Context, Session, error) {
}

func (e *BaseEventHandler) OnMessage(gcontext.Context, Session, *message.Message) error {
	return nil
}
