package server

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxyrpc/server/registery"
	"github.com/zylikedream/galaxy/core/gxyrpc/server/transport"
)

type serverConfig struct {
	Network   string `toml:"network"`
	Registery string `toml:"registery"`
}

type gxyrpcServer struct {
	conf   *serverConfig
	trans  transport.Transport
	regist registery.Registery
	svr    *server.Server
}

type GxyrpcService interface {
	Name() string
	Meta() string
}

func NewGrpcServer(configFile string) (*gxyrpcServer, error) {
	conf := &serverConfig{}
	configure := gconfig.New(configFile)
	if err := configure.UnmarshalKey("gxyrpc", conf); err != nil {
		return nil, err
	}
	gxyrpc := &gxyrpcServer{
		conf: conf,
	}
	trans, err := transport.NewTransport(conf.Network, configure)
	if err != nil {
		return nil, err
	}
	regist, err := registery.NewRegistery(registery.RegisterAddr(trans), conf.Registery, configure)
	if err != nil {
		return nil, err
	}
	gxyrpc.trans = trans
	gxyrpc.regist = regist
	gxyrpc.svr = server.NewServer(trans.Option())
	return gxyrpc, nil
}

func (g *gxyrpcServer) Start(services ...GxyrpcService) error {
	plug := g.regist.GetPlugin()
	if plug != nil {
		g.svr.Plugins.Add(plug)
	}
	if err := g.regist.Start(); err != nil {
		return err
	}
	for _, service := range services {
		g.svr.RegisterName(service.Name(), service, service.Meta())
	}
	if err := g.svr.Serve(g.trans.Network(), g.trans.Addr()); err != nil {
		return err
	}
	return nil
}

func (g *gxyrpcServer) Stop() error {
	return g.svr.Close()
}
