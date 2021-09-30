package encoder

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap/zapcore"
)

type EncoderHook interface {
	UnmarshalText(text []byte) error
}

func stringToEncoder(hook EncoderHook) mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		hookType := reflect.TypeOf(hook).Elem()
		if t != hookType {
			return data, nil
		}

		// Convert it by parsing
		v := reflect.New(hookType).Interface()
		if err := v.(EncoderHook).UnmarshalText([]byte(data.(string))); err != nil {
			return nil, err
		}
		return v, nil
	}
}

func stringToLevelEncoder() mapstructure.DecodeHookFunc {
	return stringToEncoder((*galaxyEncodeLevel)(nil))
}

func stringToTimeEncoder() mapstructure.DecodeHookFunc {
	return stringToEncoder((*zapcore.TimeEncoder)(nil))
}

func stringToDurationEncoder() mapstructure.DecodeHookFunc {
	return stringToEncoder((*zapcore.DurationEncoder)(nil))
}

func stringToCallerEncoder() mapstructure.DecodeHookFunc {
	return stringToEncoder((*zapcore.CallerEncoder)(nil))
}

func stringToNameEncoder() mapstructure.DecodeHookFunc {
	return stringToEncoder((*zapcore.NameEncoder)(nil))
}
