package mon

import "github.com/privatix/dappctrl/util/errors"

// Errors.
const (
	// CRC16("github.com/privatix/dapp-openvpn/adapter/mon") = 0xABB7
	ErrServerOutdated errors.Error = 0xABB7 + iota
)

var errMsgs = errors.Messages{
	ErrServerOutdated: "server outdated",
}

func init() { errors.InjectMessages(errMsgs) }
