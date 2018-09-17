package statik

import "github.com/privatix/dappctrl/util/errors"

// Errors.
const (
	// CRC16("github.com/privatix/dapp-openvpn/statik") = 0x8300
	ErrOpenFS errors.Error = 0x8300<<8 + iota
	ErrOpenFile
	ErrReadFile
)

var errMsgs = errors.Messages{
	ErrOpenFS:   "failed to open statik filesystem",
	ErrOpenFile: "failed to open statik file",
	ErrReadFile: "failed to read statik file",
}

func init() { errors.InjectMessages(errMsgs) }
