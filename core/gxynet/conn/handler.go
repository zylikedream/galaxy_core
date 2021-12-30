package conn

import (
	"context"

	"github.com/zylikedream/galaxy/core/gxynet/message"
)

type EventHandler interface {
	OnOpen(context.Context, Conn) error
	OnClose(context.Context, Conn)
	OnError(context.Context, Conn, error)
	OnMessage(context.Context, Conn, *message.Message) error
}

type BaseEventHandler struct {
}

func (e *BaseEventHandler) OnOpen(context.Context, Conn) error {
	return nil
}

func (e *BaseEventHandler) OnClose(context.Context, Conn) {
}

func (e *BaseEventHandler) OnError(context.Context, Conn, error) {
}

func (e *BaseEventHandler) OnMessage(context.Context, Conn, *message.Message) error {
	return nil
}
