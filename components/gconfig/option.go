package gconfig

import "github.com/spf13/viper"

type Option = func(c *Configuration)

func WithKeyDelimiter(deli string) Option {
	return func(c *Configuration) {
		c.options = append(c.options, viper.KeyDelimiter(deli))
	}
}

func WithConfigType(t string) Option {
	return func(c *Configuration) {
		c.configType = t
	}
}

func WithTagName(tag string) Option {
	return func(c *Configuration) {
		c.tag = tag
	}
}
