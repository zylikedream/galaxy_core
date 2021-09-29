package register

import (
	"fmt"

	"github.com/zylikedream/galaxy/components/gconfig"
)

type FuncType = func(c *gconfig.Configuration) (interface{}, error)

type register struct {
	nodeMap map[string]FuncType
}

func NewRegister() *register {
	return &register{
		nodeMap: map[string]FuncType{},
	}
}

func (r *register) Register(t string, nfun FuncType) {
	r.nodeMap[t] = nfun
}

func (r *register) NewNode(t string, c *gconfig.Configuration) (interface{}, error) {
	Func, ok := r.nodeMap[t]
	if !ok {
		return nil, fmt.Errorf("no node for type:%s", t)
	}
	return Func(c)
}
