package component

import (
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Component interface {
	Init(ctx *gscontext.Context) error
}

type IDCreatetor interface {
	CreateByID(ctx *gscontext.Context, ID primitive.ObjectID)
}
