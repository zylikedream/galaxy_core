package message

import (
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
)

type MessageCodec interface {
	Decode(msg interface{}, data []byte) error
	Encode(packet interface{}) ([]byte, error)
	Type() string
}

type Message struct {
	ID      int
	Type    uint64
	Payload []byte
	Msg     interface{}
}

const (
	MESSAGE_JSON     = "message.json"
	MESSAGE_PROTOBUF = "message.protobuf"
)

func NewMessageCodec(t string, c *gconfig.Configuration) (MessageCodec, error) {
	if node, err := gregister.NewNode("message."+t, c); err != nil {
		return nil, err
	} else {
		return node.(MessageCodec), nil
	}
}
