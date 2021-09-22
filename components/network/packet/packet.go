package packet

import (
	"encoding/binary"

	"github.com/zylikedream/galaxy/components/network/message"
)

type PacketCodec interface {
	MsgLenLength() int // 长度字段的字节长度
	Decode(payLoad []byte) (*message.Message, error)
	Encode(msg *message.Message) ([]byte, error)
	ByteOrder() binary.ByteOrder
	Uint(data []byte) uint64
}
