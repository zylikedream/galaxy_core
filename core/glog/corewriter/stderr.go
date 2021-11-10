package corewriter

import (
	"fmt"
	"io"
	"os"

	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/glog/corewriter/encoder"
	"github.com/zylikedream/galaxy/core/gregister"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type stderrWriter struct {
	zapcore.Core
	io.Closer
}

type stderrConfig struct {
	EncoderType string `toml:"encoder_type"`
}

// Load constructs a zapcore.Core with stderr syncer
func new(c *gconfig.Configuration, atomiclv zap.AtomicLevel) (*stderrWriter, error) {
	// Debug output to console and file by default
	w := &stderrWriter{}
	conf := &stderrConfig{
		EncoderType: "json",
	}
	if err := c.UnmarshalKeyWithParent(w.Type(), conf); err != nil {
		return nil, err
	}
	encoder, err := encoder.NewZapEncoder(conf.EncoderType, c)
	if err != nil {
		return nil, err
	}
	w.Core = zapcore.NewCore(encoder, os.Stderr, atomiclv)
	w.Closer = CloseFunc(noopCloseFunc)
	return w, nil
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
	return new(c, atomiclv)

}

func init() {
	gregister.Register((*stderrWriter)(nil))
}
