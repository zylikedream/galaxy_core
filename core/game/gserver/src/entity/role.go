package entity

import (
	"context"
	"reflect"

	"github.com/zylikedream/galaxy/core/game/gserver/src/component"
	"github.com/zylikedream/galaxy/core/gmongo"
	"go.mongodb.org/mongo-driver/bson"
)

type RoleEntity struct {
	ID       uint64
	roleInfo *component.RoleInfo
}

func NewRoleEntity(enityID uint64) *RoleEntity {
	entity := &RoleEntity{
		ID: enityID,
	}
	return entity
}

func (r *RoleEntity) Load(ctx context.Context) error {
	components := make([]component.Component, 0)
	comType := reflect.TypeOf((*component.Component)(nil)).Elem()
	val := reflect.ValueOf(r).Elem()
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f.Addr().Type().Implements(comType) {
			components = append(components, f.Interface().(component.Component))
		}
	}
	for _, comp := range components {
		gmongo.FindOne(ctx, comp.Name(), bson.M{"_id": r.ID})
	}
	return nil
}
