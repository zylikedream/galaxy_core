package module

import (
	"reflect"
	"testing"
)

func Test_suitableMethods(t *testing.T) {
	lm := &LoginModule{}
	Register(lm)
	name := reflect.TypeOf(lm).Elem().Name()
	t.Log(gmodules[name].String())
	t.Log(len(gmodules[name].Methods))
}
