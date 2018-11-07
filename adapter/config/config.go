package config

import (
	"time"

	"github.com/privatix/dappctrl/svc/connector"
	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/adapter/mon"
	"github.com/privatix/dapp-openvpn/adapter/msg"
	"github.com/privatix/dapp-openvpn/adapter/tc"
)

type ovpnConfig struct {
	Name       string        // Name of OvenVPN executable.
	Args       []string      // Extra arguments for OpenVPN executable.
	ConfigRoot string        // Root path for OpenVPN channel configs.
	StartDelay time.Duration // Delay to ensure OpenVPN is ready, in milliseconds.
}

// Config is dapp-openvpn adapter configuration.
type Config struct {
	ChannelDir      string // Directory for common-name -> channel mappings.
	ClientMode      bool
	HeartbeatPeriod time.Duration // In milliseconds.
	FileLog         *log.FileConfig
	Monitor         *mon.Config
	OpenVPN         *ovpnConfig // OpenVPN settings for client mode.
	Pusher          *msg.Config
	Connector       *connector.Config
	TC              *tc.Config
}

// NewConfig creates default dapp-openvpn configuration.
func NewConfig() *Config {
	return &Config{
		ChannelDir: ".",
		FileLog:    log.NewFileConfig(),
		Monitor:    mon.NewConfig(),
		OpenVPN: &ovpnConfig{
			Name:       "openvpn",
			ConfigRoot: "/etc/openvpn/config",
			StartDelay: 1000,
		},
		Pusher:    msg.NewConfig(),
		Connector: connector.DefaultConfig(),
		TC:        tc.NewConfig(),
	}
}
