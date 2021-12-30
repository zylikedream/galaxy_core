package logger

import "github.com/zylikedream/galaxy/core/gxylog"

var Nlog *gxylog.GalaxyLog

func SetLogger(l *gxylog.GalaxyLog) {
	Nlog = l
}
