package transport

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
)

type Transport interface {
	Addr() string
	Network() string
	Option() server.OptionFn
}

const (
	TRANSPORT_TYPE_TCP  = "tcp"
	TRANSPORT_TYPE_TLS  = "tls"
	TRANSPORT_TYPE_UNIX = "unix"
)

func NewTransport(t string, c *gconfig.Configuration) (Transport, error) {
	if node, err := gregister.NewNode(t, c); err != nil {
		return nil, err
	} else {
		return node.(Transport), nil
	}
}
