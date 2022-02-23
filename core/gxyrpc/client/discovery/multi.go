package discovery

import (
	"github.com/smallnest/rpcx/client"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type multiDiscovery struct {
	conf *multiConfig
	d    client.ServiceDiscovery
}

type multiConfig struct {
	Peers []string `toml:"peers"` // ex:tcp@localhost:2786
}

func newmultiDiscovery(c *gxyconfig.Configuration) (*multiDiscovery, error) {
	conf := &multiConfig{}
	multi := &multiDiscovery{
		conf: conf,
	}
	if err := c.UnmarshalKey(multi.Type(), conf); err != nil {
		return nil, err
	}
	pairs := []*client.KVPair{}
	for _, svr := range conf.Peers {
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

func (m *multiDiscovery) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newmultiDiscovery(c)
}

func init() {
	gxyregister.Register((*multiDiscovery)(nil))
}
