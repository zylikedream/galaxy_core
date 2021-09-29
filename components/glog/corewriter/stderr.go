package corewriter

import (
	"fmt"
	"io"
	"os"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/glog/corewriter/encoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type stderrWriter struct {
	zapcore.Core
	io.Closer
}

// Load constructs a zapcore.Core with stderr syncer
func new(c *gconfig.Configuration, atomiclv zap.AtomicLevel) *stderrWriter {
	// Debug output to console and file by default
	w := &stderrWriter{}
	debug := true
	var econfig zapcore.EncoderConfig
	if debug {
		econfig = *encoder.DefaultDebugConfig()
	} else {
		econfig = *encoder.DefaultZapConfig()
	}
	w.Core = zapcore.NewCore(zapcore.NewJSONEncoder(econfig), os.Stderr, atomiclv)
	w.Closer = CloseFunc(noopCloseFunc)
	return w
}

func (s *stderrWriter) Type() string {
	return "stderr"
}

func (s *stderrWriter) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("params num error")
	}
	atomiclv, ok := args[0].(zap.AtomicLevel)
	if !ok {
		return nil, fmt.Errorf("need param type (*zap.AtomicLevel)")
	}
	return new(c, atomiclv), nil

}
