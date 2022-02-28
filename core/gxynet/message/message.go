package message

import (
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type MessageCodec interface {
	Decode(msg interface{}, data []byte) error
	Encode(packet interface{}) ([]byte, error)
	Type() string
}

type Message struct {
	Key     string
	Type    uint64
	Payload []byte
	Msg     interface{}
}

type messageOption = func(m *Message)

func WithKey(key string) messageOption {
	return func(m *Message) {
		m.Msg = key
	}
}

func WithType(t int) messageOption {
	return func(m *Message) {
		m.Msg = t
	}
}

func NewMessage(raw interface{}, opts ...messageOption) *Message {
	msg := &Message{
		Msg: raw,
	}
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

const (
	MESSAGE_JSON     = "message.json"
	MESSAGE_PROTOBUF = "message.protobuf"
)

func NewMessageCodec(t string, c *gxyconfig.Configuration) (MessageCodec, error) {
	if node, err := gxyregister.NewNode("message."+t, c); err != nil {
		return nil, err
	} else {
		return node.(MessageCodec), nil
	}
}
