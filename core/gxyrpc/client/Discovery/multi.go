package discovery

import (
	"github.com/smallnest/rpcx/client"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
)

type multiDiscovery struct {
	conf *multiConfig
	d    client.ServiceDiscovery
}

type multiConfig struct {
	servers []string `toml:"servers"` // ex:tcp@localhost:2786
}

func newmultiDiscovery(c *gconfig.Configuration) (*multiDiscovery, error) {
	conf := &multiConfig{}
	multi := &multiDiscovery{
		conf: conf,
	}
	if err := c.UnmarshalKeyWithPrefix(multi.Type(), conf); err != nil {
		return nil, err
	}
	pairs := []*client.KVPair{}
	for _, svr := range conf.servers {
		pairs = append(pairs, &client.KVPair{Key: svr})
	}
	d, err := client.NewMultipleServersDiscovery(pairs)
	if err != nil {
		return nil, err
	}
	multi.d = d
	return multi, nil
}

func (m *multiDiscovery) Type() string {
	return DISCOVERY_TYPE_MUTISERVER
}

func (m *multiDiscovery) GetDiscovery() client.ServiceDiscovery {
	return m.d
}

func (m *multiDiscovery) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newmultiDiscovery(c)
}

func init() {
	gregister.Register((*multiDiscovery)(nil))
}
