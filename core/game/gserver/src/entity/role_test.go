package entity

import (
	"bytes"
	"context"
	"testing"

	"github.com/zylikedream/galaxy/core/game/gserver/src/gscontext"
	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gxymongo"
)

var config = []byte(`
[mongo]
	addr = "mongodb://root:test@localhost:27017/admin"
	database = "test"
	pool_size = {min = 10, max = 20}
`)

func TestRoleEntity_Create(t *testing.T) {
	r := NewRoleEntity()
	ctx := gscontext.NewContext(context.Background())
	mgo, err := gxymongo.NewMongoClient(ctx, gconfig.NewWithReader(bytes.NewBuffer(config), gconfig.WithConfigType("toml")))
	if err != nil {
		t.Error(err)
		return
	}
	ctx.SetMongo(mgo)
	account := "zhangyi1"
	if err := r.Create(ctx, account); err != nil {
		t.Error(err)
		return
	}
	t.Logf("role = %+v", r)
	return
}

func TestRoleEntity_Load(t *testing.T) {
	r := NewRoleEntity()
	ctx := gscontext.NewContext(context.Background())
	mgo, err := gxymongo.NewMongoClient(ctx, gconfig.NewWithReader(bytes.NewBuffer(config), gconfig.WithConfigType("toml")))
	if err != nil {
		t.Error(err)
		return
	}
	ctx.SetMongo(mgo)
	account := "zhangyi1"
	if err := r.LoadByAccount(ctx, account); err != nil {
		t.Error(err)
		return
	}
	t.Logf("role = %#v", r)
	return
}
