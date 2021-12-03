package glog

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/glog/color"
	"github.com/zylikedream/galaxy/core/glog/corewriter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	Level           string `toml:"level"`            // 日志初始等级，默认info级别
	EnableAddCaller bool   `toml:"enable_addcaller"` // 是否添加调用者信息，默认不加调用者信息
	Writer          string `toml:"writer"`           // 使用哪种Writer，可选[rotate_file|stderr]，默认file
	CallerSkip      int    `toml:"caller_skip"`      // 跳过的堆栈层数，一般默认都为1
	Watch           bool   `toml:"watch"`            // 是否监听日志等级变化
}
type GalaxyLog struct {
	core          zapcore.Core
	name          string
	desugar       *zap.Logger
	logLevel      *zap.AtomicLevel
	conf          *config
	sugar         *zap.SugaredLogger
	asyncStopFunc func() error
}

func DefaultConfig() *config {
	return &config{
		Level:           "info",
		CallerSkip:      1,
		EnableAddCaller: false,
		Writer:          "rotate_file",
	}
}

func NewLogger(name string, configure *gconfig.Configuration, opts ...Option) *GalaxyLog {
	return newLogger(name, configure, opts...)
}

func newLogger(name string, configure *gconfig.Configuration, opts ...Option) *GalaxyLog {
	conf := DefaultConfig()
	gl := &GalaxyLog{
		conf:          conf,
		name:          name,
		asyncStopFunc: func() error { return nil },
	}
	for _, opt := range opts {
		opt(gl)
	}
	if err := configure.UnmarshalKey("log", conf); err != nil {
		panic(err)
	}
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if conf.EnableAddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(conf.CallerSkip))
	}

	// 默认日志级别
	logLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	gl.logLevel = &logLevel
	if err := gl.logLevel.UnmarshalText([]byte(conf.Level)); err != nil {
		panic(err)
	}

	// 如果用户没有设置core。那么就选择官方默认的core。
	if gl.core == nil {
		w, err := corewriter.NewCoreWriter(conf.Writer, configure, logLevel)
		if err != nil {
			panic(err)
		}
		gl.core = w
		gl.asyncStopFunc = w.Close
	}

	gl.desugar = zap.New(gl.core, zapOptions...)
	gl.sugar = gl.desugar.Sugar()

	if gl.conf.Watch {
		// 如果名字不为空，加载动态配置
		gl.AutoLevel(configure)
	}
	return gl

}

// AutoLevel ...
func (g *GalaxyLog) AutoLevel(c *gconfig.Configuration) {
	c.Watch(func(c *gconfig.Configuration) {
		lvText := strings.ToLower(c.GetString("log.level"))
		if lvText == "" {
			return
		}
		g.Infof("config level change %s->%s", g.logLevel.String(), lvText)
		_ = g.logLevel.UnmarshalText([]byte(lvText))
	})
}

// SetLevel ...
func (g *GalaxyLog) SetLevel(lv Level) {
	g.logLevel.SetLevel(lv)
}

// Flush ...
// When use os.Stdout or os.Stderr as zapcore.WriteSyncer
// g.desugar.Sync() maybe return an error like this: 'sync /dev/stdout: The handle is invalid.'
// Because os.Stdout and os.Stderr is a non-normal file, maybe not support 'fsync' in different os platform
// So ignored Sync() return value
// About issues: https://github.com/uber-go/zap/issues/328
// About 'fsync': https://man7.org/linux/man-pages/man2/fsync.2.html
func (g *GalaxyLog) Flush() error {
	if g.asyncStopFunc != nil {
		if err := g.asyncStopFunc(); err != nil {
			return err
		}
	}

	_ = g.desugar.Sync()
	return nil
}

// Debug ...
func (g *GalaxyLog) Debug(msg string, fields ...Field) {
	g.desugar.Debug(msg, fields...)
}

// Debugw ...
func (g *GalaxyLog) Debugw(msg string, keysAndValues ...interface{}) {
	g.sugar.Debugw(msg, keysAndValues...)
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}

// StdLog ...
func (g *GalaxyLog) StdLog() *log.Logger {
	return zap.NewStdLog(g.desugar)
}

// Debugf ...
func (g *GalaxyLog) Debugf(template string, args ...interface{}) {
	g.sugar.Debugw(sprintf(template, args...))
}

