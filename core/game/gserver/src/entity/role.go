package entity

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/pkg/errors"
	"github.com/zylikedream/galaxy/core/game/gserver/src/component"
	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type compoentField struct {
	tableName  string
	autoload   bool
	autocreate bool
}

var compFields = make(map[reflect.Type]compoentField)

func parseTag(tag string) (string, []string) {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return "", nil
	}
	return tags[0], tags[1:]
}

func init() {
	entity := &RoleEntity{}
	typ := reflect.TypeOf(entity).Elem()
	comType := reflect.TypeOf((*component.Component)(nil)).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if !typ.Field(i).Type.Implements(comType) {
			continue
		}
		t := typ.Field(i).Tag.Get("table")
		if t == "" {
			continue
		}
		tableName, opts := parseTag(t)
		fmt.Println(typ.Field(i).Name, tableName, opts)
		compFields[typ.Field(i).Type.Elem()] = compoentField{
			tableName:  tableName,
			autoload:   arrutil.Contains(opts, "autoload"),
			autocreate: arrutil.Contains(opts, "autocreate"),
		}
	}
}

func getComponetTable(comp interface{}) string {
	return compFields[reflect.TypeOf(comp).Elem()].tableName
}

type RoleEntity struct {
	RoleID primitive.ObjectID
	acc    *component.RoleAccount `table:"account"`
	basic  *component.RoleBasic   `table:"role_basic,autoload,autocreate"`
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
	return r.load(ctx, bson.M{"account": acc})
}

func (r *RoleEntity) LoadByID(ctx *gscontext.Context, roleid uint64) error {
	return r.load(ctx, bson.M{"roleid": roleid})
}

func (r *RoleEntity) load(ctx *gscontext.Context, filter interface{}) error {
	gmongo := ctx.GetMongo()
	roleInfo := &component.RoleAccount{}
	err := gmongo.FindOne(ctx, roleInfo, getComponetTable(roleInfo), filter)
	if err != nil {
		return err
	}
	r.acc = roleInfo
	r.RoleID = roleInfo.RoleID
	if err := r.autoload(ctx); err != nil {
		return err
	}
	return nil
}

func (r *RoleEntity) autocreate(ctx *gscontext.Context, comp reflect.Value) {
	if ac, ok := comp.Interface().(component.AutoCreate); ok {
		ac.Create(ctx, r.RoleID)
	}
}

func (r *RoleEntity) autoload(ctx *gscontext.Context) error {
	val := reflect.ValueOf(r).Elem()
	autoload := make([]reflect.Value, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if cf, ok := compFields[f.Type().Elem()]; ok && cf.autoload {
			autoload = append(autoload, f)
		}
	}
	gmongo := ctx.GetMongo()
	for _, comp := range autoload {
		compIns := reflect.New(comp.Elem().Type())
		err := gmongo.FindOne(ctx, compIns, getComponetTable(compIns), bson.M{"_id": r.RoleID})
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				if cf, ok := compFields[comp.Type().Elem()]; ok && cf.autocreate {
					r.autocreate(ctx, compIns)
				}
			} else {
				return err
			}
		}
		comp.Set(compIns)
	}
	return nil
}

func (r *RoleEntity) Create(ctx *gscontext.Context, account string) error {
	gmongo := ctx.GetMongo()
	accInfo := &component.RoleAccount{
		Account: account,
	}
	_, err := gmongo.InsertOne(ctx, getComponetTable(accInfo), accInfo)
	if err != nil {
		return errors.Wrap(err, "create new role failed")
	}
	err = r.LoadByAccount(ctx, account)
	if err != nil {
		return err
	}

	return nil
}
