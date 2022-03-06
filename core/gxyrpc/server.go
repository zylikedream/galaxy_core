package gxyrpc

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxyrpc/server/registery"
	"github.com/zylikedream/galaxy/core/gxyrpc/server/transport"
)

type serverConfig struct {
	Transport string `toml:"transport"`
	Registery string `toml:"registery"`
}

type RpcServer struct {
	conf   *serverConfig
	trans  transport.Transport
	regist registery.Registery
	svr    *server.Server
}

type serviceOptionFunc func(opt *serviceOption)
type serviceOption struct {
	Name string
	Meta string
}

func WithName(name string) serviceOptionFunc {
	return func(opt *serviceOption) {
		opt.Name = name
	}
}

func WithMeta(meta string) serviceOptionFunc {
	return func(opt *serviceOption) {
		opt.Meta = meta
	}
}

func NewGrpcServer(configFile string) (*RpcServer, error) {
	conf := &serverConfig{}
	configure := gxyconfig.New(configFile)
	if err := configure.UnmarshalKey("gxyrpc_server", conf); err != nil {
		return nil, err
	}
	gxyrpc := &RpcServer{
		conf: conf,
	}
	trans, err := transport.NewTransport(conf.Transport, configure)
	if err != nil {
		return nil, err
	}
	regist, err := registery.NewRegistery(conf.Registery, registery.RegisterAddr(trans), configure)
	if err != nil {
		return nil, err
	}
	gxyrpc.trans = trans
	gxyrpc.regist = regist
	gxyrpc.svr = server.NewServer(trans.Option())

	plug := regist.GetPlugin()
	if plug != nil {
		gxyrpc.svr.Plugins.Add(plug)
	}
	return gxyrpc, nil
}

func (g *RpcServer) ReigsterService(service interface{}, opts ...serviceOptionFunc) error {
	defaultOpt := &serviceOption{}
	for _, opt := range opts {
		opt(defaultOpt)
	}
	if defaultOpt.Name == "" {
		return g.svr.Register(service, defaultOpt.Meta)
	} else {
		return g.svr.RegisterName(defaultOpt.Name, service, defaultOpt.Meta)
	}
}

func (g *RpcServer) RegisterAgent(path string, agent server.ServiceAgent) {
	g.svr.RegisterAgent(path, agent)
}

func (g *RpcServer) Start() error {
	if err := g.regist.Start(); err != nil {
		return err
	}
	if err := g.svr.Serve(g.trans.Network(), g.trans.Addr()); err != nil {
		return err
	}
	return nil
}

func (g *RpcServer) GetRawServer() *server.Server {
	return g.svr
}

func (g *RpcServer) Stop() error {
	return g.svr.Close()
}
