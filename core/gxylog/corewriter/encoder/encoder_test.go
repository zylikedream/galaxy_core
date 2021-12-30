package encoder

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/zylikedream/galaxy/core/gxyconfig"
	"github.com/zylikedream/galaxy/core/gxylog/color"
	"go.uber.org/zap/zapcore"
)

type tconfig struct {
	LevelEncoder zapcore.LevelEncoder `mapstructure:"level_encoder"`
}

func TestUmarshal(t *testing.T) {
	viper.SetConfigType("toml")
	configToml := []byte(`
	[test]
		level_encoder = "capital"
	`)
	viper.ReadConfig(bytes.NewBuffer(configToml))
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

func TestHook(t *testing.T) {
	configToml := []byte(`
	[log]
		level_encoder = "capital"
	[log.encoder_config] # 完整字段参考github.com/zylikedream/galaxy/core/gxylog/corewriter/encoder/encoder.go的config
		message_key = "msg1"
		level_key = "level1"
		encode_level = "lower" # level的颜色大小写控制(capital|capitalColor|color|lower), 默认为lower
		encode_time = "rfc3339" # 时间格式(rfc3339|iso8601|mills|nanos|sec)，默认为rfc3339
		encode_duration = "string" # duration写入格式(string|nanas|ms|sec), 默认为sec
		encode_caller = "full" # caller的格式（full|short), 默认为short
		encode_name = "full" # log的格式(full), 目前只有full
	`)
	configure := gxyconfig.NewWithReader(bytes.NewBuffer(configToml), gxyconfig.WithConfigType("toml"))
	econfig, err := newZapEncoderConfig(configure)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("message_key:%s", configure.GetString("log.encoder_config.message_key"))
	t.Logf("%#v", econfig)
}

func TestColor(t *testing.T) {
	fmt.Printf("\033[1;31;40m%s\033[0m\n", "Red.")
	fmt.Println(color.Red("hello"))
}
