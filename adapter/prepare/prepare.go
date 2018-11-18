package prepare

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/privatix/dappctrl/sesssrv"
	"github.com/privatix/dappctrl/svc/connector"
	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/adapter/config"
	"github.com/privatix/dapp-openvpn/adapter/msg"
)

// ClientConfig prepares configuration for Client.
// By the channel ID, finds a endpoint on a session server.
// Creates client configuration files for using a product.
func ClientConfig(logger log.Logger, channel string, conn connector.Connector,
	adapterConfig *config.Config) error {
	logger = logger.Add("method", "ClientConfig", "channel", channel)

	args := &sesssrv.EndpointMsgArgs{ChannelID: channel}

	endpoint, err := conn.GetEndpointMessage(args)
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
		specificOptions(adapterConfig))
	if err != nil {
		return ErrMakeConfig
	}
	return nil
}

// SpecificOptions returns specific options for dappvpn.
// These options will be used to create a product configuration.
func specificOptions(cfg *config.Config) (options map[string]interface{}) {
	options = make(map[string]interface{})

	// Reads Windows TAP device name.
	if cfg.OpenVPN.TapInterface != "" {
		options[msg.TapInterface] = cfg.OpenVPN.TapInterface
	}

	// Reads OpenVpn management interface address from configuration.
	params := strings.Split(cfg.Monitor.Addr, ":")
	if len(params) != 2 {
		return options
	}

	// Reads OpenVpn management interface server port from configuration.
	port, err := strconv.ParseUint(params[1], 10, 16)
	if err != nil {
		return options
	}

	options[msg.VpnManagementPort] = uint16(port)
	return options
}
