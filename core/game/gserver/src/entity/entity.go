package entity

import "github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"

type Entity interface {
	Load(ctx *gscontext.Context) error
}
