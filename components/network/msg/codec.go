package msg

import (
	"encoding/binary"

	"github.com/zylikedream/galaxy/components/network/packet"
)

type Codec interface {
	MsgLenLength() int  // 长度字段的字节长度
	MsgTypeLength() int // 类型字段的字节长度
	MsgIDLength() int   // msgID的字节长度
	Decode(payLoad []byte) (*packet.Packet, error)
	ByteOrder() binary.ByteOrder
	Uint(data []byte) uint64
}
