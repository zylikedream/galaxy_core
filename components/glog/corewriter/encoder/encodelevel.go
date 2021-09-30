package encoder

import (
	"github.com/zylikedream/galaxy/components/glog/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zap默认的颜色处理不好，这儿改写掉
type galaxyEncodeLevel zapcore.LevelEncoder

func (e *galaxyEncodeLevel) UnmarshalText(text []byte) error {
	switch string(text) {
	case "capital":
		*e = zapcore.CapitalLevelEncoder
	case "capitalColor":
		*e = galaxyCapitalLevelEncoder
	case "color":
		*e = galaxyLowercaseLevelEncoder
	default:
		*e = zapcore.LowercaseLevelEncoder
	}
	return nil
}

func getLevelColor(lv zapcore.Level) color.Color {
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
	return colorize
}

func galaxyCapitalLevelEncoder(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	colorize := getLevelColor(lv)
	enc.AppendByteString([]byte(colorize(lv.CapitalString())))
}

func galaxyLowercaseLevelEncoder(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	colorize := getLevelColor(lv)
	enc.AppendByteString([]byte(colorize(lv.String())))

}
