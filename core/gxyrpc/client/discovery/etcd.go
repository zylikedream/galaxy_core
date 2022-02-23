package discovery

import (
	etcd_client "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type etcdDiscovery struct {
	conf *etcdConfig
	d    client.ServiceDiscovery
}

type etcdConfig struct {
	EtcdServers      []string `toml:"etcd_servers"`
	BasePath         string   `toml:"base_path"`
	AllowKeyNotFound bool     `toml:"allow_key_not_found"`
}

func newEtcdDiscovery(c *gxyconfig.Configuration) (*etcdDiscovery, error) {
	conf := &etcdConfig{}
	etcd := &etcdDiscovery{
		conf: conf,
	}
	if err := c.UnmarshalKey(etcd.Type(), conf); err != nil {
		return nil, err
	}
	d, err := etcd_client.NewEtcdV3DiscoveryTemplate(conf.BasePath, conf.EtcdServers, conf.AllowKeyNotFound, nil)
	if err != nil {
		return nil, err
	}
	etcd.d = d
	return etcd, nil
}

func (e *etcdDiscovery) Type() string {
	return DISCOVERY_TYPE_ETCD
}

func (e *etcdDiscovery) GetDiscovery() client.ServiceDiscovery {
	return e.d
}

func (e *etcdDiscovery) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newEtcdDiscovery(c)
}

func init() {
	gxyregister.Register((*etcdDiscovery)(nil))
}
