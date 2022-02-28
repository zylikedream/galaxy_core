package processor

import (
	"strconv"

	"github.com/gookit/goutil/strutil"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/message"
	"github.com/zylikedream/galaxy/core/gxynet/packet"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type idProcessor struct {
	pktCodec packet.PacketCodec
	msgCodec message.MessageCodec
	conf     *processorConfig
}

func newIDProcessor(c *gxyconfig.Configuration) (Processor, error) {
	proc := &idProcessor{}
	conf := &processorConfig{}
	var err error
	if err = c.UnmarshalKey(proc.Type(), conf); err != nil {
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

func (p *idProcessor) Type() string {
	return "processor.id"
}

func (p *idProcessor) Decode(data []byte) (uint64, *message.Message, error) {
	pkgLen, msg, err := p.pktCodec.Decode(data)
	if err == packet.ErrPkgBodyNotEnough || err == packet.ErrPkgHeadNotEnough { // 数据不足够，不算错误
		return 0, nil, nil
	}
	if msg != nil {
		msgMeta := message.MessageMetaByID(strutil.MustInt(msg.Key))
		msg.Msg = msgMeta.NewInstance()
		if err := p.msgCodec.Decode(msg.Msg, msg.Payload); err != nil {
			return 0, nil, err
		}
	}
	return pkgLen, msg, err
}

func (p *idProcessor) Encode(msg *message.Message) ([]byte, error) {
	var err error
	msgMeta := message.MessageMetaByMsg(msg.Msg)
	msg.Key = strconv.Itoa(msgMeta.ID)
	msg.Payload, err = p.msgCodec.Encode(msg.Msg)
	if err != nil {
		return nil, err
	}
	data, err := p.pktCodec.Encode(msg)
	if err != nil {
		return nil, err
	}
	return data, err

}

func (p *idProcessor) GetMessageCodec() message.MessageCodec {
	return p.msgCodec
}

func (p *idProcessor) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newIDProcessor(c)
}

func init() {
	gxyregister.Register((*idProcessor)(nil))
}
