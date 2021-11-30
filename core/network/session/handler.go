package session

import (
	"context"

	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/network/message"
)

type EventHandler interface {
	OnOpen(*gcontext.Context, Session) error
	OnClose(context.Context, Session)
	OnError(context.Context, Session, error)
	OnMessage(context.Context, Session, *message.Message) error
}

type BaseEventHandler struct {
}

func (e *BaseEventHandler) OnOpen(context.Context, Session) error {
	return nil
}

func (e *BaseEventHandler) OnClose(context.Context, Session) {
}

func (e *BaseEventHandler) OnError(context.Context, Session, error) {
}

func (e *BaseEventHandler) OnMessage(context.Context, Session, *message.Message) error {
	return nil
}
