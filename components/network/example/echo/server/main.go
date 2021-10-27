package main

import (
	"github.com/zylikedream/galaxy/components/glog"
	"github.com/zylikedream/galaxy/components/network"
	"github.com/zylikedream/galaxy/components/network/message"
	"go.uber.org/zap"
)

func main() {
	EchoServer()
}

func EchoServer() {
	p, err := network.NewNetwork("config/config.toml")
	if err != nil {
		glog.Error("network", zap.Namespace("new failed"), zap.Error(err))
		return
	}
	p.Start(func(msg *message.Message) error {
		glog.Info("network", zap.Namespace("recive package"), zap.Any("package", msg))
		return nil
	})

	done := make(chan struct{})
	<-done
}
