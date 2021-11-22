package logic

import (
	"context"
	"reflect"
	"testing"
)

type module struct {
}

type Echo struct {
	msg string
}

func (m *module) TestFunc(ctx context.Context, e *Echo) {

}

func TestReflect(t *testing.T) {
	m := &module{}
	val := reflect.ValueOf(m)
	typ := reflect.TypeOf(m)
	for m := 0; m < val.NumMethod(); m++ {
		method := typ.Method(m)
		t.Log("***", method.PkgPath)
	}
}
