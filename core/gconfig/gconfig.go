package gconfig

import (
	"io"
	"path"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Configuration struct {
	vp         *viper.Viper
	watch      int32
	onChanges  []OnChangeCallback
	parent     string
	options    []viper.Option
	hooks      viper.DecoderConfigOption
	configType string
	tag        string
}

type OnChangeCallback = func(c *Configuration)

func defaultConfig() *Configuration {
	return &Configuration{}
}

func New(configFile string, opts ...Option) *Configuration {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}
	v := viper.NewWithOptions(conf.options...)
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	conf.configType = path.Ext(configFile)[1:]
	conf.vp = v
	return conf
}

func NewWithReader(r io.Reader, opts ...Option) *Configuration {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}
	v := viper.NewWithOptions(conf.options...)
	v.SetConfigType(conf.configType)
	if err := v.ReadConfig(r); err != nil {
		panic(err)
	}
	conf.vp = v
	return conf
}

func (c *Configuration) onConfigChange(_ fsnotify.Event) {
	for _, onChange := range c.onChanges {
		onChange(c)
	}
}

func (c *Configuration) WithParent(key string) *Configuration {
	c.parent = key
	return c
}

func (c *Configuration) Watch(onChange OnChangeCallback) {
	if atomic.LoadInt32(&c.watch) == 0 {
		atomic.StoreInt32(&c.watch, 1)
		c.vp.WatchConfig()
		c.vp.OnConfigChange(c.onConfigChange)
	}
	c.onChanges = append(c.onChanges, onChange)
}

func (c *Configuration) GetString(key string) string {
	return c.vp.GetString(key)
}

func (c *Configuration) decodeOptions() []viper.DecoderConfigOption {
	opts := []viper.DecoderConfigOption{}
	if c.hooks != nil {
		opts = append(opts, c.hooks)
	}
	opts = append(opts, func(dc *mapstructure.DecoderConfig) {
		if c.tag != "" {
			dc.TagName = c.tag
		} else {
			dc.TagName = c.configType // 默认和后缀一样
		}
	})
	return opts
}

func (c *Configuration) UnmarshalKey(key string, data interface{}) error {
	return c.vp.UnmarshalKey(key, data, c.decodeOptions()...)
}

func (c *Configuration) KeyWithParent(key string) string {
	if c.parent == "" {
		return key
	}
	return c.parent + "." + key
}

func (c *Configuration) UnmarshalKeyWithParent(key string, data interface{}) error {
	return c.UnmarshalKey(c.KeyWithParent(key), data)
}

func (c *Configuration) HookDecodeFunc(funcs ...mapstructure.DecodeHookFunc) {
	c.hooks = viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(funcs...))
}

func (c *Configuration) AllKeys() []string {
	return c.vp.AllKeys()
}
