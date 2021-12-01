package gscontext

import (
	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/gmongo"
	"github.com/zylikedream/galaxy/core/network/session"
)

// var _ context.Context = &Context{}

// Context is a rpcx customized Context that can contains multiple values.
type Context struct {
	*gcontext.Context
}

type contextKey struct {
	name string
}

func NewContext(ctx Context) *Context {
	return &Context{
		gcontext.NewContext(ctx),
	}
}

func (k *contextKey) String() string { return "rpcx context value " + k.name }

var (
	// RemoteConnContextKey is a context key. It can be used in
	// services with context.WithValue to access the connection arrived on.
	// The associated value will be of type net.Conn.
	sessionCtxKey = &contextKey{"session"}
	mongoCtxKey   = &contextKey{"mongo"}
	loggerCtxKey  = &contextKey{"logger"}
)

func (ctx *Context) GetSession() session.Session {
	return ctx.Value(sessionCtxKey).(session.Session)
}

func (ctx *Context) SetSession(sess session.Session) {
	ctx.SetValue(sessionCtxKey, sess)
}

func (ctx *Context) GetMongo() *gmongo.MongoClient {
	return ctx.Value(mongoCtxKey).(*gmongo.MongoClient)
}

func (ctx *Context) SetMongo(mgo *gmongo.MongoClient) {
	ctx.SetValue(mongoCtxKey, mgo)
}

func (ctx *Context) GetLogger() *glog.GalaxyLog {
	return ctx.Value(loggerCtxKey).(*glog.GalaxyLog)
}

func (ctx *Context) SetLogger(log *glog.GalaxyLog) {
	ctx.SetValue(loggerCtxKey, log)
}
