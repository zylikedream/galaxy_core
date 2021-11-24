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
	"github.com/zylikedream/galaxy/core/game/gserver/define"
	"github.com/zylikedream/galaxy/core/game/gserver/util"
	"github.com/zylikedream/galaxy/core/game/proto"
	"github.com/zylikedream/galaxy/core/gcontext"
	"github.com/zylikedream/galaxy/core/glog"
	"github.com/zylikedream/galaxy/core/network/message"
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
		im:      mod,
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
		groutes[m.Info.ReqName] = RouteInfo{
			ModName:    modi.Name,
			MethodName: m.Info.Name,
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
		if mtype.NumIn() != 4 {
			continue
		}
		// First arg must be context.Context
		ctxType := mtype.In(1)
		if ctxType != typeOfGContext {
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

		argTypeV := argType
		if argTypeV.Kind() == reflect.Ptr {
			argTypeV = argTypeV.Elem()
		}
		replyType = replyType.Elem()

		methodi.ReqName = argTypeV.Name()
		methodi.Req = generateTypeDefination(methodi.ReqName, PkgPath, generateJSON(argTypeV))
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

func HandleMessage(ctx gcontext.Context, Msg interface{}) error {
	argName := reflect.TypeOf(Msg).Elem().Name()
	path, ok := groutes[argName]
	if !ok {
		return fmt.Errorf("no hanlder for message: %+v", Msg)
	}
	mod := gmodules[path.ModName]
	mtd := mod.Methods[path.MethodName]
	Reply := reflect.New(mtd.ReplyType)
	var err error
	defer func() {
		if err != nil {
			AckFail(ctx, Reply.Interface(), err.Error())
		} else {
			AckOk(ctx, Reply.Interface())
		}
	}()
	if err = mod.im.BeforeMsg(ctx, Msg); err != nil {
		return err
	}
	glog.Debugf("msg %#v", Msg)
	if mtd.ArgType.Kind() != reflect.Ptr {
		err = mod.call(ctx, mtd, reflect.ValueOf(Msg).Elem(), Reply)
	} else {
		err = mod.call(ctx, mtd, reflect.ValueOf(Msg), Reply)
	}
	if err != nil {
		return err
	}
	if err = mod.im.AfterMsg(ctx, Msg); err != nil {
		return err
	}
	return nil
}

func AckFail(ctx gcontext.Context, msg interface{}, Reason string) {
	Ack(ctx, msg, ACK_CODE_FAIL, Reason)
}

func AckOk(ctx gcontext.Context, msg interface{}) {
	Ack(ctx, msg, ACK_CODE_OK, "")
}

func Ack(ctx gcontext.Context, msg interface{}, code int, Reason string) {
	meta := message.MessageMetaByMsg(msg)
	if meta == nil {
		glog.Error("ack unkonw msg", zap.Any("msg", msg))
		return
	}

	ack := &proto.Ack{
		Code:   code,
		Reason: Reason,
		MsgID:  meta.ID,
	}
	if code == ACK_CODE_OK {
		sess := ctx.Value(define.SessionCtxKey).(session.Session)
		msgData, err := sess.GetMessageCodec().Encode(msg)
		if err != nil {
			glog.Error("ack error", zap.Error(err))
			return
		}
		ack.Data = msgData
	}
	Send(ctx, ack)
}

func Send(ctx gcontext.Context, msg interface{}) {
	sess := ctx.Value(define.SessionCtxKey).(session.Session)
	err := sess.Send(msg)
	if err != nil {
		glog.Error("send error", zap.Error(err), zap.Any("msg", msg))
	}
}

func (m *ModuleMeta) call(ctx gcontext.Context, mm *MethodMeta, argv, replyv reflect.Value) (err error) {
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
	returnValues := function.Call([]reflect.Value{reflect.ValueOf(m.im), reflect.ValueOf(ctx), argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	if errInter != nil {
		return errInter.(error)
	}

	return nil
}
