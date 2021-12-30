package gxylog

import (
	"go.uber.org/zap/zapcore"
)

// Option 可选项
type Option func(g *GalaxyLog)

func WithName(name string) Option {
	return func(g *GalaxyLog) {
		g.name = name
	}
}

// WithDebug 设置在命令行显示
// WithLevel 设置级别
func WithLevel(level string) Option {
	return func(g *GalaxyLog) {
		g.conf.Level = level
	}
}

// WithEnableAddCaller 是否添加行号，默认不添加行号
func WithEnableAddCaller(enableAddCaller bool) Option {
	return func(g *GalaxyLog) {
		g.conf.EnableAddCaller = enableAddCaller
	}
}

// WithZapCore 添加ZapCore
func WithZapCore(core zapcore.Core) Option {
	return func(g *GalaxyLog) {
		g.core = core
	}
}
