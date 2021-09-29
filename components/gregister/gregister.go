package gregister

import (
	"fmt"

	"github.com/zylikedream/galaxy/components/gconfig"
)

type Builder interface {
	Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error)
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

func (r *register) NewNode(t string, c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	builder, ok := r.nodeMap[t]
	if !ok {
		return nil, fmt.Errorf("no node for type:%s", t)
	}
	return builder.Build(c)
}
