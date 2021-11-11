package encoder

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/zylikedream/galaxy/core/gconfig"
	"go.uber.org/zap/zapcore"
)

// debugEncodeLevel ...
type config struct {
	MessageKey       string                  `toml:"message_key"`
	LevelKey         string                  `toml:"level_key"`
	TimeKey          string                  `toml:"time_key"`
	NameKey          string                  `toml:"name_key"`
	CallerKey        string                  `toml:"caller_key"`
	FunctionKey      string                  `toml:"function_key"`
	StacktraceKey    string                  `toml:"stacktrace_key"`
	LineEnding       string                  `toml:"line_ending"`
	EncodeLevel      galaxyEncodeLevel       `toml:"encode_level"`
	EncodeTime       zapcore.TimeEncoder     `toml:"encode_time"`
	EncodeDuration   zapcore.DurationEncoder `toml:"encode_duration"`
	EncodeCaller     zapcore.CallerEncoder   `toml:"encode_caller"`
	EncodeName       zapcore.NameEncoder     `toml:"encode_name"`
	ConsoleSeparator string                  `toml:"console_separator"`
}

func defaultZapConfig() *config {
	return &config{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func newZapEncoderConfig(c *gconfig.Configuration) (*zapcore.EncoderConfig, error) {
	conf := defaultZapConfig()
	c.HookDecodeFunc(stringToCallerEncoder(), stringToDurationEncoder(), stringToLevelEncoder(), stringToTimeEncoder(), stringToNameEncoder(),
		mapstructure.StringToTimeDurationHookFunc(), mapstructure.StringToSliceHookFunc(","))
	if err := c.UnmarshalKeyWithPrefix("encoder_config", conf); err != nil {
		return nil, err
	}
	return &zapcore.EncoderConfig{
		MessageKey:       conf.MessageKey,
		LevelKey:         conf.LevelKey,
		TimeKey:          conf.TimeKey,
		NameKey:          conf.NameKey,
		CallerKey:        conf.CallerKey,
		FunctionKey:      conf.FunctionKey,
		StacktraceKey:    conf.StacktraceKey,
		LineEnding:       conf.LineEnding,
		EncodeLevel:      zapcore.LevelEncoder(conf.EncodeLevel),
		EncodeTime:       conf.EncodeTime,
		EncodeDuration:   conf.EncodeDuration,
		EncodeCaller:     conf.EncodeCaller,
		EncodeName:       conf.EncodeName,
		ConsoleSeparator: conf.ConsoleSeparator,
	}, nil
}

func NewZapEncoder(encoderType string, c *gconfig.Configuration) (zapcore.Encoder, error) {
	encoderConfig, err := newZapEncoderConfig(c)
	if err != nil {
		return nil, err
	}
	switch encoderType {
	case "json":
		return zapcore.NewJSONEncoder(*encoderConfig), nil
	case "console":
		return zapcore.NewConsoleEncoder(*encoderConfig), nil
	default:
		return nil, fmt.Errorf("unkown encoder type %s", encoderType)
	}
}
