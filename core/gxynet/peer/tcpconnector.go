package peer

import (
	"context"
	"net"

	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxynet/endpoint"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type TcpConnector struct {
	endpoint.ConnBundle
	conf *tcpConnectorConfig
}

type tcpConnectorConfig struct {
	Addr string `toml:"addr"`
}

func newTcpConnetor(c *gxyconfig.Configuration) (*TcpConnector, error) {
	server := &TcpConnector{}
	conf := &tcpConnectorConfig{}
	if err := c.UnmarshalKey(server.Type(), conf); err != nil {
		return nil, err
	}
	server.conf = conf
	if err := server.BindProc(c); err != nil {
		return nil, err
	}
	return server, nil
}

func (t *TcpConnector) Init() error {
	return nil
}

func (t *TcpConnector) Start(ctx context.Context, h endpoint.EventHandler) error {
	con, err := net.Dial("tcp", t.conf.Addr)
	if err != nil {
		return err
	}
	t.BindHandler(h)
	sess := endpoint.NewTcpConn(con, t.ConnBundle)
	return sess.Start(ctx)
}

func (t *TcpConnector) Type() string {
	return PEER_TCP_CONNECTOR
}

func (t *TcpConnector) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTcpConnetor(c)
}

func (t *TcpConnector) Stop(ctx context.Context) {
}

func init() {
	gxyregister.Register((*TcpConnector)(nil))
}
