package message

import (
	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/gregister"
)

type MessageCodec interface {
	Decode(ID uint64, data []byte) (interface{}, error)
	Encode(packet interface{}) (uint64, []byte, error)
	ReisterPacket(ID uint64, packet interface{}) error
	Type() string
}

type Message struct {
	ID      uint64
	Type    uint64
	Payload []byte
	Msg     interface{}
}

const (
	MESSAGE_JSON     = "json"
	MESSAGE_PROTOBUF = "protobuf"
)

var reg = gregister.NewRegister()

func Register(builder gregister.Builder) {
	reg.Register(builder)
}

func NewMessageCodec(t string, c *gconfig.Configuration) (MessageCodec, error) {
	if node, err := reg.NewNode(t, c); err != nil {
		return nil, err
	} else {
		return node.(MessageCodec), nil
	}
}
