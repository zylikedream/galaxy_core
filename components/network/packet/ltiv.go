package packet

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/message"
)

// length + type + id + payload
type ltiv struct {
	byteOrder binary.ByteOrder `toml:"byte_order"`
	conf      *ltivConfig
}

type ltivConfig struct {
	SizeLength int    `toml:"size_length"`
	TypeLength int    `toml:"type_length"`
	IDLength   int    `toml:"id_length"`
	ByteOrder  string `toml:"byte_order"`
}

func newLtiv(c *gconfig.Configuration) (*ltiv, error) {
	l := &ltiv{}
	conf := &ltivConfig{}
	if err := c.UnmarshalKeyWithParent(l.Type(), conf); err != nil {
		return nil, err
	}
	if conf.ByteOrder == "little" {
		l.byteOrder = binary.LittleEndian
	} else {
		l.byteOrder = binary.BigEndian
	}
	l.conf = conf
	return l, nil
}

func (l *ltiv) ReadPacketLen(r io.Reader) (uint64, error) {
	sizebuf, err := io.ReadAll(io.LimitReader(r, int64(l.conf.SizeLength)))
	if err != nil {
		return 0, err
	}
	// eof
	if len(sizebuf) == 0 {
		return 0, io.EOF
	}
	packetSize, err := Uint(sizebuf, l.byteOrder)
	if err != nil {
		return 0, err
	}
	return packetSize, nil
}

func (l *ltiv) MsgLenLength() int {
	return l.conf.SizeLength
}

func (l *ltiv) DecodeBody(payLoad []byte) (*message.Message, error) {
	msg := &message.Message{}
	// 消息类型+消息id+消息内容
	pointer := 0
	if tp, err := Uint(payLoad[pointer:pointer+l.conf.TypeLength], l.byteOrder); err != nil {
		return nil, err
	} else {
		msg.Type = tp
	}
	pointer += l.conf.TypeLength
	// 消息id
	if id, err := Uint(payLoad[pointer:pointer+l.conf.IDLength], l.byteOrder); err != nil {
		return nil, err
	} else {
		msg.ID = int(id)
	}
	pointer += l.conf.IDLength

	msg.Payload = payLoad[pointer:]

	return msg, nil
}

func (l *ltiv) ByteOrder() binary.ByteOrder {
	return l.byteOrder
}

func (l *ltiv) Encode(m *message.Message) ([]byte, error) {
	payload := &bytes.Buffer{}
	// 消息长度
	// 消息类型+消息id+消息内容
	if err := binary.Write(payload, l.byteOrder, convertUint(m.Type, l.conf.TypeLength)); err != nil {
		return nil, err
	}
	if err := binary.Write(payload, l.byteOrder, convertUint(uint64(m.ID), l.conf.IDLength)); err != nil {
		return nil, err
	}
	if _, err := payload.Write(m.Payload); err != nil {
		return nil, err
	}
	m.Payload = payload.Bytes()
	buf := &bytes.Buffer{}
	payloadLen := len(m.Payload)
	if err := binary.Write(buf, l.byteOrder, convertUint(uint64(payloadLen), l.conf.SizeLength)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(m.Payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
