package transport

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type unixConfig struct {
	Addr string `toml:"addr"`
}

type unixTransport struct {
	conf *unixConfig
}

func newUnixTransport(c *gxyconfig.Configuration) (*unixTransport, error) {
	tran := &unixTransport{}
	conf := &unixConfig{}
	if err := c.UnmarshalKey(tran.Type(), conf); err != nil {
		return nil, err
	}
	tran.conf = conf
	return tran, nil
}

func (t *unixTransport) Addr() string {
	return t.conf.Addr
}

func (t *unixTransport) Network() string {
	return "unix"
}

func (t *unixTransport) Option() server.OptionFn {
	return emptyOptin
}

func (t *unixTransport) Type() string {
	return TRANSPORT_TYPE_UNIX
}

func (t *unixTransport) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newUnixTransport(c)
}

func init() {
	gxyregister.Register((*unixTransport)(nil))
}
