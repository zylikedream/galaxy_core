package discovery

import (
	"github.com/pkg/errors"
	"github.com/smallnest/rpcx/client"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type Discovery interface {
	GetDiscovery() client.ServiceDiscovery
}

const (
	DISCOVERY_TYPE_P2P        = "discovery.p2p"
	DISCOVERY_TYPE_MUTISERVER = "discovery.mutiserver"
	DISCOVERY_TYPE_ETCD       = "discovery.etcd"
	DISCOVERY_TYPE_CONSUL     = "discovery.consul"
)

func NewDisvoery(t string, c *gxyconfig.Configuration) (Discovery, error) {
	if node, err := gxyregister.NewNode("discovery."+t, c); err != nil {
		return nil, errors.Wrap(err, "new discovery failed")
	} else {
		return node.(Discovery), nil
	}
}
