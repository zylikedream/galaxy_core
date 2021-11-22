package transport

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
)

type unixConfig struct {
	Addr string `toml:"addr"`
}

type unixTransport struct {
	conf *unixConfig
}

func newUnixTransport(c *gconfig.Configuration) (*unixTransport, error) {
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

func (t *unixTransport) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newUnixTransport(c)
}

func init() {
	gregister.Register((*unixTransport)(nil))
}
