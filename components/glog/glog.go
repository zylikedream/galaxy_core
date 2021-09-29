package glog

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/glog/corewriter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GalaxyLog struct {
	name          string
	desugar       *zap.Logger
	lv            *zap.AtomicLevel
	conf          *config
	sugar         *zap.SugaredLogger
	asyncStopFunc func() error
}

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zap.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zap.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zap.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-Level logs.
	ErrorLevel = zap.ErrorLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zap.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zap.FatalLevel
)

type (
	// Field ...
	Field = zap.Field
	// Level ...
	Level = zapcore.Level
	// Component 组件
)

var (
	// String alias for zap.String
	String = zap.String
	// Any alias for zap.Any
	Any = zap.Any
	// Int64 alias for zap.Int64
	Int64 = zap.Int64
	// Int alias for zap.Int
	Int = zap.Int
	// Int32 alias for zap.Int32
	Int32 = zap.Int32
	// Uint alias for zap.Uint
	Uint = zap.Uint
	// Duration alias for zap.Duration
	Duration = zap.Duration
	// Durationp alias for zap.Duration
	Durationp = zap.Durationp
	// Object alias for zap.Object
	Object = zap.Object
	// Namespace alias for zap.Namespace
	Namespace = zap.Namespace
	// Reflect alias for zap.Reflect
	Reflect = zap.Reflect
	// Skip alias for zap.Skip()
	Skip = zap.Skip()
	// ByteString alias for zap.ByteString
	ByteString = zap.ByteString
)

func NewLogger(name string, configFile string) *GalaxyLog {
	configure := gconfig.New(configFile)
	conf := &config{}
	if err := configure.UnmarshalKey("log", conf); err != nil {
		panic(err)
	}
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if conf.EnableAddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(conf.CallerSkip))
	}
	if len(conf.fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(conf.fields...))
	}

	// 默认日志级别
	conf.al = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := conf.al.UnmarshalText([]byte(conf.Level)); err != nil {
		panic(err)
	}

	// 如果用户没有设置core。那么就选择官方默认的core。
	if conf.core == nil {
		w, err := corewriter.NewCoreWriter(conf.Writer, configure)
		if err != nil {
			panic(err)
		}
		conf.core = w
		conf.asyncStopFunc = w.Close
	}

	zapLogger := zap.New(conf.core, zapOptions...)
	l := &GalaxyLog{
		desugar:       zapLogger,
		lv:            &conf.al,
		conf:          conf,
		sugar:         zapLogger.Sugar(),
		name:          name,
		asyncStopFunc: conf.asyncStopFunc,
	}

	// 如果名字不为空，加载动态配置
	l.AutoLevel(configure.WithParent("log"))
	return l

}

// AutoLevel ...
func (l *GalaxyLog) AutoLevel(c *gconfig.Configuration) {
	c.Watch(func(c *gconfig.Configuration) {
		lvText := strings.ToLower(c.GetString("level"))
		if lvText != "" {
			l.Info("update level", String("level", lvText), String("name", l.conf.Name))
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

func defaultZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "l",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func defaultDebugConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "l",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    debugEncodeLevel,
		EncodeTime:     timeDebugEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// debugEncodeLevel ...
func debugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = Red
	switch lv {
	case zapcore.DebugLevel:
		colorize = Blue
	case zapcore.InfoLevel:
		colorize = Green
	case zapcore.WarnLevel:
		colorize = Yellow
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = Red
	default:
	}
	enc.AppendString(colorize(lv.CapitalString()))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

func timeDebugEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// IsDebugMode ...
func (l *GalaxyLog) IsDebugMode() bool {
	return l.conf.Debug
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
	fmt.Printf("%s: \n    %s: %s\n", Red("panic"), Red("msg"), msg)
	if _, file, line, ok := runtime.Caller(3); ok {
		fmt.Printf("    %s: %s:%d\n", Red("loc"), file, line)
	}
	for key, val := range enc.Fields {
		fmt.Printf("    %s: %s\n", Red(key), fmt.Sprintf("%+v", val))
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

// ConfigDir 获取日志路径
func (l *GalaxyLog) ConfigDir() string {
	return l.conf.Dir
}

// ConfigName 获取日志名称
func (l *GalaxyLog) ConfigName() string {
	return l.conf.Name
}
