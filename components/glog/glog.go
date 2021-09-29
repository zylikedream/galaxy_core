package glog

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/glog/color"
	"github.com/zylikedream/galaxy/components/glog/corewriter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	Level           string `mapstructure:"level"`            // 日志初始等级，默认info级别
	EnableAddCaller bool   `mapstructure:"enable_addcaller"` // 是否添加调用者信息，默认不加调用者信息
	Writer          string `mapstructure:"writer"`           // 使用哪种Writer，可选[rotate_file|stderr]，默认file
	CallerSkip      int    `mapstructure:"caller_skip"`      // 跳过的堆栈层数，一般默认都为1
	Watch           bool   `mapstructure:"watch"`            // 是否监听日志等级变化
}
type GalaxyLog struct {
	core          zapcore.Core
	name          string
	desugar       *zap.Logger
	lv            zap.AtomicLevel
	conf          *config
	sugar         *zap.SugaredLogger
	asyncStopFunc func() error
}

func DefaultConfig() *config {
	return &config{
		Level:           "info",
		CallerSkip:      1,
		EnableAddCaller: false,
		Writer:          "file",
	}
}

func NewLogger(name string, configFile string, opts ...Option) *GalaxyLog {
	configure := gconfig.New(configFile)
	conf := DefaultConfig()
	gl := &GalaxyLog{
		conf:          DefaultConfig(),
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
	gl.lv = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := gl.lv.UnmarshalText([]byte(conf.Level)); err != nil {
		panic(err)
	}

	// 如果用户没有设置core。那么就选择官方默认的core。
	if gl.core == nil {
		w, err := corewriter.NewCoreWriter(conf.Writer, configure.WithParent("log"), gl.lv)
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
		gl.AutoLevel(configure.WithParent("log"))
	}
	return gl

}

// AutoLevel ...
func (l *GalaxyLog) AutoLevel(c *gconfig.Configuration) {
	c.Watch(func(c *gconfig.Configuration) {
		lvText := strings.ToLower(c.GetString("level"))
		if lvText != "" {
			l.Info("update level", String("level", lvText), String("name", l.name))
			_ = l.lv.UnmarshalText([]byte(lvText))
		}
	})
}

// SetLevel ...
func (l *GalaxyLog) SetLevel(lv Level) {
	l.lv.SetLevel(lv)
}

// Flush ...
// When use os.Stdout or os.Stderr as zapcore.WriteSyncer
// l.desugar.Sync() maybe return an error like this: 'sync /dev/stdout: The handle is invalid.'
// Because os.Stdout and os.Stderr is a non-normal file, maybe not support 'fsync' in different os platform
// So ignored Sync() return value
// About issues: https://github.com/uber-go/zap/issues/328
// About 'fsync': https://man7.org/linux/man-pages/man2/fsync.2.html
func (l *GalaxyLog) Flush() error {
	if l.asyncStopFunc != nil {
		if err := l.asyncStopFunc(); err != nil {
			return err
		}
	}

	_ = l.desugar.Sync()
	return nil
}

// IsDebugMode ...
func (l *GalaxyLog) IsDebugMode() bool {
	return true
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-32s", msg)
}

// Debug ...
func (l *GalaxyLog) Debug(msg string, fields ...Field) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.desugar.Debug(msg, fields...)
}

// Debugw ...
func (l *GalaxyLog) Debugw(msg string, keysAndValues ...interface{}) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.sugar.Debugw(msg, keysAndValues...)
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
func (l *GalaxyLog) StdLog() *log.Logger {
	return zap.NewStdLog(l.desugar)
}

// Debugf ...
func (l *GalaxyLog) Debugf(template string, args ...interface{}) {
	l.sugar.Debugw(sprintf(template, args...))
}

// Info ...
func (l *GalaxyLog) Info(msg string, fields ...Field) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.desugar.Info(msg, fields...)
}

// Infow ...
func (l *GalaxyLog) Infow(msg string, keysAndValues ...interface{}) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.sugar.Infow(msg, keysAndValues...)
}

// Infof ...
func (l *GalaxyLog) Infof(template string, args ...interface{}) {
	l.sugar.Infof(sprintf(template, args...))
}

// Warn ...
func (l *GalaxyLog) Warn(msg string, fields ...Field) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.desugar.Warn(msg, fields...)
}

// Warnw ...
func (l *GalaxyLog) Warnw(msg string, keysAndValues ...interface{}) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.sugar.Warnw(msg, keysAndValues...)
}

// Warnf ...
func (l *GalaxyLog) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(sprintf(template, args...))
}

// Error ...
func (l *GalaxyLog) Error(msg string, fields ...Field) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.desugar.Error(msg, fields...)
}

// Errorw ...
func (l *GalaxyLog) Errorw(msg string, keysAndValues ...interface{}) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.sugar.Errorw(msg, keysAndValues...)
}

// Errorf ...
func (l *GalaxyLog) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(sprintf(template, args...))
}

// Panic ...
func (l *GalaxyLog) Panic(msg string, fields ...Field) {
	panicDetail(msg, fields...)
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.desugar.Panic(msg, fields...)
}

// Panicw ...
func (l *GalaxyLog) Panicw(msg string, keysAndValues ...interface{}) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.sugar.Panicw(msg, keysAndValues...)
}

// Panicf ...
func (l *GalaxyLog) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(sprintf(template, args...))
}

// DPanic ...
func (l *GalaxyLog) DPanic(msg string, fields ...Field) {
	if l.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
	}
	l.desugar.DPanic(msg, fields...)
}

// DPanicw ...
func (l *GalaxyLog) DPanicw(msg string, keysAndValues ...interface{}) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.sugar.DPanicw(msg, keysAndValues...)
}

// DPanicf ...
func (l *GalaxyLog) DPanicf(template string, args ...interface{}) {
	l.sugar.DPanicf(sprintf(template, args...))
}

// Fatal ...
func (l *GalaxyLog) Fatal(msg string, fields ...Field) {
	if l.IsDebugMode() {
		panicDetail(msg, fields...)
		//msg = normalizeMessage(msg)
		return
	}
	l.desugar.Fatal(msg, fields...)
}

// Fatalw ...
func (l *GalaxyLog) Fatalw(msg string, keysAndValues ...interface{}) {
	if l.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	l.sugar.Fatalw(msg, keysAndValues...)
}

// Fatalf ...
func (l *GalaxyLog) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(sprintf(template, args...))
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
func (l *GalaxyLog) With(fields ...Field) *GalaxyLog {
	desugarLogger := l.desugar.With(fields...)
	return &GalaxyLog{
		desugar: desugarLogger,
		lv:      l.lv,
		sugar:   desugarLogger.Sugar(),
		conf:    l.conf,
	}
}

// WithCallerSkip ...
func (l *GalaxyLog) WithCallerSkip(callerSkip int, fields ...Field) *GalaxyLog {
	l.conf.CallerSkip = callerSkip
	desugarLogger := l.desugar.WithOptions(zap.AddCallerSkip(callerSkip)).With(fields...)
	return &GalaxyLog{
		desugar: desugarLogger,
		lv:      l.lv,
		sugar:   desugarLogger.Sugar(),
		conf:    l.conf,
	}
}
