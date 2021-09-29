package tcp

import (
	"net"
	"time"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network/peer"
	"github.com/zylikedream/galaxy/components/network/peer/processor"
)

type TcpListener struct {
	proc     *processor.Processor
	listener net.Listener

	addr string `mapstructure:"addr"`
}

func newTcpListener(c *gconfig.Configuration) (*TcpListener, error) {
	proc, err := processor.NewProcessor(c)
	if err != nil {
		return nil, err
	}
	tcplistener := &TcpListener{
		proc: proc,
	}
	if err := c.UnmarshalKey("network.tcp_acceptor", tcplistener); err != nil {
		return nil, err
	}
	return tcplistener, nil

}

func (t *TcpListener) Init() error {
	return nil
}

func (t *TcpListener) Start() error {
	var err error
	t.listener, err = net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	go t.accept()
	return nil
}

func (t *TcpListener) accept() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				time.Sleep(time.Millisecond)
				continue
			}
			break
		}
		sess := NewTcpSession(conn, t.proc)
		go sess.Start()
	}
}

func (t *TcpListener) Type() string {
	return peer.PEER_TCP_ACCEPTOR
}

func (t *TcpListener) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTcpListener(c)
}

func (t *TcpListener) Stop() {

}

func init() {
	peer.Register(&TcpListener{})
}
