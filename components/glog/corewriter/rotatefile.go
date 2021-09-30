package corewriter

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zylikedream/galaxy/components/gconfig"
<<<<<<< HEAD
=======
	"github.com/zylikedream/galaxy/components/glog/corewriter/encoder"
>>>>>>> 94b59071e8de9f193171c73707897a0a8492681f
	"github.com/zylikedream/galaxy/components/glog/corewriter/rotate"
)

type rotateFileWriter struct {
	zapcore.Core
	io.Closer
}

// config ...
type config struct {
	Dir                 string        `mapstructure:"dir"`
	File                string        `mapstructure:"file"`
	MaxSize             int           `mapstructure:"max_size"`              // [fileWriter]日志输出文件最大长度，超过改值则截断，默认500M
	MaxAge              int           `mapstructure:"max_age"`               // [fileWriter]日志存储最大时间，默认最大保存天数为7天
	MaxBackup           int           `mapstructure:"max_backup"`            // [fileWriter]日志存储最大数量，默认最大保存文件个数为10个
	RotateInterval      time.Duration `mapstructure:"rotate_interval"`       // [fileWriter]日志轮转时间，默认1天
	FlushBufferSize     int           `mapstructure:"flush_buffer_size"`     // 缓冲大小，默认256 * 1024B
	FlushBufferInterval time.Duration `mpastructure:"flush_buffer_interval"` // 缓冲时间，默认5秒
	EnableAsync         bool          `mapstructure:"enable_sync"`           // 是否异步，默认异步
	Encoder             string        `mapstructure:"encoder"`               // console|json 使用可读或者json格式
	Stdout              bool          `mapstructure:"stdout"`                // 是否同时输出到控制台
}

func defaultConfig() *config {
	return &config{
		Dir:                 "log",
		File:                "default.log",
		MaxSize:             500, // 500M
		MaxAge:              7,   // 1 week
		MaxBackup:           10,  // 10 backup
		RotateInterval:      24 * time.Hour,
		FlushBufferSize:     256 * 1024,
		FlushBufferInterval: 5 * time.Second,
	}
}

// Load constructs a zapcore.Core with stderr syncer
func newRotateFileWriter(c *gconfig.Configuration, atomiclv zap.AtomicLevel) CoreWriter {
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

	var encoder zapcore.Encoder
	if conf.EnableAsync {
		ws, cf = rotate.BufferWriteSyncer(ws, conf.FlushBufferSize, conf.FlushBufferInterval)
	}
	if conf.Stdout {
		ws = zap.CombineWriteSyncers(os.Stdout, ws)
	}
	if conf.Encoder == "console" {
		encoder = zapcore.NewConsoleEncoder(*encoder.DefaultDebugConfig())
	} else if conf.Encoder == "json" {
		encoder = zapcore.NewJSONEncoder(*encoder.DefaultZapConfig())
	} else {
		panic(fmt.Errorf("unkonw encoder %s", conf.Encoder))
	}
	w.Closer = CloseFunc(cf)
	w.Core = zapcore.NewCore(encoder, ws, atomiclv)
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
		return nil, fmt.Errorf("need param type (*zap.AtomicLevel)")
	}
	return newRotateFileWriter(c, atomiclv), nil
}
