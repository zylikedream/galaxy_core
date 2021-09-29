package packet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/message"
)

// length + type + id + payload
type ltiv struct {
	sizeLength int              `mapstructure:"size_length"`
	typeLength int              `mapstructure:"type_length"`
	IDLength   int              `mapstructure:"id_length"`
	byteOrder  binary.ByteOrder `mapstructure:"byte_order"`
}

func newLtiv(c *gconfig.Configuration) (*ltiv, error) {
	l := &ltiv{}
	if err := c.UnmarshalKeyWithParent(l.Type(), l); err != nil {
		return nil, err
	}
	return l, nil
}

func (l *ltiv) MsgLenLength() int {
	return l.sizeLength
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
		msg.Type = tp
	}
	pointer += l.typeLength
	// 消息id
	if id, err := l.Uint(payLoad[pointer : pointer+l.IDLength]); err != nil {
		return nil, err
	} else {
		msg.ID = id
	}
	pointer += l.IDLength

	msg.Payload = payLoad[pointer:]

	return msg, nil
}

func (l *ltiv) ByteOrder() binary.ByteOrder {
	return l.byteOrder
}

func (l *ltiv) convertUint(v uint64, len int) interface{} {
	switch len {
	case 1:
		return uint8(v)
	case 2:
		return uint16(v)
	case 4:
		return uint32(v)
	case 8:
		return v
	}
	return v
}

func (l *ltiv) Encode(m *message.Message) ([]byte, error) {
	payload := bytes.Buffer{}
	// 消息类型+消息id+消息内容
	if err := binary.Write(&payload, l.byteOrder, l.convertUint(m.Type, l.typeLength)); err != nil {
		return nil, err
	}
	if err := binary.Write(&payload, l.byteOrder, l.convertUint(m.ID, l.IDLength)); err != nil {
		return nil, err
	}
	if _, err := payload.Write(m.Payload); err != nil {
		return nil, err
	}
	return payload.Bytes(), nil
}

func (l *ltiv) Type() string {
	return PACKET_LTIV
}

func (l *ltiv) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newLtiv(c)
}

func init() {
	Register(&ltiv{})
}
