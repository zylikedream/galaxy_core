package packet

import (
	"bytes"
	"encoding/binary"
	"strconv"

	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/message"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

// length + type + id + payload
type ltiv struct {
	byteOrder binary.ByteOrder `toml:"byte_order"`
	conf      *ltivConfig
}

type ltivConfig struct {
	SizeLength int    `toml:"size_length"`
	MaxSize    int    `toml:"max_size"`
	TypeLength int    `toml:"type_length"`
	IDLength   int    `toml:"id_length"`
	ByteOrder  string `toml:"byte_order"`
}

func newLtiv(c *gxyconfig.Configuration) (*ltiv, error) {
	l := &ltiv{}
	conf := &ltivConfig{}
	if err := c.UnmarshalKey(l.Type(), conf); err != nil {
		return nil, errors.WithStack(err)
	}
	if conf.ByteOrder == "little" {
		l.byteOrder = binary.LittleEndian
	} else {
		l.byteOrder = binary.BigEndian
	}
	l.conf = conf
	return l, nil
}

func (l *ltiv) decodeBody(payLoad []byte) (*message.Message, error) {
	msg := &message.Message{}
	// 消息类型+消息id+消息内容
	if tp, err := uintDecode(payLoad[:l.conf.TypeLength], l.byteOrder); err != nil {
		return nil, errors.WithStack(err)
	} else {
		msg.Type = tp
	}
	payLoad = payLoad[l.conf.TypeLength:]
	// 消息id
	if id, err := uintDecode(payLoad[:l.conf.IDLength], l.byteOrder); err != nil {
		return nil, errors.WithStack(err)
	} else {
		msg.Key = strconv.Itoa(int(id))
	}
	payLoad = payLoad[l.conf.IDLength:]

	msg.Payload = payLoad

	return msg, nil
}

func (l *ltiv) ByteOrder() binary.ByteOrder {
	return l.byteOrder
}

func (l *ltiv) Decode(data []byte) (uint64, *message.Message, error) {
	if len(data) < l.conf.SizeLength {
		return 0, nil, ErrPkgHeadNotEnough
	}
	dataSize, err := uintDecode(data[:l.conf.SizeLength], l.byteOrder)
	if err != nil {
		return 0, nil, err
	}
	if dataSize > uint64(l.conf.MaxSize) {
		return 0, nil, errors.WithStack(errors.Errorf("packet too big, %d(%d)", dataSize, l.conf.MaxSize))
	}
	data = data[l.conf.SizeLength:]
	if len(data) < int(dataSize) {
		return 0, nil, errors.WithStack(ErrPkgBodyNotEnough)
	}
	msg, err := l.decodeBody(data[:dataSize])
	return dataSize + uint64(l.conf.SizeLength), msg, err
}

func (l *ltiv) Encode(m *message.Message) ([]byte, error) {
	payload := &bytes.Buffer{}
	// 消息长度
	// 消息类型+消息id+消息内容
	if err := binary.Write(payload, l.byteOrder, convertUint(m.Type, l.conf.TypeLength)); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := binary.Write(payload, l.byteOrder, convertUint(uint64(strutil.MustInt(m.Key)), l.conf.IDLength)); err != nil {
		return nil, errors.WithStack(err)
	}
	if _, err := payload.Write(m.Payload); err != nil {
		return nil, errors.WithStack(err)
	}
	m.Payload = payload.Bytes()
	buf := &bytes.Buffer{}
	payloadLen := len(m.Payload)
	if err := binary.Write(buf, l.byteOrder, convertUint(uint64(payloadLen), l.conf.SizeLength)); err != nil {
		return nil, errors.WithStack(err)
	}
	if _, err := buf.Write(m.Payload); err != nil {
		return nil, errors.WithStack(err)
	}
	return buf.Bytes(), nil
}

func (l *ltiv) Type() string {
	return PACKET_LTIV
}

func (l *ltiv) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newLtiv(c)
}

func init() {
	gxyregister.Register((*ltiv)(nil))
}
