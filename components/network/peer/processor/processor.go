package processor

import (
	"io"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/packet"
)

type Processor struct {
	pktCodec packet.PacketCodec
	msgCodec message.MessageCodec
}

type config struct {
	PacketCodecType  string `mapstructure:"packet"`
	MessageCodecType string `mapstructure:"message"`
}

func NewProcessor(c *gconfig.Configuration) (*Processor, error) {
	conf := &config{}
	if err := c.UnmarshalKey("network.processor", conf); err != nil {
		return nil, err
	}
	packetCodec, err := packet.NewPacketCodec(conf.PacketCodecType, c)
	if err != nil {
		return nil, err
	}
	messageCodec, err := message.NewMessageCodec(conf.MessageCodecType, c)
	if err != nil {
		return nil, err
	}
	return &Processor{
		pktCodec: packetCodec,
		msgCodec: messageCodec,
	}, nil
}

func (p *Processor) ReadAndDecode(r io.Reader) (*message.Message, error) {
	sizebuf, err := io.ReadAll(io.LimitReader(r, int64(p.pktCodec.MsgLenLength())))
	if err != nil {
		return nil, err
	}
	// eof
	if len(sizebuf) == 0 {
		return nil, io.EOF
	}
	size, err := p.pktCodec.Uint(sizebuf)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(io.LimitReader(r, int64(size)))
	if err != nil {
		return nil, err
	}
	if len(data) < int(size) {
		return nil, err
	}
	msg, err := p.pktCodec.Decode(data)
	if err != nil {
		return nil, err
	}
	msg.Msg, err = p.msgCodec.Decode(msg.ID, msg.Payload)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (p *Processor) Encode(rawMsg interface{}) ([]byte, error) {
	var err error
	msg := &message.Message{
		Msg: rawMsg,
	}
	msg.ID, msg.Payload, err = p.msgCodec.Encode(rawMsg)
	if err != nil {
		return nil, err
	}
	data, err := p.pktCodec.Encode(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}
