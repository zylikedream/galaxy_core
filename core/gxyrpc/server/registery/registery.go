package registery

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
	"github.com/zylikedream/galaxy/core/gxyrpc/server/transport"
)

type Registery interface {
	GetPlugin() server.Plugin
	Start() error
}

const (
	REGISTERY_TYPE_ETCD   = "registery.etcd"
	REGISTERY_TYPE_CONSUL = "registery.consul"
	REGISTERY_TYPE_PEER   = "registery.peer" // 客户端直连(比如点对点，点对多)
)

func RegisterAddr(t transport.Transport) string {
	return fmt.Sprintf("%s@%s", t.Network(), t.Addr())
}

func NewRegistery(t string, ServerAddr string, c *gxyconfig.Configuration) (Registery, error) {
	if node, err := gxyregister.NewNode("registery."+t, c, ServerAddr); err != nil {
		return nil, errors.Wrap(err, "new registery failed")
	} else {
		return node.(Registery), nil
	}
}
