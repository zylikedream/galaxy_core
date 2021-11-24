package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"text/template"

	"github.com/ChimeraCoder/gojson"
	"github.com/zylikedream/galaxy/core/game/gclient/define"
	"github.com/zylikedream/galaxy/core/game/gclient/util"
	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/network/session"
	"go.uber.org/zap"
)

const (
	ACK_CODE_OK = iota
	ACK_CODE_FAIL
)

type IModule interface {
	// 定义handler是为了避免一些module 会有和导出类型相似的方法被意外导出，这时可以实现该方法，导出一个子类型handler, 一般的module可以不用管
	Handler(IModule) interface{}
	BeforeMsg(ctx gcontext.Context, msg interface{}) error // 像一些玩法的开启验证，和玩家的验证可以在这儿做
	AfterMsg(ctx gcontext.Context, msg interface{}) error
}

type BaseModule struct {
}

func (*BaseModule) BeforeMsg(ctx gcontext.Context, msg interface{}) error {
	return nil
}

func (b *BaseModule) Handler(mod IModule) interface{} {
	return mod
}

func (*BaseModule) AfterMsg(ctx gcontext.Context, msg interface{}) error {
	return nil
}

type Cookie struct {
	Sess session.Session
	Role interface{}
}

type ModuleMeta struct {
	Info    ModuleInfo
	Methods map[string]*MethodMeta
	modType reflect.Type
	im      IModule
}

type ModuleInfo struct {
	Name    string
	PkgPath string
}

type MethodMeta struct {
	Info    MethodInfo
	Method  reflect.Method
	MsgType reflect.Type
}

type MethodInfo struct {
	Name      string
	ReplyName string
	Reply     string
}

type NilReply struct {
}

var siTemplate = `package {{.Info.PkgPath}}

type {{.Info.Name}} struct{}
{{$name := .Info.Name}}
{{range .Methods}}
{{.Info.Req}}
{{.Info.Reply}}
type (s *{{$name}}) {{.Info.Name}}(ctx context.Context, arg *{{.Info.ReqName}}, reply *{{.Info.ReplyName}}) error {
	return nil
}
{{end}}
`

func (mm ModuleMeta) String() string {
	tpl := template.Must(template.New("service").Parse(siTemplate))
	var buf bytes.Buffer
	_ = tpl.Execute(&buf, mm)
	return buf.String()
}

var (
	typeOfError    = reflect.TypeOf((*error)(nil)).Elem()
	typeOfGContext = reflect.TypeOf((*gcontext.Context)(nil)).Elem()
)

type RouteInfo struct {
	ModName    string
	MethodName string
}

var gmodules map[string]*ModuleMeta = make(map[string]*ModuleMeta)
var groutes map[string]RouteInfo = make(map[string]RouteInfo)

func Register(mod IModule) error {
	handler := mod.Handler(mod)
	mval := reflect.ValueOf(handler)
	mtyp := reflect.TypeOf(handler)
	mtypIdr := reflect.Indirect(mval).Type()
	modMeta := &ModuleMeta{
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
		groutes[m.MsgType.Name()] = RouteInfo{
			ModName:    modi.Name,
			MethodName: m.MsgType.Name(),
		}
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
		if mtype.NumIn() != 3 {
			continue
		}
		// First arg must be context.Context
		ctxType := mtype.In(1)
		if ctxType != typeOfGContext {
			continue
		}

		// Second arg need be a pointer.
		replyType := mtype.In(2)
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

		replyType = replyType.Elem()

		methodi.ReplyName = replyType.Name()
		methodi.Reply = generateTypeDefination(methodi.ReplyName, PkgPath, generateJSON(replyType))

		methods[method.Name] = &MethodMeta{
			Info:    methodi,
			Method:  method,
			MsgType: replyType,
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

func HandleMessage(ctx gcontext.Context, Msg interface{}) error {
	argName := reflect.TypeOf(Msg).Name()
	path, ok := groutes[argName]
	if !ok {
		return fmt.Errorf("no hanlder for message %s", argName)
	}
	mod := gmodules[path.ModName]
	mtd := mod.Methods[path.MethodName]
	var err error
	if err = mod.im.BeforeMsg(ctx, Msg); err != nil {
		return err
	}
	if mtd.MsgType.Kind() != reflect.Ptr {
		err = mod.call(ctx, mtd, reflect.ValueOf(Msg).Elem())
	} else {
		err = mod.call(ctx, mtd, reflect.ValueOf(Msg))
	}
	if err != nil {
		return err
	}
	if err = mod.im.AfterMsg(ctx, Msg); err != nil {
		return err
	}
	return nil
}

func Send(ctx gcontext.Context, msg interface{}) {
	sess := ctx.Value(define.SessionCtxKey).(session.Session)
	err := sess.Send(msg)
	if err != nil {
		glog.Error("send error", zap.Error(err), zap.Any("msg", msg))
	}
}

func (m *ModuleMeta) call(ctx gcontext.Context, mm *MethodMeta, argv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]

			err = fmt.Errorf("module internal error: %v, method: %s, argv: %+v, stack: %s",
				r, mm.Method.Name, argv.Interface(), buf)
		}
	}()

	function := mm.Method.Func
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{reflect.ValueOf(m.im), reflect.ValueOf(ctx), argv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}

	return nil
}
