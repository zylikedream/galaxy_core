package gxyregister

import (
	"errors"
	"fmt"

	"github.com/zylikedream/galaxy/core/gxyconfig"
)

var (
	ErrParamNotEnough = errors.New("param is not enough")
	ErrParamErrType   = errors.New("param type error")
)

type Builder interface {
	Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error)
	Type() string
}

type register struct {
	nodeMap map[string]Builder
}

func NewRegister() *register {
	return &register{
		nodeMap: map[string]Builder{},
	}
}

func (r *register) Register(b Builder) {
	r.nodeMap[b.Type()] = b
}

func (r *register) NewNode(t string, c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	builder, ok := r.nodeMap[t]
	if !ok {
		return nil, fmt.Errorf("no node for type:%s", t)
	}
	return builder.Build(c, args...)
}

var gxyregister *register

func init() {
	gxyregister = NewRegister()
}

func Register(b Builder) {
	gxyregister.Register(b)
}

func NewNode(t string, c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	builder, ok := gxyregister.nodeMap[t]
	if !ok {
		return nil, fmt.Errorf("no node for type:%s", t)
	}
	return builder.Build(c, args...)
}
