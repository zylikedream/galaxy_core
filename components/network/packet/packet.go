package packet

import "github.com/zylikedream/galaxy/components/network/packet/codec"

type Packet struct {
	ID      int
	Type    int
	Payload []byte
	Msg     interface{}
	encoder codec.PacketEncoder
	decoder codec.PacketDecoder
}

func (p *Packet) Decode() error {
	var err error
	p.Msg, err = p.decoder.Decode(p.ID, p.Payload)
	return err
}
