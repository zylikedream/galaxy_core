package registery

import (
	"time"

	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type consulRegistery struct {
	conf   *consulConfig
	plugin *serverplugin.ConsulRegisterPlugin
}

type consulConfig struct {
	consulServers  []string      `toml:"consul_servers"`
	BasePath       string        `toml:"base_path"`
	UpdateInterval time.Duration `toml:"update_interval"`
}

func newconsulRegistery(ServerAddr string, c *gxyconfig.Configuration) (*consulRegistery, error) {
	conf := &consulConfig{}
	regist := &consulRegistery{
		conf: conf,
	}
	if err := c.UnmarshalKey(regist.Type(), conf); err != nil {
		return nil, err
	}
	regist.plugin = &serverplugin.ConsulRegisterPlugin{
		ServiceAddress: ServerAddr,
		ConsulServers:  conf.consulServers,
		BasePath:       conf.BasePath,
		UpdateInterval: conf.UpdateInterval,
	}
	return regist, nil
}

func (r *consulRegistery) Type() string {
	return REGISTERY_TYPE_CONSUL
}

func (r *consulRegistery) Start() error {
	return r.plugin.Start()
}

func (r *consulRegistery) GetPlugin() server.Plugin {
	return r.plugin
}

func (t *consulRegistery) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, gxyregister.ErrParamNotEnough
	}
	ServerAddr, ok := args[0].(string)
	if !ok {
		return nil, gxyregister.ErrParamErrType
	}
	return newconsulRegistery(ServerAddr, c)
}

func init() {
	gxyregister.Register((*consulRegistery)(nil))
}
