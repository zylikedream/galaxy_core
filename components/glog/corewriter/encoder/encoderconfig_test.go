package encoder

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type tconfig struct {
	LevelEncoder zapcore.LevelEncoder `mapstructure:"level_encoder"`
}

func StringToLeveEncoder(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	fmt.Printf("f.type = %s, t.type=%s\n", f.Name(), t.Name())
	if f.Kind() != reflect.String {
		return data, nil
	}
	var encoder zapcore.LevelEncoder
	if t != reflect.TypeOf(encoder) {
		return data, nil
	}

	// Convert it by parsing
	if err := encoder.UnmarshalText([]byte(data.(string))); err != nil {
		return nil, err
	}
	return encoder, nil
}

func TestUmarshal(t *testing.T) {
	viper.SetConfigFile("config.toml")
	viper.ReadInConfig()
	conf := &tconfig{}
	var encoder zapcore.LevelEncoder
	fmt.Printf("encoder:%v\n", encoder)
	opt := viper.DecodeHook(stringToEncoder((*zapcore.LevelEncoder)(nil)))
	// opt := viper.DecodeHook(stringToEncoder(&encoder))
	err := viper.UnmarshalKey("test", conf, opt)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("encoder=%v", reflect.ValueOf(conf.LevelEncoder).Pointer() == reflect.ValueOf(zapcore.CapitalLevelEncoder).Pointer())
	t.Logf("encoder=%v", reflect.ValueOf(conf.LevelEncoder).Pointer() == reflect.ValueOf(zapcore.LowercaseColorLevelEncoder).Pointer())
}
