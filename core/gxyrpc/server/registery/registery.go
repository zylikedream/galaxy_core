package registery

import (
	"fmt"

	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
	"github.com/zylikedream/galaxy/core/gxyrpc/server/transport"
)

type Registery interface {
	GetPlugin() server.Plugin
	Start() error
}

const (
	REGISTERY_TYPE_ETCD   = "etcd"
	REGISTERY_TYPE_CONSUL = "consul"
	REGISTERY_TYPE_PEER   = "peer" // 客户端直连(比如点对点，点对多)
)

func RegisterAddr(t transport.Transport) string {
	return fmt.Sprintf("%s@%s", t.Network(), t.Addr())
}

func NewRegistery(t string, ServerAddr string, c *gconfig.Configuration) (Registery, error) {
	if node, err := gregister.NewNode(t, c.WithPrefix("registery"), ServerAddr); err != nil {
		return nil, err
	} else {
		return node.(Registery), nil
	}
}
