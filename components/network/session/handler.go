package session

import "github.com/zylikedream/galaxy/components/network/message"

type EventHandler interface {
	OnOpen(Session) error
	OnClose(Session)
	OnError(Session, error)
	OnMessage(Session, *message.Message) error
}

type BaseEventHandler struct {
}

func (e *BaseEventHandler) OnOpen(Session) error {
	return nil
}

func (e *BaseEventHandler) OnClose(Session) {
}

func (e *BaseEventHandler) OnError(Session, error) {
}

func (e *BaseEventHandler) OnMessage(Session, *message.Message) error {
	return nil
}
