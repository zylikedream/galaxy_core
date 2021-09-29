package glog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type config struct {
	Debug           bool   // 是否双写至文件控制日志输出到终端
	Level           string // 日志初始等级，默认info级别
	EnableAddCaller bool   // 是否添加调用者信息，默认不加调用者信息
	EnableAsync     bool   // 是否异步，默认异步
	Writer          string // 使用哪种Writer，可选[file|ali|stderr]，默认file
	core            zapcore.Core
	asyncStopFunc   func() error
	fields          []zap.Field // 日志初始化字段
	CallerSkip      int
	al              zap.AtomicLevel
}
