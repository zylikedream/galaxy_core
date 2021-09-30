package encoder

import (
	"time"

	"github.com/zylikedream/galaxy/components/gconfig"
	"github.com/zylikedream/galaxy/components/glog/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

func timeDebugEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// debugEncodeLevel ...
func debugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = color.Red
	switch lv {
	case zapcore.DebugLevel:
		colorize = color.Blue
	case zapcore.InfoLevel:
		colorize = color.Green
	case zapcore.WarnLevel:
		colorize = color.Yellow
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = color.Red
	default:
	}
	enc.AppendString(colorize(lv.String()))
}

func DefaultZapConfig() *zapcore.EncoderConfig {
	zap.NewProductionEncoderConfig()
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func DefaultDebugConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    debugEncodeLevel,
		EncodeTime:     timeDebugEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

type config struct {
	MessageKey       string                  `toml:"message_key"`
	LevelKey         string                  `toml:"level_key"`
	TimeKey          string                  `toml:"time_key"`
	NameKey          string                  `toml:"name_key"`
	CallerKey        string                  `toml:"caller_key"`
	FunctionKey      string                  `toml:"function_key"`
	StacktraceKey    string                  `toml:"stacktrace_key"`
	LineEnding       string                  `toml:"line_ending"`
	EncodeLevel      zapcore.LevelEncoder    `toml:"encode_level"`
	EncodeTime       zapcore.TimeEncoder     `toml:"encode_time"`
	EncodeDuration   zapcore.DurationEncoder `toml:"encode_duration"`
	EncodeCaller     zapcore.CallerEncoder   `toml:"encode_caller"`
	EncodeName       zapcore.NameEncoder     `toml:"encode_name"`
	ConsoleSeparator string                  `toml:"console_separator"`
}

func newZapConfig(c *gconfig.Configuration) *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    debugEncodeLevel,
		EncodeTime:     timeDebugEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
