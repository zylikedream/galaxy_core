package transport

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
)

type unixConfig struct {
	Socket string `toml:"socket"`
}

type unixTransport struct {
	conf *unixConfig
}

func newUnixTransport(c *gconfig.Configuration) (*unixTransport, error) {
	tran := &unixTransport{}
	conf := &unixConfig{}
	if err := c.UnmarshalKeyWithParent(tran.Type(), conf); err != nil {
		return nil, err
	}
	tran.conf = conf
	return tran, nil
}

func (t *unixTransport) Addr() string {
	return t.conf.Socket
}

func (t *unixTransport) Network() string {
	return "unix"
}

func (t *unixTransport) Options() server.OptionFn {
	return nil
}

func (t *unixTransport) Type() string {
	return TRANSPORT_TYPE_TCP
}

func (t *unixTransport) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newUnixTransport(c)
}

func init() {
	gregister.Register((*unixTransport)(nil))
}
