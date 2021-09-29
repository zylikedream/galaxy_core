package packet

import (
	"encoding/binary"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/register"
)

type PacketCodec interface {
	MsgLenLength() int // 长度字段的字节长度
	Decode(payLoad []byte) (*message.Message, error)
	Encode(msg *message.Message) ([]byte, error)
	ByteOrder() binary.ByteOrder
	Uint(data []byte) (uint64, error)
	Type() string
}

const (
	PACKET_LTIV = "ltiv"
)

var reg = register.NewRegister()

func Register(t string, f register.FuncType) {
	reg.Register(t, f)
}

func NewPacketCodec(t string, c *gconfig.Configuration) (PacketCodec, error) {
	if node, err := reg.NewNode(t, c); err != nil {
		return nil, err
	} else {
		return node.(PacketCodec), nil
	}
}
