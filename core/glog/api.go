package glog

import "bytes"

// PackageName 包名
// defaultLogger defines default logger, it's usually used in application business logic
var defaultLogger *GalaxyLog

var defaultConfig = []byte(`
[log]
	level = "debug"
	enable_addcaller = true
	writer = "stderr" # 写入类型(rotate_file|stderr), 默认为rotate_file
	watch = false
	caller_skip = 2

[writer.stderr]
	encoder_type = "console"

[writer.encoder_config] # 完整字段参考github.com/zylikedream/galaxy/core/glog/corewriter/encoder/encoder.go的config
	encode_level = "capital" # level的颜色大小写控制(capital|capitalColor|color|lower), 默认为lower
	encode_time = "rfc3339" # 时间格式(rfc3339|iso8601|mills|nanos|sec)，默认为rfc3339
	encode_duration = "sec" # duration写入格式(string|nanas|ms|sec), 默认为sec
	encode_caller = "short" # caller的格式（full|short), 默认为short

`)

func init() {
	defaultLogger = NewLoggerWithReader("default", bytes.NewBuffer(defaultConfig))
}

func SetDefaultLogger(l *GalaxyLog) {
	defaultLogger = l
}

func DefaultLogger() *GalaxyLog {
	return defaultLogger
}

// Info ...
func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

// Debug ...
func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

// Warn ...
func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

// Error ...
func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}

// Panic ...
func Panic(msg string, fields ...Field) {
	defaultLogger.Panic(msg, fields...)
}

// DPanic ...
func DPanic(msg string, fields ...Field) {
	defaultLogger.DPanic(msg, fields...)
}

// Fatal ...
func Fatal(msg string, fields ...Field) {
	defaultLogger.Fatal(msg, fields...)
}

// Debugw ...
func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

// Infow ...
func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

// Warnw ...
func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

// Panicw ...
func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

// DPanicw ...
func DPanicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.DPanicw(msg, keysAndValues...)
}

// Fatalw ...
func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

// Debugf ...
func Debugf(msg string, args ...interface{}) {
	defaultLogger.Debugf(msg, args...)
}

// Infof ...
func Infof(msg string, args ...interface{}) {
	defaultLogger.Infof(msg, args...)
}

// Warnf ...
func Warnf(msg string, args ...interface{}) {
	defaultLogger.Warnf(msg, args...)
}

// Errorf ...
func Errorf(msg string, args ...interface{}) {
	defaultLogger.Errorf(msg, args...)
}

// Panicf ...
func Panicf(msg string, args ...interface{}) {
	defaultLogger.Panicf(msg, args...)
}

// DPanicf ...
func DPanicf(msg string, args ...interface{}) {
	defaultLogger.DPanicf(msg, args...)
}

// Fatalf ...
func Fatalf(msg string, args ...interface{}) {
	defaultLogger.Fatalf(msg, args...)
}

// With ...
func With(fields ...Field) *GalaxyLog {
	return defaultLogger.With(fields...)
}
