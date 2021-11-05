package processor

import (
	"fmt"
	"io"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/message"
	"github.com/zylikedream/galaxy/components/network/packet"
)

type Processor struct {
	pktCodec packet.PacketCodec
	msgCodec message.MessageCodec
	conf     *processorConfig
}

type processorConfig struct {
	PacketCodecType  string `toml:"packet"`
	MessageCodecType string `toml:"message"`
	PacketMaxSize    uint64 `toml:"packet_max_size"`
}

func NewProcessor(c *gconfig.Configuration) (*Processor, error) {
	proc := &Processor{}
	conf := &processorConfig{}
	var err error
	if err = c.UnmarshalKeyWithParent(Type(), conf); err != nil {
		return nil, err
	}
	proc.conf = conf
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

func (p *Processor) ReadMsg(r io.Reader) (*message.Message, error) {
	bodySize, err := p.pktCodec.ReadPacketLen(r)
	if err != nil {
		return nil, err
	}
	if bodySize > p.conf.PacketMaxSize {
		return nil, fmt.Errorf("packet size too big, %d(%d)", bodySize, p.conf.PacketMaxSize)
	}
	data, err := io.ReadAll(io.LimitReader(r, int64(bodySize)))
	if err != nil {
		return nil, err
	}
	if len(data) < int(bodySize) {
		return nil, fmt.Errorf("unexpect size %d(%d)", len(data), bodySize)
	}
	msg, err := p.pktCodec.DecodeBody(data)
	if err != nil {
		return nil, err
	}
	msgMeta := message.MessageMetaByID(msg.ID)
	msg.Msg = msgMeta.NewInstance()
	if err := p.msgCodec.Decode(msg.Msg, msg.Payload); err != nil {
		return nil, err
	}
	return msg, nil
}

func (p *Processor) WriteMsg(w io.Writer, rawMsg interface{}) (int, error) {
	var err error
	msg := &message.Message{
		Msg: rawMsg,
	}
	msgMeta := message.MessageMetaByMsg(rawMsg)
	msg.ID = msgMeta.ID
	msg.Payload, err = p.msgCodec.Encode(rawMsg)
	if err != nil {
		return 0, err
	}
	data, err := p.pktCodec.Encode(msg)
	if err != nil {
		return 0, err
	}
	return w.Write(data)
}
