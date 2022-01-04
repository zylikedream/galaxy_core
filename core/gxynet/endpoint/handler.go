package endpoint

import (
	"context"

	"github.com/zylikedream/galaxy/core/gxynet/message"
)

type EventHandler interface {
	OnOpen(context.Context, Endpoint) (context.Context, error)
	OnClose(context.Context, Endpoint)
	OnError(context.Context, Endpoint, error)
	OnMessage(context.Context, Endpoint, *message.Message) error
}

type BaseEventHandler struct {
}

func (e *BaseEventHandler) OnOpen(ctx context.Context, _ Endpoint) (context.Context, error) {
	return ctx, nil
}

func (e *BaseEventHandler) OnClose(context.Context, Endpoint) {
}

func (e *BaseEventHandler) OnError(context.Context, Endpoint, error) {
}

func (e *BaseEventHandler) OnMessage(context.Context, Endpoint, *message.Message) error {
	return nil
}
