package gxylog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
