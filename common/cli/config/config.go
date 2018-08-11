package config

import (
	"flag"

	"github.com/privatix/dapp-openvpn/common/path"
)

const (
	name = "config"
)

// Config is a object that is used to process the configuration flag.
type Config struct {
	run  func(file *string, config interface{}) error
	val  *string
	conf interface{}
}

// NewConfigFlag initializes the object to process the configuration flag.
func NewConfigFlag(conf interface{}) *Config {
	return &Config{
		run: func(file *string, config interface{}) error {
			return path.ReadJSONFile(*file, &config)
		},
		val: flag.String(
			"config", "adapter.config",
			"Configuration file"),
		conf: conf,
	}

}

// Name returns flag name.
func (c *Config) Name() string {
	return name
}

// Value returns value of configuration flag.
func (c *Config) Value() interface{} {
	return c.val
}

// Process performs processing of the configuration flag.
func (c *Config) Process(flag interface{}) error {
	return c.run(c.val, c.conf)
}
