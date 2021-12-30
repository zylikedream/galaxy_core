package gxylog

import (
	"fmt"
	"testing"
	"time"

	"github.com/zylikedream/galaxy/core/gconfig"
	"go.uber.org/zap"
)

func TestGlog(t *testing.T) {
	logger := NewLogger("test", gconfig.New("config/config.example.toml"))
	for {
		logger.Debug("pvp", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
		logger.Warn("pvp", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
		logger.Info("pvp", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
		logger.Error("pvp", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
		logger.Infof("pvp, pid=%s", "234")
		time.Sleep(2 * time.Second)
	}
}

func TestDefaultLog(t *testing.T) {

	Debug("pvp1", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
	Warn("pvp1", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
	Info("pvp1", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
	Error("pvp1", zap.String("pid", "123"), zap.Duration("tm", time.Hour))
	Infof("pvp1, pid=%s", "234")
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
