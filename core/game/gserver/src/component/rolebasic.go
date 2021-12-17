package component

import (
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleBasic struct {
	RoleID primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
}

func (r *RoleBasic) CreateByID(ctx *gscontext.Context, ID primitive.ObjectID) {
	r.RoleID = ID
	r.Name = ID.String()
}

func (r *RoleBasic) Init() error {
	return nil
}
