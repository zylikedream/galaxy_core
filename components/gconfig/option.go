package gconfig

import "github.com/spf13/viper"

type Option = func(c *Configuration)

func WithKeyDelimiter(deli string) Option {
	return func(c *Configuration) {
		c.options = append(c.options, viper.KeyDelimiter(deli))
	}
}
