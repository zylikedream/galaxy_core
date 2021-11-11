package logger

import "github.com/zylikedream/galaxy/core/glog"

var Nlog *glog.GalaxyLog

func SetLogger(l *glog.GalaxyLog) {
	Nlog = l
}
