package corewriter

import (
	"io"

	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/gregister"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WriterBuilder 根据key初始化writer
// Writer 日志interface
type CoreWriter interface {
	zapcore.Core
	io.Closer
}

type CloseFunc func() error

// Close 关闭
func (c CloseFunc) Close() error {
	return c()
}

var noopCloseFunc = func() error {
	return nil
}

var reg = gregister.NewRegister()

func Register(builder gregister.Builder) {
	reg.Register(builder)
}

func NewCoreWriter(t string, c *gconfig.Configuration, atomiclv zap.AtomicLevel) (CoreWriter, error) {
	if node, err := reg.NewNode(t, c, atomiclv); err != nil {
		return nil, err
	} else {
		return node.(CoreWriter), nil
	}
}
