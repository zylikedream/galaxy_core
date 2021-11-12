package server

import (
	"github.com/smallnest/rpcx/server"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxyrpc/server/registery"
	"github.com/zylikedream/galaxy/core/gxyrpc/server/transport"
)

type serverConfig struct {
	Transport string `toml:"transport"`
	Registery string `toml:"registery"`
}

type gxyrpcServer struct {
	conf   *serverConfig
	trans  transport.Transport
	regist registery.Registery
	svr    *server.Server
}

type GxyrpcService interface {
	Service() interface{}
	Name() string
	Meta() string
}

func NewGrpcServer(configFile string) (*gxyrpcServer, error) {
	conf := &serverConfig{}
	configure := gconfig.New(configFile)
	if err := configure.UnmarshalKey("gxyrpc_server", conf); err != nil {
		return nil, err
	}
	gxyrpc := &gxyrpcServer{
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

func (g *gxyrpcServer) ReigsterService(service GxyrpcService) error {
	return g.svr.RegisterName(service.Name(), service.Service(), service.Meta())
}

func (g *gxyrpcServer) Start() error {
	if err := g.regist.Start(); err != nil {
		return err
	}
	if err := g.svr.Serve(g.trans.Network(), g.trans.Addr()); err != nil {
		return err
	}
	return nil
}

func (g *gxyrpcServer) Stop() error {
	return g.svr.Close()
}
