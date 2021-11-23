package gcontext

import (
	"context"
	"fmt"
	"reflect"
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
