// Package cfg implements service adapter configuration methods.
package cfg

import (
	"github.com/privatix/dapp-openvpn/common/dappctrl/connector"
)

// Config is a custom installer configuration.
type Config struct {
	Connector *connector.Config
}

// DefaultConfig is a default config for installer.
func DefaultConfig() *Config {
	return &Config{
		Connector: connector.DefaultConfig(),
	}
}
