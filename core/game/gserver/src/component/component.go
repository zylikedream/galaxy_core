package component

import (
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AutoCreate interface {
	Create(ctx *gscontext.Context, ID primitive.ObjectID)
}