// Info ...
func (g *GalaxyLog) Info(msg string, fields ...Field) {
	g.desugar.Info(msg, fields...)
}

// Infow ...
func (g *GalaxyLog) Infow(msg string, keysAndValues ...interface{}) {
	g.sugar.Infow(msg, keysAndValues...)
}

// Infof ...
func (g *GalaxyLog) Infof(template string, args ...interface{}) {
	g.sugar.Infof(sprintf(template, args...))
}

// Warn ...
func (g *GalaxyLog) Warn(msg string, fields ...Field) {
	g.desugar.Warn(msg, fields...)
}

// Warnw ...
func (g *GalaxyLog) Warnw(msg string, keysAndValues ...interface{}) {
	g.sugar.Warnw(msg, keysAndValues...)
}

// Warnf ...
func (g *GalaxyLog) Warnf(template string, args ...interface{}) {
	g.sugar.Warnf(sprintf(template, args...))
}

// Error ...
func (g *GalaxyLog) Error(msg string, fields ...Field) {
	g.desugar.Error(msg, fields...)
}

// Errorw ...
func (g *GalaxyLog) Errorw(msg string, keysAndValues ...interface{}) {
	g.sugar.Errorw(msg, keysAndValues...)
}

// Errorf ...
func (g *GalaxyLog) Errorf(template string, args ...interface{}) {
	g.sugar.Errorf(sprintf(template, args...))
}

// Panic ...
func (g *GalaxyLog) Panic(msg string, fields ...Field) {
	panicDetail(msg, fields...)
	g.desugar.Panic(msg, fields...)
}

// Panicw ...
func (g *GalaxyLog) Panicw(msg string, keysAndValues ...interface{}) {
	g.sugar.Panicw(msg, keysAndValues...)
}

// Panicf ...
func (g *GalaxyLog) Panicf(template string, args ...interface{}) {
	g.sugar.Panicf(sprintf(template, args...))
}

// DPanic ...
func (g *GalaxyLog) DPanic(msg string, fields ...Field) {
	g.desugar.DPanic(msg, fields...)
}

// DPanicw ...
func (g *GalaxyLog) DPanicw(msg string, keysAndValues ...interface{}) {
	g.sugar.DPanicw(msg, keysAndValues...)
}

// DPanicf ...
func (g *GalaxyLog) DPanicf(template string, args ...interface{}) {
	g.sugar.DPanicf(sprintf(template, args...))
}

// Fatal ...
func (g *GalaxyLog) Fatal(msg string, fields ...Field) {
	g.desugar.Fatal(msg, fields...)
}

// Fatalw ...
func (g *GalaxyLog) Fatalw(msg string, keysAndValues ...interface{}) {
	g.sugar.Fatalw(msg, keysAndValues...)
}

// Fatalf ...
func (g *GalaxyLog) Fatalf(template string, args ...interface{}) {
	g.sugar.Fatalf(sprintf(template, args...))
}

func panicDetail(msg string, fields ...Field) {
	enc := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(enc)
	}

	// 控制台输出
	fmt.Printf("%s: \n    %s: %s\n", color.Red("panic"), color.Red("msg"), msg)
	if _, file, line, ok := runtime.Caller(3); ok {
		fmt.Printf("    %s: %s:%d\n", color.Red("loc"), file, line)
	}
	for key, val := range enc.Fields {
		fmt.Printf("    %s: %s\n", color.Red(key), fmt.Sprintf("%+v", val))
	}
}

// With ...
func (g *GalaxyLog) With(fields ...Field) *GalaxyLog {
	desugarLogger := g.desugar.With(fields...)
	return &GalaxyLog{
		desugar:  desugarLogger,
		logLevel: g.logLevel,
		sugar:    desugarLogger.Sugar(),
		conf:     g.conf,
	}
}

// WithCallerSkip ...
func (g *GalaxyLog) WithCallerSkip(callerSkip int, fields ...Field) *GalaxyLog {
	g.conf.CallerSkip = callerSkip
	desugarLogger := g.desugar.WithOptions(zap.AddCallerSkip(callerSkip)).With(fields...)
	return &GalaxyLog{
		desugar:  desugarLogger,
		logLevel: g.logLevel,
		sugar:    desugarLogger.Sugar(),
		conf:     g.conf,
	}
}
