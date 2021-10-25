package session

import (
	"github.com/zylikedream/galaxy/components/network/handler"
	"github.com/zylikedream/galaxy/components/network/message"
)

type Session interface {
	Send(msg interface{}) error
	BindHandler(h handler.Handler)
	Close()
}

type BaseSession struct {
	handler handler.Handler
}

func (s *BaseSession) BindHandler(h handler.Handler) {
	s.handler = h
}

func (s *BaseSession) handle(msg *message.Message) error {
	return s.handler.Handle(msg)
}
