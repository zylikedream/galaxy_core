package handler

import "github.com/zylikedream/galaxy/components/network/message"

type Handler interface {
	Handle(msg *message.Message) error
}
