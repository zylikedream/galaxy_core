package gconfig

import (
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Configuration struct {
	vp        *viper.Viper
	watch     int32
	onChanges []OnChangeCallback
	parent    string
	options   []viper.Option
}

type OnChangeCallback = func(c *Configuration)

func New(ConfigFile string, opts ...Option) *Configuration {
	conf := &Configuration{}
	for _, opt := range opts {
		opt(conf)
	}
	v := viper.NewWithOptions(conf.options...)
	v.SetConfigFile(ConfigFile)
	if err := v.ReadInConfig(); err != nil {
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

func (c *Configuration) UnmarshalKey(key string, data interface{}) error {
	return c.vp.UnmarshalKey(key, data)
}

func (c *Configuration) KeyWithParent(key string) string {
	if c.parent == "" {
		return key
	}
	return c.parent + "." + key
}

func (c *Configuration) UnmarshalKeyWithParent(key string, data interface{}) error {
	return c.vp.UnmarshalKey(c.KeyWithParent(key), data)
}
