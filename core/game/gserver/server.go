package main

import (
	"context"

	"github.com/zylikedream/galaxy/core/game/gserver/src/app"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	_ "github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/gxylog"
	"go.uber.org/zap"
)

func main() {
	ctx := gscontext.NewContext(context.Background())
	s := app.NewServer(ctx)
	if err := s.Run(ctx); err != nil {
		gxylog.Error("server run err", zap.Error(err))
	}
}
