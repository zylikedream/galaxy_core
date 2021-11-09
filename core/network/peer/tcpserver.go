package peer

import (
	"net"
	"time"

	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/network/logger"
	"github.com/zylikedream/galaxy/core/network/session"
	"go.uber.org/zap"
)

type TcpServer struct {
	session.SessionBundle
	listener net.Listener
	conf     *tcpServerConfig
}

type tcpServerConfig struct {
	Addr string `toml:"addr"`
}

func newTcpServer(c *gconfig.Configuration) (*TcpServer, error) {
	server := &TcpServer{}
	conf := &tcpServerConfig{}
	if err := c.UnmarshalKeyWithParent(server.Type(), conf); err != nil {
		return nil, err
	}
	server.conf = conf
	if err := server.BindProc(c); err != nil {
		return nil, err
	}
	return server, nil
}

func (t *TcpServer) Init() error {
	return nil
}

func (t *TcpServer) Start(el session.EventHandler) error {
	var err error
	t.listener, err = net.Listen("tcp", t.conf.Addr)
	if err != nil {
		return err
	}
	t.BindHandler(el)
	go t.accept()
	return nil
}

func (t *TcpServer) accept() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				logger.Nlog.Warn("tcpserver", zap.String("msg", "accept temporary error"))
				time.Sleep(time.Millisecond)
				continue
			}
			break
		}
		sess := session.NewTcpSession(conn, t.SessionBundle)
		go sess.Start()
	}
}

func (t *TcpServer) Type() string {
	return PEER_TCP_SERVER
}

func (t *TcpServer) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTcpServer(c)
}

func (t *TcpServer) Stop() {

}

func init() {
	Register(&TcpServer{})
}
