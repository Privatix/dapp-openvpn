package prepare

import "github.com/privatix/dappctrl/util/errors"

// Errors.
const (
	// CRC16("github.com/privatix/dapp-openvpn/adapter/prepare") = 0x8C92
	ErrGetEndpoint errors.Error = 0x8C92<<8 + iota
	ErrMakeConfig
)

var errMsgs = errors.Messages{
	ErrGetEndpoint: "failed to get endpoint",
	ErrMakeConfig:  "failed to make client configuration files",
}

func init() { errors.InjectMessages(errMsgs) }
