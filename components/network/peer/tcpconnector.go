package peer

import (
	"net"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/processor"
	"github.com/zylikedream/galaxy/components/network/session"
)

type TcpConnector struct {
	processor.ProcessorBundle
	conf *tcpConnectorConfig
}

type tcpConnectorConfig struct {
	Addr string `toml:"addr"`
}

func newTcpConnetor(c *gconfig.Configuration) (*TcpConnector, error) {
	server := &TcpConnector{}
	conf := &tcpConnectorConfig{}
	if err := c.UnmarshalKeyWithParent(server.Type(), conf); err != nil {
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

func (t *TcpConnector) Start() error {
	conn, err := net.Dial("tcp", t.conf.Addr)
	if err != nil {
		return err
	}

	sess := session.NewTcpSession(conn, t.ProcessorBundle)
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
	Register(&TcpConnector{})
}
