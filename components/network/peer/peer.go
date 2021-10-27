package peer

import (
	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/gregister"
	"github.com/zylikedream/galaxy/components/network/processor"
)

const (
	PEER_TCP_SERVER    = "tcp_server"
	PEER_TCP_CONNECTOR = "tcp_connector"
)

type Peer interface {
	Start(h processor.MsgHandler) error
	Stop()
	Type() string
}

var reg = gregister.NewRegister()

func Register(builder gregister.Builder) {
	reg.Register(builder)
}

func NewPeer(t string, c *gconfig.Configuration) (Peer, error) {
	if node, err := reg.NewNode(t, c); err != nil {
		return nil, err
	} else {
		return node.(Peer), nil
	}
}
