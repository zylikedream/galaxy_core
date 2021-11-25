package module

import (
	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/network/session"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "rpcx context value " + k.name }

var (
	// RemoteConnContextKey is a context key. It can be used in
	// services with context.WithValue to access the connection arrived on.
	// The associated value will be of type net.Conn.
	SessionCtxKey = &contextKey{"session"}
)

func GetSessionFromCtx(ctx gcontext.Context) session.Session {
	return ctx.Value(SessionCtxKey).(session.Session)
}
