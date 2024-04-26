package peer

import (
	"context"
	"fmt"
	"net"

	"github.com/panjf2000/gnet/v2"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/gxynet/endpoint"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type TcpServer struct {
	endpoint.CoreBundle
	listener net.Listener
	conf     *tcpServerConfig
	gnet.BuiltinEventEngine
	engine gnet.Engine
}

type tcpServerConfig struct {
	Addr     string `toml:"addr"`
	ProcType string `toml:"processor"`
}

func newTcpServer(c *gxyconfig.Configuration) (*TcpServer, error) {
	server := &TcpServer{}
	conf := &tcpServerConfig{}
	if err := c.UnmarshalKey(server.Type(), conf); err != nil {
		return nil, err
	}
	conf.Addr = fmt.Sprintf("tcp://%s", conf.Addr)
	server.conf = conf
	if err := server.BindProc(c, conf.ProcType); err != nil {
		return nil, err
	}
	return server, nil
}

func (t *TcpServer) Init() error {
	return nil
}

func (t *TcpServer) Start(ctx context.Context, el endpoint.EventHandler) error {
	t.Handler = el
	return gnet.Run(t, t.conf.Addr, gnet.WithMulticore(true), gnet.WithReuseAddr(true))
}

func (t *TcpServer) Type() string {
	return PEER_TCP_SERVER
}

func (t *TcpServer) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTcpServer(c)
}

func (t *TcpServer) Stop(ctx context.Context) {
	t.engine.Stop(ctx)
}

func init() {
	gxyregister.Register((*TcpServer)(nil))
}

func (t *TcpServer) OnBoot(eng gnet.Engine) gnet.Action {
	t.engine = eng
	return gnet.None
}

func (t *TcpServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	endPoint := endpoint.NewTcpEndPoint(c, t.Processor)
	c.SetContext(endPoint)
	t.Handler.OnOpen(endPoint)
	return nil, gnet.None
}

func (t *TcpServer) OnTraffic(c gnet.Conn) gnet.Action {
	endPoint := c.Context().(*endpoint.TcpEndpoint)
	data, err := c.Next(-1)
	gxylog.Debugf("receive data %s", string(data))
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

func (t *TcpServer) OnClose(c gnet.Conn, err error) gnet.Action {
	gxylog.Errorf("conn close %s, error %s", c.RemoteAddr().String(), err.Error())
	return gnet.None
}
