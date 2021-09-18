package peer

import "github.com/zylikedream/galaxy/components/network/msg"

type Peer interface {
	Init() error
	Start() error
	Stop()
	Name() string
	MsgCodec() msg.Codec
}
