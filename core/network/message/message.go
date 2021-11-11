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
	MESSAGE_JSON     = "json"
	MESSAGE_PROTOBUF = "protobuf"
)

func NewMessageCodec(t string, c *gconfig.Configuration) (MessageCodec, error) {
	if node, err := gregister.NewNode(t, c.WithPrefix("message")); err != nil {
		return nil, err
	} else {
		return node.(MessageCodec), nil
	}
}
