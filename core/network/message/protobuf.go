package message

import (
	"github.com/zylikedream/galaxy/core/gconfig"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoBuf struct {
}

func newProtobuf(_ *gconfig.Configuration) (*ProtoBuf, error) {
	return &ProtoBuf{}, nil
}

func (p *ProtoBuf) Decode(msg interface{}, data []byte) error {
	return proto.Unmarshal(data, msg.(protoreflect.ProtoMessage))
}

func (p *ProtoBuf) Encode(raw interface{}) ([]byte, error) {
	data, err := proto.Marshal(raw.(proto.Message))
	if err != nil {
		return nil, err
	}
	return data, nil
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
