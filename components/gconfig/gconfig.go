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
}

type OnChangeCallback = func(c *Configuration)

func New(ConfigFile string) *Configuration {
	v := viper.New()
	v.SetConfigFile(ConfigFile)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return &Configuration{
		vp: v,
	}
}

func (c *Configuration) onConfigChange(_ fsnotify.Event) {
	for _, onChange := range c.onChanges {
		onChange(c)
	}
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
