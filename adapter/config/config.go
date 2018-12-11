package config

import (
	"time"

	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/adapter/mon"
	"github.com/privatix/dapp-openvpn/adapter/msg"
	"github.com/privatix/dapp-openvpn/adapter/tc"
)

type ovpnConfig struct {
	Name         string        // Name of OvenVPN executable.
	Args         []string      // Extra arguments for OpenVPN executable.
	ConfigRoot   string        // Root path for OpenVPN channel configs.
	StartDelay   time.Duration // Delay to ensure OpenVPN is ready, in milliseconds.
	TapInterface string        // Windows TAP device name.
}

type sessConfig struct {
	Endpoint string
	Origin   string
	Product  string
	Password string
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
	Sess            *sessConfig
	TC              *tc.Config
}

// NewConfig creates default dapp-openvpn configuration.
func NewConfig() *Config {
	return &Config{
		ChannelDir:      ".",
		ClientMode:      false,
		HeartbeatPeriod: 2000,
		FileLog:         log.NewFileConfig(),
		Monitor:         mon.NewConfig(),
		OpenVPN: &ovpnConfig{
			Name:       "openvpn",
			ConfigRoot: "/etc/openvpn/config",
			StartDelay: 1000,
		},
		Pusher: msg.NewConfig(),
		Sess: &sessConfig{
			Endpoint: "ws://localhost:8000/ws",
		},
		TC: tc.NewConfig(),
	}
}
