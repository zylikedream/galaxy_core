package message

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoBuf struct {
	protoMap map[int]protoreflect.MessageType
	idMap    map[protoreflect.MessageType]int
}

func (p *ProtoBuf) ReisterPacket(ID int, raw interface{}) error {
	msgType := raw.(protoreflect.MessageType)
	p.protoMap[ID] = msgType
	p.idMap[msgType] = ID
	return nil
}

func (p *ProtoBuf) Decode(ID int, data []byte) (interface{}, error) {
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

func (p *ProtoBuf) Encode(raw interface{}) ([]byte, int, error) {
	msg := raw.(proto.Message)
	id, Ok := p.idMap[msg.ProtoReflect().Type()]
	if !Ok {
		return nil, 0, fmt.Errorf("not found ID for proto:%s", proto.MessageName(msg))
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, 0, err
	}
	return data, id, nil
}
