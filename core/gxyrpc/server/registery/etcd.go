package registery

import (
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type etcdRegistery struct {
	conf   *etcdConfig
	plugin *serverplugin.EtcdV3RegisterPlugin
}

type etcdConfig struct {
	EtcdServers    []string      `toml:"etcd_servers"`
	BasePath       string        `toml:"base_path"`
	UpdateInterval time.Duration `toml:"update_interval"`
}

func newEtcdRegistery(ServerAddr string, c *gconfig.Configuration) (*etcdRegistery, error) {
	conf := &etcdConfig{}
	regist := &etcdRegistery{
		conf: conf,
	}
	if err := c.UnmarshalKey(regist.Type(), conf); err != nil {
		return nil, err
	}
	regist.plugin = &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: ServerAddr,
		EtcdServers:    conf.EtcdServers,
		BasePath:       conf.BasePath,
		UpdateInterval: conf.UpdateInterval,
	}
	return regist, nil
}

func (r *etcdRegistery) Type() string {
	return REGISTERY_TYPE_ETCD
}

func (r *etcdRegistery) Start() error {
	return r.plugin.Start()
}

func (r *etcdRegistery) GetPlugin() server.Plugin {
	return r.plugin
}

func (t *etcdRegistery) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, gxyregister.ErrParamNotEnough
	}
	ServerAddr, ok := args[0].(string)
	if !ok {
		return nil, gxyregister.ErrParamErrType
	}
	return newEtcdRegistery(ServerAddr, c)
}

func init() {
	gxyregister.Register((*etcdRegistery)(nil))
}
