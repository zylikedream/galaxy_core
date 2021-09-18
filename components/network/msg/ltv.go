package msg

import (
	"encoding/binary"
	"fmt"

	"github.com/zylikedream/galaxy/components/network/packet"
)

// length + id + payload
type Ltv struct {
	lenLength  int
	typeLength int
	IDLength   int
	byteOrder  binary.ByteOrder
}

func NewLtv(ll, tl, idl int, bo binary.ByteOrder) *Ltv {
	return &Ltv{
		lenLength:  ll,
		typeLength: tl,
		IDLength:   idl,
		byteOrder:  bo,
	}
}

func (l *Ltv) MsgLenLength() int {
	return l.lenLength
}

func (l *Ltv) MsgTypeLength() int {
	return l.typeLength
}

func (l *Ltv) MsgIDLength() int {
	return l.IDLength
}

func (l *Ltv) Uint(data []byte) (uint64, error) {
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

func (l *Ltv) Decode(payLoad []byte) (*packet.Packet, error) {
	packet := &packet.Packet{}
	// 字节类型+消息id+消息内容
	pointer := 0
	if tp, err := l.Uint(payLoad[pointer : pointer+l.typeLength]); err != nil {
		return nil, err
	} else {
		packet.Type = int(tp)
	}
	pointer += l.typeLength
	// 消息id
	if id, err := l.Uint(payLoad[pointer : pointer+l.IDLength]); err != nil {
		return nil, err
	} else {
		packet.ID = int(id)
	}
	pointer += l.IDLength

	packet.Payload = payLoad[pointer:]
	return packet, nil
}
