package component

import (
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleSign struct {
	RoleID   primitive.ObjectID `json:"_id"`
	SignTime int64              `json:"sign_time"`
	SignDay  int                `json:"sign_day"`
}

func (r *RoleSign) CreateByID(ctx *gscontext.Context, ID primitive.ObjectID) {
	r.RoleID = ID
	r.SignTime = 0
	r.SignDay = 0
}
