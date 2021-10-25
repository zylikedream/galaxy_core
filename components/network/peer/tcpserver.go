package peer

import (
	"net"
	"time"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/handler"
	"github.com/zylikedream/galaxy/components/network/processor"
	"github.com/zylikedream/galaxy/components/network/session"
)

type TcpServer struct {
	proc     *processor.Processor
	listener net.Listener
	h        handler.Handler
	conf     *config
}

type config struct {
	addr string `toml:"addr"`
}

func newTcpServer(c *gconfig.Configuration) (*TcpServer, error) {
	server := &TcpServer{}
	conf := &config{}
	if err := c.UnmarshalKeyWithParent(server.Type(), conf); err != nil {
		return nil, err
	}
	server.conf = conf
	proc, err := processor.NewProcessor(c)
	if err != nil {
		return nil, err
	}
	server.proc = proc
	return server, nil
}

func (t *TcpServer) BindHandler(h handler.Handler) {
	t.h = h
}

func (t *TcpServer) Init() error {
	return nil
}

func (t *TcpServer) Start() error {
	var err error
	t.listener, err = net.Listen("tcp", t.conf.addr)
	if err != nil {
		return err
	}
	go t.accept()
	return nil
}

func (t *TcpServer) accept() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				time.Sleep(time.Millisecond)
				continue
			}
			break
		}
		sess := session.NewTcpSession(conn, t.proc)
		sess.BindHandler(t.h)
		go sess.Start()
	}
}

func (t *TcpServer) Type() string {
	return PEER_TCP_ACCEPTOR
}

func (t *TcpServer) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTcpServer(c)
}

func (t *TcpServer) Stop() {

}

func init() {
	Register(&TcpServer{})
}
