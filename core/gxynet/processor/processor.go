package processor

import (
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/message"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type Processor interface {
	Decode(data []byte) (uint64, *message.Message, error)
	Encode(msg *message.Message) ([]byte, error)
	GetMessageCodec() message.MessageCodec
}

type processorConfig struct {
	PacketCodecType  string `toml:"packet"`
	MessageCodecType string `toml:"message"`
}

func NewProcessor(t string, c *gxyconfig.Configuration) (Processor, error) {
	if node, err := gxyregister.NewNode("processor."+t, c); err != nil {
		return nil, err
	} else {
		return node.(Processor), nil
	}
}
