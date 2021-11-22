package discovery

import (
	"github.com/smallnest/rpcx/client"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
)

type consulDiscovery struct {
	conf *consulConfig
	d    client.ServiceDiscovery
}

type consulConfig struct {
	ConsulServers []string `toml:"etcd_servers"`
	BasePath      string   `toml:"base_path"`
}

func newconsulDiscovery(c *gconfig.Configuration) (*consulDiscovery, error) {
	conf := &consulConfig{}
	consul := &consulDiscovery{
		conf: conf,
	}
	if err := c.UnmarshalKey(consul.Type(), conf); err != nil {
		return nil, err
	}
	d, err := client.NewConsulDiscoveryTemplate(conf.BasePath, conf.ConsulServers, nil)
	if err != nil {
		return nil, err
	}
	consul.d = d
	return consul, nil
}

func (c *consulDiscovery) Type() string {
	return DISCOVERY_TYPE_CONSUL
}

func (c *consulDiscovery) GetDiscovery() client.ServiceDiscovery {
	return c.d
}

func (cd *consulDiscovery) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newconsulDiscovery(c)
}

func init() {
	gregister.Register((*consulDiscovery)(nil))
}
