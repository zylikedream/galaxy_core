package peer

import (
	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/register"
)

const (
	PEER_TCP_ACCEPTOR = "tcp_acceptor"
)

type Peer interface {
	Init() error
	Start() error
	Stop()
	Type() string
}

var reg = register.NewRegister()

func Register(t string, f func(c *gconfig.Configuration) (interface{}, error)) {
	reg.Register(t, f)
}

func NewPeer(t string, c *gconfig.Configuration) (Peer, error) {
	if node, err := reg.NewNode(t, c); err != nil {
		return nil, err
	} else {
		return node.(Peer), nil
	}
}
