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
	PacketCodecType  string `toml:"packet"`
	MessageCodecType string `toml:"message"`
}

func NewProcessor(c *gconfig.Configuration) (*Processor, error) {
	proc := &Processor{}
	conf := &config{}
	var err error
	if err = c.UnmarshalKeyWithParent(Type(), conf); err != nil {
		return nil, err
	}
	proc.pktCodec, err = packet.NewPacketCodec(conf.PacketCodecType, c)
	if err != nil {
		return nil, err
	}
	proc.msgCodec, err = message.NewMessageCodec(conf.MessageCodecType, c)
	if err != nil {
		return nil, err
	}
	return proc, nil
}

func Type() string {
	return "processor"
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
	msgMeta := message.MessageMetaByID(msg.ID)
	msg.Msg = msgMeta.NewInstance()
	if err = p.msgCodec.Decode(msg.Msg, msg.Payload); err != nil {
		return nil, err
	}
	return msg, nil
}

func (p *Processor) Encode(rawMsg interface{}) ([]byte, error) {
	var err error
	msg := &message.Message{
		Msg: rawMsg,
	}
	msgMeta := message.MessageMetaByMsg(rawMsg)
	msg.ID = msgMeta.ID
	msg.Payload, err = p.msgCodec.Encode(rawMsg)
	if err != nil {
		return nil, err
	}
	data, err := p.pktCodec.Encode(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}
