package corewriter

import (
	"io"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/gregister"
	"go.uber.org/zap/zapcore"
)

// WriterBuilder 根据key初始化writer
// Writer 日志interface
type CoreWriter interface {
	zapcore.Core
	io.Closer
}

// Close 关闭
func (c CloseFunc) Close() error {
	return c()
}

var noopCloseFunc = func() error { return nil }

// CloseFunc should be called when the caller exits to clean up buffers.
type CloseFunc func() error

var reg = gregister.NewRegister()

func Register(builder gregister.Builder) {
	reg.Register(builder)
}

func NewCoreWriter(t string, c *gconfig.Configuration) (CoreWriter, error) {
	if node, err := reg.NewNode(t, c); err != nil {
		return nil, err
	} else {
		return node.(CoreWriter), nil
	}
}
