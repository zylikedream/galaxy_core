package tcp

import (
	"net"
	"sync"
	"time"

	log "github.com/AlexStocks/log4go"
	gxnet "github.com/dubbogo/gost/net"
	gxtime "github.com/dubbogo/gost/time"
	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/network"
	"github.com/zylikedream/galaxy/components/network/peer"
	"github.com/zylikedream/galaxy/components/network/peer/processor"
	"github.com/zylikedream/galaxy/components/network/session"
)

var (
	errSelfConnect        = errors.New("connect self!")
	serverFastFailTimeout = time.Second * 1
	serverID              = 0
	wheel                 *gxtime.Wheel
)

type TcpServer struct {
	proc     *processor.Processor
	listener net.Listener
	conf     *config
	wg       sync.WaitGroup
	done     chan struct{}
}

type sslConfig struct {
	Enable bool `toml:"enable"`
}

type config struct {
	addr string    `toml:"addr"`
	ssl  sslConfig `toml:"ssl"`
}

func newTcpServer(c *gconfig.Configuration) (*TcpServer, error) {
	tcplistener := &TcpServer{}
	conf := &config{}
	if err := c.UnmarshalKeyWithParent(tcplistener.Type(), conf); err != nil {
		return nil, err
	}
	tcplistener.conf = conf
	proc, err := processor.NewProcessor(c)
	if err != nil {
		return nil, err
	}
	tcplistener.proc = proc
	if err := c.UnmarshalKey("network.tcp_acceptor", tcplistener); err != nil {
		return nil, err
	}
	return tcplistener, nil

}

func (t *TcpServer) Init() error {
	return nil
}

func (t *TcpServer) Start() error {
	if err := t.listen(); err != nil {
		return err
	}
	go t.run()
	return nil
}

func (t *TcpServer) listen() error {
	var err error
	t.listener, err = net.Listen("tcp", t.conf.addr)
	if err != nil {
		return err
	}
	addr := t.conf.addr
	if t.conf.ssl.Enable {
		// if sslConfig, err := s.tlsConfigBuilder.BuildTlsConfig(); err == nil && sslConfig != nil {
		// 	t.listener, err = tls.Listen("tcp", t.conf.addr, sslConfig)
		// }
	} else {
		t.listener, err = net.Listen("tcp", addr)
	}
	if err != nil {
		return errors.Wrapf(err, "net.Listen(tcp, addr:%s)", addr)
	}
	return nil
}

func (t *TcpServer) run() {
	t.wg.Add(1)
	var newSession network.NewSessionCallback
	go func() {
		defer t.wg.Done()
		var err error
		var client network.Session
		var delay time.Duration
		for {
			if t.IsClosed() {
				return
			}
			if delay != 0 {
				<-wheel.After(delay)
			}
			client, err = t.accept(newSession)
			if err != nil {
				if netErr, ok := errors.Cause(err).(net.Error); ok && netErr.Temporary() {
					if delay == 0 {
						delay = 5 * time.Millisecond
					} else {
						delay *= 2
					}
					if max := 1 * time.Second; delay > max {
						delay = max
					}
					continue
				}
				continue
			}
			delay = 0
			client.Run()
		}
	}()
}

func (s *TcpServer) accept(newSession network.NewSessionCallback) (network.Session, error) {
	conn, err := s.listener.Accept()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if gxnet.IsSameAddr(conn.RemoteAddr(), conn.LocalAddr()) {
		log.Warn("conn.localAddr{%s} == conn.RemoteAddr", conn.LocalAddr().String(), conn.RemoteAddr().String())
		return nil, errors.WithStack(errSelfConnect)
	}

	ss := session.NewTCPSession(conn, s)
	err = newSession(ss)
	if err != nil {
		conn.Close()
		return nil, errors.WithStack(err)
	}

	return ss, nil
}

func (t *TcpServer) Type() string {
	return peer.PEER_TCP_ACCEPTOR
}

func (t *TcpServer) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTcpServer(c)
}

func (t *TcpServer) IsClosed() bool {
	select {
	case <-t.done:
		return true
	default:
		return false
	}
}

func (t *TcpServer) Stop() {

}

func init() {
	peer.Register(&TcpServer{})
}
