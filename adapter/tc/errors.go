package tc

import "github.com/privatix/dappctrl/util/errors"

// Errors.
const (
	// CRC16("github.com/privatix/dapp-openvpn/adapter/tc") = 0x63FE
	ErrBadClientIP errors.Error = 0x63FE<<8 + iota
)

var errMsgs = errors.Messages{
	ErrBadClientIP: "bad client IP",
}

func init() { errors.InjectMessages(errMsgs) }
