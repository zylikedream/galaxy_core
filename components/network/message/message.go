package message

import (
	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/register"
	"github.com/zylikedream/galaxy/components/network/session"
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
	Sess    session.Session
}

const (
	MESSAGE_JSON     = "json"
	MESSAGE_PROTOBUF = "protobuf"
)

var reg = register.NewRegister()

func Register(t string, f register.FuncType) {
	reg.Register(t, f)
}

func NewMessageCodec(t string, c *gconfig.Configuration) (MessageCodec, error) {
	if node, err := reg.NewNode(t, c); err != nil {
		return nil, err
	} else {
		return node.(MessageCodec), nil
	}
}
