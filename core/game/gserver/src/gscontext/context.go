package gscontext

import (
	"context"
	"fmt"
	"reflect"

	"github.com/zylikedream/galaxy/core/game/gserver/src/gsconfig"
	"github.com/zylikedream/galaxy/core/gxylog"
	"github.com/zylikedream/galaxy/core/gxymongo"
	"github.com/zylikedream/galaxy/core/gxynet/peer"
	"github.com/zylikedream/galaxy/core/gxynet/session"
)

// var _ context.Context = &Context{}

// Context is a rpcx customized Context that can contains multiple values.

type Context struct {
	cookie map[interface{}]interface{}
	context.Context
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
		cookie:  make(map[interface{}]interface{}),
	}
}

func (c *Context) Value(key interface{}) interface{} {
	if c.cookie == nil {
		c.cookie = make(map[interface{}]interface{})
	}

	if v, ok := c.cookie[key]; ok {
		return v
	}
	return c.Context.Value(key)
}

func (c *Context) SetValue(key, val interface{}) {
	if c.cookie == nil {
		c.cookie = make(map[interface{}]interface{})
	}
	c.cookie[key] = val
}

func (c *Context) String() string {
	return fmt.Sprintf("%v.WithValue(%v)", c.Context, c.cookie)
}

func (c *Context) Clone() *Context {
	return &Context{
		cookie:  c.cookie,
		Context: c.Context,
	}
}

func WithValue(parent context.Context, key, val interface{}) *Context {
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	tags := make(map[interface{}]interface{})
	tags[key] = val
	return &Context{Context: parent, cookie: tags}
}

func WithLocalValue(ctx *Context, key, val interface{}) *Context {
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}

	if ctx.cookie == nil {
		ctx.cookie = make(map[interface{}]interface{})
	}

	ctx.cookie[key] = val
	return ctx
}

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "rpcx context value " + k.name }

var (
	// RemoteConnContextKey is a context key. It can be used in
	// services with context.WithValue to access the connection arrived on.
	// The associated value will be of type net.Conn.
	sessionCtxKey    = &contextKey{"session"}
	mongoCtxKey      = &contextKey{"mongo"}
	loggerCtxKey     = &contextKey{"logger"}
	peerCtxKey       = &contextKey{"peer"}
	GameConfigCtxKey = &contextKey{"game_config"}
)

func (ctx *Context) GetSession() session.Session {
	return ctx.Value(sessionCtxKey).(session.Session)
}

func (ctx *Context) SetSession(sess session.Session) {
	ctx.SetValue(sessionCtxKey, sess)
}

func (ctx *Context) GetMongo() *gxymongo.MongoClient {
	return ctx.Value(mongoCtxKey).(*gxymongo.MongoClient)
}

func (ctx *Context) SetMongo(mgo *gxymongo.MongoClient) {
	ctx.SetValue(mongoCtxKey, mgo)
}

func (ctx *Context) GetLogger() *gxylog.GalaxyLog {
	return ctx.Value(loggerCtxKey).(*gxylog.GalaxyLog)
}

func (ctx *Context) SetLogger(log *gxylog.GalaxyLog) {
	ctx.SetValue(loggerCtxKey, log)
}

func (ctx *Context) GetPeer() peer.Peer {
	return ctx.Value(peerCtxKey).(peer.Peer)
}

func (ctx *Context) SetPeer(p peer.Peer) {
	ctx.SetValue(peerCtxKey, p)
}

func (ctx *Context) GetGameConfig() *gsconfig.GameConfig {
	return ctx.Value(GameConfigCtxKey).(*gsconfig.GameConfig)
}

func (ctx *Context) SetGameConfig(gs *gsconfig.GameConfig) {
	ctx.SetValue(GameConfigCtxKey, gs)
}
