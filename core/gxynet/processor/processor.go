package processor

import (
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/message"
	"github.com/zylikedream/galaxy/core/gxynet/packet"
)

type Processor struct {
	pktCodec packet.PacketCodec
	msgCodec message.MessageCodec
	conf     *processorConfig
}

type processorConfig struct {
	PacketCodecType  string `toml:"packet"`
	MessageCodecType string `toml:"message"`
}

func NewProcessor(c *gxyconfig.Configuration) (*Processor, error) {
	proc := &Processor{}
	conf := &processorConfig{}
	var err error
	if err = c.UnmarshalKey(Type(), conf); err != nil {
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
	return "peer.processor"
}

func (p *Processor) Decode(data []byte) (uint64, *message.Message, error) {
	pkgLen, msg, err := p.pktCodec.Decode(data)
	if err == packet.ErrPkgBodyNotEnough || err == packet.ErrPkgHeadNotEnough { // 数据不足够，不算错误
		return 0, nil, nil
	}
	if msg != nil {
		msgMeta := message.MessageMetaByID(msg.ID)
		msg.Msg = msgMeta.NewInstance()
		if err := p.msgCodec.Decode(msg.Msg, msg.Payload); err != nil {
			return 0, nil, err
		}
	}
	return pkgLen, msg, err
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
	return data, err

}

func (p *Processor) GetMessageCodec() message.MessageCodec {
	return p.msgCodec
}
