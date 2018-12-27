package prepare

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/privatix/dappctrl/data"
	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/adapter/config"
	"github.com/privatix/dapp-openvpn/adapter/msg"
)

// GetEndpointFunc gets controller's channel endpoint for a give client key.
type GetEndpointFunc func(clientKey string) (*data.Endpoint, error)

// ClientConfig prepares configuration for Client. By the channel ID, finds a
// endpoint on a session server. Creates client configuration files for using a
// product.
func ClientConfig(logger log.Logger, channel string,
	adapterConfig *config.Config, getEndpoint GetEndpointFunc) error {
	logger = logger.Add("method", "ClientConfig", "channel", channel)

	endpoint, err := getEndpoint(channel)
	if err != nil {
		logger.Error(err.Error())
		return ErrGetEndpoint
	}

	save := func(str *string) string {
		if str != nil {
			return *str
		}
		return ""
	}

	target := filepath.Join(
		adapterConfig.OpenVPN.ConfigRoot, endpoint.Channel)

	err = msg.MakeFiles(logger, target,
		save(endpoint.ServiceEndpointAddress), save(endpoint.Username),
		save(endpoint.Password), endpoint.AdditionalParams,
		specificOptions(logger, adapterConfig))
	if err != nil {
		return ErrMakeConfig
	}
	return nil
}

// findTapInterface finds Windows TAP device name.
func findTapInterface(logger log.Logger,
	cfg *config.Config, options map[string]interface{}) {
	logger = logger.Add("tapInterface", cfg.OpenVPN.TapInterface)

	if cfg.OpenVPN.TapInterface == "" {
		logger.Debug("TAP interface not found")
		return
	}
	options[msg.TapInterface] = cfg.OpenVPN.TapInterface
	logger.Debug("Tap interface found")
}

// findVpnManagementPort finds OpenVpn management interface server port in
// configuration.
func findVpnManagementPort(logger log.Logger,
	cfg *config.Config, options map[string]interface{}) {
	logger = logger.Add("monitorAddress", cfg.Monitor.Addr)

	// Reads OpenVpn management interface address from configuration.
	params := strings.Split(cfg.Monitor.Addr, ":")
	if len(params) != 2 {
		logger.Debug("OpenVPN monitor address is in the wrong format")
		return
	}

	port, err := strconv.ParseUint(params[1], 10, 16)
	if err != nil {
		logger.Debug("OpenVpn management port not found")
		return
	}

	options[msg.VpnManagementPort] = uint16(port)
	logger.Debug("OpenVpn management port found")
}

// findLogDir finds directory for log files.
func findLogDir(logger log.Logger,
	cfg *config.Config, options map[string]interface{}) {
	logger = logger.Add("logFileName", cfg.FileLog.Filename)

	if cfg.FileLog.Filename == "" {
		logger.Debug("log file not set in the configuration")
		return
	}

	options[msg.LogDir] = filepath.Dir(cfg.FileLog.Filename)
	logger.Debug("directory for log files found")
}

// SpecificOptions returns specific options for dappvpn.
// These options will be used to create a product configuration.
func specificOptions(logger log.Logger,
	cfg *config.Config) map[string]interface{} {
	options := make(map[string]interface{})

	findTapInterface(logger, cfg, options)
	findVpnManagementPort(logger, cfg, options)
	findLogDir(logger, cfg, options)
	return options
}
