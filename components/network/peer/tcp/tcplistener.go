package tcp

import (
	"net"
	"time"

	"github.com/spf13/viper"
)

type TcpListener struct {
	addr     string
	listener net.Listener
}

func (t *TcpListener) Init(v *viper.Viper) error {
	t.addr = v.GetString("addr")
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
		sess := NewTcpSession(conn)
		go sess.Start()
	}
}
