package prepare

import (
	"path/filepath"

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
		msg.SpecificOptions(adapterConfig.Monitor))
	if err != nil {
		return ErrMakeConfig
	}
	return nil
}
