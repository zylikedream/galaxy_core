package module

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"text/template"

	"github.com/ChimeraCoder/gojson"
	"github.com/zylikedream/galaxy/core/game/gserver/util"
)

type ModuleMeta struct {
	Info    ModuleInfo
	Methods map[string]*MethodMeta
	mod     reflect.Value
	modType reflect.Type
}

type ModuleInfo struct {
	Name    string
	PkgPath string
}

type MethodMeta struct {
	Info      MethodInfo
	Method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
}

type MethodInfo struct {
	Name      string
	ReqName   string
	Req       string
	ReplyName string
	Reply     string
}

var siTemplate = `package {{.PkgPath}}

type {{.Name}} struct{}
{{$name := .Name}}
{{range .Methods}}
{{.Req}}
{{.Reply}}
type (s *{{$name}}) {{.Name}}(ctx context.Context, arg *{{.ReqName}}, reply *{{.ReplyName}}) error {
	return nil
}
{{end}}
`

func (modi ModuleInfo) String() string {
	tpl := template.Must(template.New("service").Parse(siTemplate))
	var buf bytes.Buffer
	_ = tpl.Execute(&buf, modi)
	return buf.String()
}

var (
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
)

var gmodules map[string]*ModuleMeta = make(map[string]*ModuleMeta)
var gmethods map[string]string = make(map[string]string)

func Register(mod interface{}) error {

	mval := reflect.ValueOf(mod)
	mtyp := reflect.TypeOf(mod)
	mtypIdr := reflect.Indirect(mval).Type()
	modMeta := &ModuleMeta{
		mod:     mval,
		modType: mtyp,
	}
	modi := ModuleInfo{}
	modi.Name = mtypIdr.Name()
	pkg := mtypIdr.PkgPath()
	if strings.Index(pkg, ".") > 0 {
		pkg = pkg[strings.LastIndex(pkg, ".")+1:]
	}
	pkg = strings.ReplaceAll(filepath.Base(pkg), "-", "_")
	modi.PkgPath = pkg
	modMeta.Info = modi
	modMeta.Methods = suitableMethods(mtyp, pkg)
	for _, m := range modMeta.Methods {
		gmethods[m.ArgType.Name()] = modi.Name + "/" + m.ArgType.Name()
	}
	gmodules[modi.Name] = modMeta
	return nil
}

func suitableMethods(typ reflect.Type, PkgPath string) map[string]*MethodMeta {
	methods := make(map[string]*MethodMeta)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type

		if method.PkgPath != "" {
			continue
		}
		if mtype.NumIn() != 4 {
			continue
		}
		// First arg must be context.Context
		ctxType := mtype.In(1)
		if !ctxType.Implements(typeOfContext) {
			continue
		}

		// Second arg need not be a pointer.
		argType := mtype.In(2)
		if !util.IsExportedOrBuiltinType(argType) {
			continue
		}
		// Third arg must be a pointer.
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Ptr {
			continue
		}
		// Reply type must be exported.
		if !util.IsExportedOrBuiltinType(replyType) {
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			continue
		}

		methodi := MethodInfo{}
		methodi.Name = method.Name

		if argType.Kind() == reflect.Ptr {
			argType = argType.Elem()
		}
		replyType = replyType.Elem()

		methodi.ReqName = argType.Name()
		methodi.Req = generateTypeDefination(methodi.ReqName, PkgPath, generateJSON(argType))
		methodi.ReplyName = replyType.Name()
		methodi.Reply = generateTypeDefination(methodi.ReplyName, PkgPath, generateJSON(replyType))

		methods[method.Name] = &MethodMeta{
			Info:      methodi,
			Method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
	}
	return methods
}

func generateJSON(typ reflect.Type) string {
	v := reflect.New(typ).Interface()

	data, _ := json.Marshal(v)
	return string(data)
}

func generateTypeDefination(name, pkg string, jsonValue string) string {
	jsonValue = strings.TrimSpace(jsonValue)
	if jsonValue == "" || jsonValue == `""` {
		return ""
	}
	r := strings.NewReader(jsonValue)
	output, err := gojson.Generate(r, gojson.ParseJson, name, pkg, nil, false, true)
	if err != nil {
		return ""
	}
	rt := strings.ReplaceAll(string(output), "``", "")
	return strings.ReplaceAll(rt, "package "+pkg+"\n\n", "")
}

func HandleMessage(ctx context.Context, Arg interface{}) error {
	argName := reflect.TypeOf(Arg).Name()
	path, ok := gmethods[argName]
	if !ok {
		return fmt.Errorf("no hanlder for message %s", argName)
	}
	names := strings.Split(path, "/")
	modName := names[0]
	methodName := names[1]
	mod := gmodules[modName]
	mtd := mod.Methods[methodName]
	Reply := reflect.New(mtd.ReplyType)
	// todo
	var err error
	if mtd.ArgType.Kind() != reflect.Ptr {
		err = mod.call(ctx, mtd, reflect.ValueOf(Arg).Elem(), Reply)
	} else {
		err = mod.call(ctx, mtd, reflect.ValueOf(Arg), Reply)
	}
	return err

}

func (m *ModuleMeta) call(ctx context.Context, mm *MethodMeta, argv, replyv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("[module internal error]: %v, method: %s, argv: %+v, stack: %s",
				r, mm.Method.Name, argv.Interface(), buf)
		}
	}()

	function := mm.Method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{m.mod, reflect.ValueOf(ctx), argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}

	return nil
}
