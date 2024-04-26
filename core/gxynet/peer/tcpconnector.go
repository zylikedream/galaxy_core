package peer

import (
	"context"

	"github.com/panjf2000/gnet/v2"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/gxynet/endpoint"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type TcpConnector struct {
	endpoint.CoreBundle
	conf *tcpConnectorConfig
	gnet.BuiltinEventEngine
	engine gnet.Engine
}

type tcpConnectorConfig struct {
	Addr     string `toml:"addr"`
	ProcType string `toml:"processor"`
}

func newTcpConnetor(c *gxyconfig.Configuration) (*TcpConnector, error) {
	server := &TcpConnector{}
	conf := &tcpConnectorConfig{}
	if err := c.UnmarshalKey(server.Type(), conf); err != nil {
		return nil, err
	}
	server.conf = conf
	if err := server.BindProc(c, conf.ProcType); err != nil {
		return nil, err
	}
	return server, nil
}

func (t *TcpConnector) Init() error {
	return nil
}

func (t *TcpConnector) Start(ctx context.Context, h endpoint.EventHandler) error {
	t.BindHandler(h)
	cli, err := gnet.NewClient(t)
	if err != nil {
		return err
	}
	if err := cli.Start(); err != nil {
		return err
	}
	_, err = cli.Dial("tcp", t.conf.Addr)
	if err != nil {
		return err
	}
	err = cli.Start()
	if err != nil {
		return err
	}
	return nil
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

func (t *TcpConnector) OnBoot(eng gnet.Engine) gnet.Action {
	t.engine = eng
	return gnet.None
}

func (t *TcpConnector) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	gxylog.Info("connect conn info")
	endPoint := endpoint.NewTcpEndPoint(c, t.Processor)
	c.SetContext(endPoint)
	t.Handler.OnOpen(endPoint)
	return nil, gnet.None
}

func (t *TcpConnector) OnTraffic(c gnet.Conn) gnet.Action {
	endPoint := c.Context().(*endpoint.TcpEndpoint)
	data, err := c.Next(-1)
	if err != nil {
		gxylog.Errorf("get traffic data failed %s", err.Error())
		return gnet.Close
	}
	for {
		msg, err := endPoint.DecodeMsg(data)
		if err != nil {
			break
		}
		t.Handler.OnMessage(endPoint, msg)
	}
	return gnet.None
}

func (t *TcpConnector) OnClose(c gnet.Conn, err error) gnet.Action {
	gxylog.Errorf("conn close %s, error %s", c.RemoteAddr().String(), err.Error())
	return gnet.None
}
