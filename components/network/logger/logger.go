package logger

import "github.com/zylikedream/galaxy/components/glog"

var Nlog *glog.GalaxyLog

func SetLogger(l *glog.GalaxyLog) {
	Nlog = l
}
