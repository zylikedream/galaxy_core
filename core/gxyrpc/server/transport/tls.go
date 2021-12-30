package transport

import (
	"crypto/tls"

	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyregister"
)

type tlsConfig struct {
	Addr    string `toml:"addr"`
	PemFile string `toml:"pem_file"`
	KeyFile string `toml:"key_file"`
}

type tlsTransport struct {
	conf    *tlsConfig
	tlsConf *tls.Config
}

func newTlsTransport(c *gxyconfig.Configuration) (*tlsTransport, error) {
	tran := &tlsTransport{}
	conf := &tlsConfig{}
	if err := c.UnmarshalKey(tran.Type(), conf); err != nil {
		return nil, err
	}
	tran.conf = conf
	if err := tran.initTlsConfig(); err != nil {
		return nil, err
	}
	return tran, nil
}

func (t *tlsTransport) initTlsConfig() error {
	cert, err := tls.LoadX509KeyPair(t.conf.PemFile, t.conf.KeyFile)
	if err != nil {
		return err
	}
	t.tlsConf = &tls.Config{Certificates: []tls.Certificate{cert}}
	return nil
}

func (t *tlsTransport) Addr() string {
	return t.conf.Addr
}

func (t *tlsTransport) Network() string {
	return "tcp"
}

func (t *tlsTransport) Option() server.OptionFn {
	return server.WithTLSConfig(t.tlsConf)
}

func (t *tlsTransport) Type() string {
	return TRANSPORT_TYPE_TLS
}

func (t *tlsTransport) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	return newTlsTransport(c)
}

func init() {
	gxyregister.Register((*tlsTransport)(nil))
}
