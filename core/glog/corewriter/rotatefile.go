package corewriter

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zylikedream/galaxy/core/gconfig"
	"github.com/zylikedream/galaxy/core/glog/corewriter/encoder"
	"github.com/zylikedream/galaxy/core/glog/corewriter/rotate"
)

type rotateFileWriter struct {
	zapcore.Core
	io.Closer
}

// config ...
type rotateFileConfig struct {
	Dir                 string        `toml:"dir"`
	File                string        `toml:"file"`
	MaxSize             int           `toml:"max_size"`              // [fileWriter]日志输出文件最大长度，超过改值则截断，默认500M
	MaxAge              int           `toml:"max_age"`               // [fileWriter]日志存储最大时间，默认最大保存天数为7天
	MaxBackup           int           `toml:"max_backup"`            // [fileWriter]日志存储最大数量，默认最大保存文件个数为10个
	RotateInterval      time.Duration `toml:"rotate_interval"`       // [fileWriter]日志轮转时间，默认1天
	FlushBufferSize     int           `toml:"flush_buffer_size"`     // 缓冲大小，默认256 * 1024B
	FlushBufferInterval time.Duration `toml:"flush_buffer_interval"` // 缓冲时间，默认5秒
	EnableAsync         bool          `toml:"enable_async"`          // 是否异步，默认异步
	EncoderType         string        `toml:"encoder_type"`          // console|json 使用可读或者json格式
	Stdout              bool          `toml:"stdout"`                // 是否同时输出到控制台
}

func defaultConfig() *rotateFileConfig {
	return &rotateFileConfig{
		Dir:                 "log",
		File:                "default.log",
		MaxSize:             500, // 500M
		MaxAge:              7,   // 1 week
		MaxBackup:           10,  // 10 backup
		RotateInterval:      24 * time.Hour,
		FlushBufferSize:     256 * 1024,
		FlushBufferInterval: 5 * time.Second,
		EncoderType:         "json",
	}
}

// Load constructs a zapcore.Core with stderr syncer
func newRotateFileWriter(c *gconfig.Configuration, atomiclv zap.AtomicLevel) *rotateFileWriter {
	w := &rotateFileWriter{}
	conf := defaultConfig()
	if err := c.UnmarshalKeyWithParent(w.Type(), &conf); err != nil {
		panic(err)
	}
	// NewRotateFileCore constructs a zapcore.Core with rotate file syncer
	// Debug output to console and file by default
	cf := noopCloseFunc
	var ws = zapcore.AddSync(&rotate.RLogger{
		Filename:   path.Join(conf.Dir, conf.File),
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackup,
		LocalTime:  true,
		Compress:   false,
		Interval:   conf.RotateInterval,
	})

	if conf.EnableAsync {
		ws, cf = rotate.BufferWriteSyncer(ws, conf.FlushBufferSize, conf.FlushBufferInterval)
	}
	if conf.Stdout {
		ws = zap.CombineWriteSyncers(os.Stdout, ws)
	}
	zapEncoder, err := encoder.NewZapEncoder(conf.EncoderType, c)
	if err != nil {
		panic(err)
	}
	w.Closer = CloseFunc(cf)
	w.Core = zapcore.NewCore(zapEncoder, ws, atomiclv)
	return w
}

func (r *rotateFileWriter) Type() string {
	return "rotate_file"
}

func (r *rotateFileWriter) Build(c *gconfig.Configuration, args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("params num error")
	}
	atomiclv, ok := args[0].(zap.AtomicLevel)
	if !ok {
		return nil, fmt.Errorf("need param type (zap.AtomicLevel)")
	}
	return newRotateFileWriter(c, atomiclv), nil
}

func init() {
	Register(&rotateFileWriter{})
}
