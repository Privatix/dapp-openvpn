package msg

import "github.com/privatix/dappctrl/util/errors"

// Errors.
const (
	// CRC16("github.com/privatix/dapp-openvpn/adapter/msg") = 0x6D7F
	ErrContextIsDone errors.Error = 0x6D7F<<8 + iota
	ErrCreateAccessFile
	ErrCreateConfig
	ErrCreateDir
	ErrDecodeParams
	ErrFindCert
	ErrGenConfig
	ErrParseConfigTemplate
	ErrReadCert
	ErrReadConfig
	ErrServiceEndpointAddr
)

var errMsgs = errors.Messages{
	ErrContextIsDone:       "context is done",
	ErrCreateAccessFile:    "failed to create access file",
	ErrCreateConfig:        "failed to create config",
	ErrCreateDir:           "failed to create directory",
	ErrDecodeParams:        "failed to decode additional params",
	ErrFindCert:            "failed to find certificate authority",
	ErrGenConfig:           "failed to generate config",
	ErrParseConfigTemplate: "failed to parse template for config",
	ErrReadCert:            "failed to read certificate authority",
	ErrReadConfig:          "failed to read config",
	ErrServiceEndpointAddr: "invalid service endpoint address",
}

func init() { errors.InjectMessages(errMsgs) }
