package peer

import (
	"net"

	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
	"github.com/zylikedream/galaxy/core/network/session"
)

type TcpConnector struct {
	session.SessionBundle
	conf *tcpConnectorConfig
}

type tcpConnectorConfig struct {
	Addr string `toml:"addr"`
}

func newTcpConnetor(c *gconfig.Configuration) (*TcpConnector, error) {
	server := &TcpConnector{}
	conf := &tcpConnectorConfig{}
	if err := c.UnmarshalKeyWithPrefix(server.Type(), conf); err != nil {
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

func (t *TcpConnector) Start(h session.EventHandler) error {
	conn, err := net.Dial("tcp", t.conf.Addr)
	if err != nil {
		return err
	}
	t.BindHandler(h)
	sess := session.NewTcpSession(conn, t.SessionBundle)
	sess.Start()
	return nil
}

func (t *TcpConnector) Type() string {
	return PEER_TCP_CONNECTOR
}

func (t *TcpConnector) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTcpConnetor(c)
}

func (t *TcpConnector) Stop() {
}

func init() {
	gregister.Register((*TcpConnector)(nil))
}
