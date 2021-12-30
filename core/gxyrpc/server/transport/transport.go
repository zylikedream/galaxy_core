package transport

import (
	"github.com/pkg/errors"
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type Transport interface {
	Addr() string
	Network() string
	Option() server.OptionFn
}

var emptyOptin = server.OptionFn(func(*server.Server) {})

const (
	TRANSPORT_TYPE_TCP  = "transport.tcp"
	TRANSPORT_TYPE_TLS  = "transport.tls"
	TRANSPORT_TYPE_UNIX = "transport.unix"
)

func NewTransport(t string, c *gconfig.Configuration) (Transport, error) {
	if node, err := gxyregister.NewNode("transport."+t, c); err != nil {
		return nil, errors.Wrap(err, "new transport failed")
	} else {
		return node.(Transport), nil
	}
}
