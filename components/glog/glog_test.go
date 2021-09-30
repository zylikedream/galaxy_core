package glog

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestGlog(t *testing.T) {
	logger := NewLogger("test", "config/config.toml")
	for {
		logger.Debug("pvp", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
		logger.Warn("pvp", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
		logger.Info("pvp", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
		logger.Error("pvp", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
		logger.Infof("pvp, pid=%s", "234")
		time.Sleep(2 * time.Second)
	}
}

func TestT(t *testing.T) {
	var lv interface{} = zap.NewAtomicLevel()
	atomic, ok := lv.(zap.AtomicLevel)
	if !ok {
		t.Errorf("not type")
		return
	}
	fmt.Println(atomic)
}
