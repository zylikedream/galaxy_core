package transport

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
)

type tcpConfig struct {
	Addr string `toml:"addr"`
}

type tcpTransport struct {
	conf *tcpConfig
}

func newTcpTransport(c *gconfig.Configuration) (*tcpTransport, error) {
	tran := &tcpTransport{}
	conf := &tcpConfig{}
	if err := c.UnmarshalKeyWithParent(tran.Type(), conf); err != nil {
		return nil, err
	}
	tran.conf = conf
	return tran, nil
}

func (t *tcpTransport) Addr() string {
	return t.conf.Addr
}

func (t *tcpTransport) Network() string {
	return "tcp"
}

func (t *tcpTransport) Options() server.OptionFn {
	return nil
}

func (t *tcpTransport) Type() string {
	return TRANSPORT_TYPE_TCP
}

func (t *tcpTransport) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTcpTransport(c)
}

func init() {
	gregister.Register((*tcpTransport)(nil))
}
