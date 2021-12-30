package corewriter

import (
	"io"
	"os"

	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxylog/corewriter/encoder"
	"github.com/zylikedream/galaxy/core/gxyregister"
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
func new(c *gxyconfig.Configuration, atomiclv zap.AtomicLevel) (*stderrWriter, error) {
	// Debug output to console and file by default
	w := &stderrWriter{}
	conf := &stderrConfig{
		EncoderType: "json",
	}
	if err := c.UnmarshalKey(w.Type(), conf); err != nil {
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
	return WRITER_TYPE_STDERR
}

func (s *stderrWriter) Build(c *gxyconfig.Configuration, args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, gxyregister.ErrParamNotEnough
	}
	atomiclv, ok := args[0].(zap.AtomicLevel)
	if !ok {
		return nil, gxyregister.ErrParamErrType
	}
	return new(c, atomiclv)

}

func init() {
	gxyregister.Register((*stderrWriter)(nil))
}
