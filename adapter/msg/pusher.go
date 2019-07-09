package msg

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/rdegges/go-ipify"

	"github.com/privatix/dappctrl/util/log"
)

const (
	caDataParameter        = "caData"
	defaultIP              = "127.0.0.1"
	serverAddressParameter = "externalIP"

	// PushedFile the name of a file that indicates that
	// the configuration is already loaded on the server.
	PushedFile = "configPushed"
	filePerm   = 0644
)

// Config is configuration to Pusher.
type Config struct {
	CaCertPath       string
	ConfigPath       string
	ExportConfigKeys []string
	TimeOut          int64
}

// SetProductConfigFunc sets controller's product configuration.
type SetProductConfigFunc func(config map[string]string) error

// Pusher updates the product configuration.
type Pusher struct {
	config           *Config
	setProductConfig SetProductConfigFunc
	ip               string
	password         string
	username         string
	logger           log.Logger
}

// NewConfig creates a default configuration.
func NewConfig() *Config {
	return &Config{
		ExportConfigKeys: []string{"proto", "cipher", "ping-restart",
			"ping", "connect-retry", "ca",
			"compress", "keepalive", "port", "server"},
		TimeOut: 12,
	}
}

// NewPusher creates a new Pusher object.
// Argument conf to parsing vpn configuration. Arguments srv, user, pass
// to send configuration to session service.
func NewPusher(conf *Config, logger log.Logger,
	setProductConfig SetProductConfigFunc) *Pusher {
	var ip string
	ip, err := externalIP()
	if err != nil {
		logger.Warn("couldn't get my IP address")
		ip = defaultIP
	}

	return &Pusher{
		config:           conf,
		setProductConfig: setProductConfig,
		logger:           logger,
		ip:               ip,
	}
}

// VpnParams parses the OpenVpn configuration file.
func (p *Pusher) VpnParams() (map[string]string, error) {
	vpnParams, err := vpnParams(p.logger, p.config.ConfigPath,
		p.config.ExportConfigKeys)
	if err != nil {
		return nil, err
	}

	ca, err := certificateAuthority(p.logger, p.config.CaCertPath)
	if err != nil {
		return nil, err
	}

	vpnParams[serverAddressParameter] = p.ip
	vpnParams[caDataParameter] = string(ca)

	return vpnParams, err
}

// PushConfiguration sends a vpn configuration to session server.
func (p *Pusher) PushConfiguration(ctx context.Context,
	params map[string]string) error {
	logger := p.logger.Add("method", "PushConfiguration")

	for {
		select {
		case <-ctx.Done():
			return ErrContextIsDone
		default:
		}

		err := p.setProductConfig(params)
		if err != nil {
			m := "failed to push app config to dappctrl"
			logger.Add("error", err.Error()).Warn(m)
			time.Sleep(time.Second *
				time.Duration(p.config.TimeOut))
			continue
		}
		logger.Info("vpn server configuration has been" +
			" successfully sent to dappctrl")
		break
	}
	return nil
}

func externalIP() (string, error) {
	return ipify.GetIp()
}

// IsDone checks if the vpn configuration is loaded to server.
func IsDone(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, PushedFile))
	return err == nil
}

// Done makes configPushed file.
func Done(dir string) error {
	file := filepath.Join(dir, PushedFile)
	return ioutil.WriteFile(file, nil, filePerm)
}
