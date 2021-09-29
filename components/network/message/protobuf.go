package message

import (
	"fmt"

	"github.com/zylikedream/galaxy/components/gconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoBuf struct {
	protoMap map[uint64]protoreflect.MessageType
	idMap    map[protoreflect.MessageType]uint64
}

func newProtobuf(_ *gconfig.Configuration) (*ProtoBuf, error) {
	return &ProtoBuf{
		protoMap: map[uint64]protoreflect.MessageType{},
		idMap:    map[protoreflect.MessageType]uint64{},
	}, nil
}

func (p *ProtoBuf) ReisterPacket(ID uint64, raw interface{}) error {
	msgType := raw.(protoreflect.MessageType)
	p.protoMap[ID] = msgType
	p.idMap[msgType] = ID
	return nil
}

func (p *ProtoBuf) Decode(ID uint64, data []byte) (interface{}, error) {
	msgType, Ok := p.protoMap[ID]
	if !Ok {
		return nil, fmt.Errorf("not found proto for ID:%d", ID)
	}
	msg := msgType.New().Interface()
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (p *ProtoBuf) Encode(raw interface{}) (uint64, []byte, error) {
	msg := raw.(proto.Message)
	id, Ok := p.idMap[msg.ProtoReflect().Type()]
	if !Ok {
		return 0, nil, fmt.Errorf("not found ID for proto:%s", proto.MessageName(msg))
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return 0, nil, err
	}
	return id, data, nil
}

func (p *ProtoBuf) Type() string {
	return MESSAGE_PROTOBUF
}

func (p *ProtoBuf) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newProtobuf(c)
}

func init() {
	Register(&ProtoBuf{})
}
