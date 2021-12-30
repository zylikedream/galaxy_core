package app

import (
	"context"

	"github.com/zylikedream/galaxy/core/game/gserver/src/gsconfig"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/game/gserver/src/logic"
	_ "github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gmongo"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/network"
	"github.com/zylikedream/galaxy/core/network/peer"
	"go.uber.org/zap"
)

type Server struct {
	p          peer.Peer
	logger     *gxylog.GalaxyLog
	mgoCli     *gmongo.MongoClient
	gameConfig *gsconfig.GameConfig
}

func NewServer(ctx *gscontext.Context) *Server {
	svr := &Server{}
	if err := svr.Init(ctx); err != nil {
		panic(err)
	}
	return svr
}

func (s *Server) Init(ctx *gscontext.Context) error {
	p, err := network.NewNetwork(gconfig.New("config/network.toml"))
	if err != nil {
		return err
	}
	s.p = p

	s.logger = gxylog.NewLogger("server", gconfig.New("config/log.toml"))
	gxylog.SetDefaultLogger(s.logger)

	cli, err := gmongo.NewMongoClient(ctx, gconfig.New("config/mongo.toml"))
	if err != nil {
		return err
	}
	s.mgoCli = cli

	s.gameConfig, err = gsconfig.NewGameConfig()
	if err != nil {
		return err
	}

	return nil

}

func (s *Server) Run(ctx *gscontext.Context) error {
	ctx.SetLogger(s.logger)
	ctx.SetMongo(s.mgoCli)
	ctx.SetPeer(s.p)
	ctx.SetGameConfig(s.gameConfig)

	if err := s.p.Start(ctx, &logic.LogicHandle{}); err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := gscontext.NewContext(context.Background())
	s := NewServer(ctx)
	if err := s.Run(ctx); err != nil {
		gxylog.Error("server run err", zap.Error(err))
	}
}
