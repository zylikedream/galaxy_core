package discovery

import (
	"github.com/smallnest/rpcx/client"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
)

type Discovery interface {
	GetDiscovery() client.ServiceDiscovery
}

const (
	DISCOVERY_TYPE_P2P        = "p2p"
	DISCOVERY_TYPE_MUTISERVER = "mutiserver"
	DISCOVERY_TYPE_ETCD       = "etcd"
	DISCOVERY_TYPE_CONSUL     = "consul"
)

func NewDisvoery(t string, c *gconfig.Configuration) (Discovery, error) {
	if node, err := gregister.NewNode(t, c.WithPrefix("discovery")); err != nil {
		return nil, err
	} else {
		return node.(Discovery), nil
	}
}
