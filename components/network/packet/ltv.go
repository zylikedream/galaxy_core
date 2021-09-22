package packet

import (
	"encoding/binary"
	"fmt"

	"github.com/zylikedream/galaxy/components/network/message"
)

// length + type + id + payload
type ltiv struct {
	lenLength  int
	typeLength int
	IDLength   int
	byteOrder  binary.ByteOrder
}

func NewLtiv(ll, tl, idl int, bo binary.ByteOrder) *ltiv {
	return &ltiv{
		lenLength:  ll,
		typeLength: tl,
		IDLength:   idl,
		byteOrder:  bo,
	}
}

func (l *ltiv) MsgLenLength() int {
	return l.lenLength
}

func (l *ltiv) Uint(data []byte) (uint64, error) {
	switch len(data) {
	case 1:
		return uint64(data[0]), nil
	case 2:
		return uint64(l.byteOrder.Uint16(data)), nil
	case 4:
		return uint64(l.byteOrder.Uint32(data)), nil
	case 8:
		return uint64(l.byteOrder.Uint64(data)), nil
	}
	return 0, fmt.Errorf("unsupport byte len:%d", len(data))
}

func (l *ltiv) Decode(payLoad []byte) (*message.Message, error) {
	msg := &message.Message{}
	// 消息类型+消息id+消息内容
	pointer := 0
	if tp, err := l.Uint(payLoad[pointer : pointer+l.typeLength]); err != nil {
		return nil, err
	} else {
		msg.Type = int(tp)
	}
	pointer += l.typeLength
	// 消息id
	if id, err := l.Uint(payLoad[pointer : pointer+l.IDLength]); err != nil {
		return nil, err
	} else {
		msg.ID = int(id)
	}
	pointer += l.IDLength

	msg.Payload = payLoad[pointer:]

	return msg, nil
}
