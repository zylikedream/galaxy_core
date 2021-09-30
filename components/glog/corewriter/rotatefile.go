package corewriter

import (
	"io"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/glog/corewriter/rotate"
)

type rotateFileWriter struct {
	zapcore.Core
	io.Closer
}

// config ...
type config struct {
	Dir                 string        `mapstructure:"dir"`
	name                string        `mapstructure:"name"`
	MaxSize             int           `mapstructure:"max_size"`              // [fileWriter]日志输出文件最大长度，超过改值则截断，默认500M
	MaxAge              int           `mapstructure:"max_age"`               // [fileWriter]日志存储最大时间，默认最大保存天数为7天
	MaxBackup           int           `mapstructure:"max_backup"`            // [fileWriter]日志存储最大数量，默认最大保存文件个数为10个
	RotateInterval      time.Duration `mapstructure:"rotate_interval"`       // [fileWriter]日志轮转时间，默认1天
	FlushBufferSize     int           `mapstructure:"flush_buffer_size"`     // 缓冲大小，默认256 * 1024B
	FlushBufferInterval time.Duration `mpastructure:"flush_buffer_interval"` // 缓冲时间，默认5秒
	encoderConfig       *zapcore.EncoderConfig
}

func defaultConfig() *config {
	return &config{
		MaxSize:             500, // 500M
		MaxAge:              7,   // 1 day
		MaxBackup:           10,  // 10 backup
		RotateInterval:      24 * time.Hour,
		FlushBufferSize:     256 * 1024,
		FlushBufferInterval: 5 * time.Second,
	}
}

// Load constructs a zapcore.Core with stderr syncer
func newRotateFileWriter(c *gconfig.Configuration) CoreWriter {
	w := &rotateFileWriter{}
	conf := defaultConfig()
	if err := c.UnmarshalKey(w.Type(), &conf); err != nil {
		panic(err)
	}
	// NewRotateFileCore constructs a zapcore.Core with rotate file syncer
	// Debug output to console and file by default
	cf := noopCloseFunc
	var ws = zapcore.AddSync(&rotate.RLogger{
		Filename:   path.Join(conf.Dir, conf.name),
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackup,
		LocalTime:  true,
		Compress:   false,
		Interval:   conf.RotateInterval,
	})

	debug := true
	if debug {
		ws = zap.CombineWriteSyncers(os.Stdout, ws)
		ws, cf = rotate.BufferWriteSyncer(ws, conf.FlushBufferSize, conf.FlushBufferInterval)
	}
	w.Closer = CloseFunc(cf)
	w.Core = zapcore.NewCore(
		func() zapcore.Encoder {
			if debug {
				return zapcore.NewConsoleEncoder(*conf.encoderConfig)
			}
			return zapcore.NewJSONEncoder(*conf.encoderConfig)
		}(),
		ws,
		Conf.AtomicLevel(),
	)
	return w
}

func (r *rotateFileWriter) Type() string {
	return "rotate_file"
}
