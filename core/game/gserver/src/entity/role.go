package entity

import (
	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/core/game/gserver/src/component"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleEntity struct {
	RoleID   primitive.ObjectID
	roleInfo *component.RoleInfo
}

func NewRoleEntity() *RoleEntity {
	return &RoleEntity{}
}

// func (r *RoleEntity) Load(ctx context.Context) error {
// 	components := make([]component.Component, 0)
// 	comType := reflect.TypeOf((*component.Component)(nil)).Elem()
// 	val := reflect.ValueOf(r).Elem()
// 	for i := 0; i < val.NumField(); i++ {
// 		f := val.Field(i)
// 		if f.Addr().Type().Implements(comType) {
// 			components = append(components, f.Interface().(component.Component))
// 		}
// 	}
// 	for _, comp := range components {
// 		fmt.Println(comp.Name())
// 	}
// 	return nil
// }

func (r *RoleEntity) LoadByAccount(ctx *gscontext.Context, acc string) error {
	gmongo := ctx.GetMongo()
	roleInfo := &component.RoleInfo{}
	err := gmongo.FindOne(ctx, roleInfo, roleInfo.GetName(), bson.M{"account": acc})
	if err != nil {
		return err
	}
	r.roleInfo = roleInfo
	r.RoleID = roleInfo.RoleID
	return nil

}

func (r *RoleEntity) LoadByID(ctx *gscontext.Context, roleid uint64) error {
	gmongo := ctx.GetMongo()
	roleInfo := &component.RoleInfo{}
	err := gmongo.FindOne(ctx, roleInfo, roleInfo.GetName(), bson.M{"roleid": roleid})
	if err != nil {
		return err
	}
	r.roleInfo = roleInfo
	r.RoleID = roleInfo.RoleID
	return nil
}

func (r *RoleEntity) Create(ctx *gscontext.Context, account string) error {
	gmongo := ctx.GetMongo()
	roleInfo := &component.RoleInfo{
		Account: account,
		Name:    account,
	}
	res, err := gmongo.InsertOne(ctx, roleInfo.GetName(), roleInfo)
	if err != nil {
		return errors.Wrap(err, "create new role failed")
	}
	roleInfo.RoleID = res.InsertedID.(primitive.ObjectID)
	r.roleInfo = roleInfo
	return nil
}
