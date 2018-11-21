package path

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/privatix/dappctrl/util"
)

// Config is a global variable to path configuration.
var Config *config

func init() {
	Config = newConfig()
	dir := filepath.Dir(os.Args[0])
	path := filepath.Join(dir, "../config/path.config.json")
	if _, err := os.Stat(path); err == nil {
		_ = util.ReadJSONFile(path, &Config)
	}
}

// config has a path configuration.
type config struct {
	// OVPN is a openvpn
	OVPN string
	// DVPN is a dappvpn
	DVPN string
	// OpenVPN file location
	OpenVPN string
	// OpenSSL file location
	OpenSSL string
	// OemVista driver file location
	OemVista string
	// TapInstall file location
	TapInstall string
	// DHParam file location
	DHParam string
	// CACertificate file location
	CACertificate string
	// CAKey file location
	CAKey string
	// ServerConfigTemplate file location
	ServerConfigTemplate string
	// DappVPN file location
	DappVPN string
	// DappVPNConfig file location
	DappVPNConfig string
	// DappCtrlConfig file location
	DappCtrlConfig string
}

// newConfig creates a default path configuration.
func newConfig() *config {
	return &config{
		OVPN:                 "openvpn",
		DVPN:                 "dappvpn",
		OpenVPN:              `bin/openvpn/openvpn`,
		OpenSSL:              `bin/openvpn/openssl`,
		OemVista:             `bin/openvpn/driver/OemVista.inf`,
		TapInstall:           `bin/openvpn/tapinstall`,
		DHParam:              `config/dh2048.pem`,
		CACertificate:        `config/ca.crt`,
		CAKey:                `config/ca.key`,
		ServerConfigTemplate: `/ovpn/templates/server-config.tpl`,
		DappVPN:              `bin/dappvpn`,
		DappVPNConfig:        `config/dappvpn.config.json`,
		DappCtrlConfig:       `../../dappctrl/dappctrl.config.json`,
	}
}

// RoleCertificate returns roles certificate path.
func RoleCertificate(role string) string {
	return "config/" + role + ".crt"
}

// RoleKey returns roles private key path.
func RoleKey(role string) string {
	return "config/" + role + ".key"
}

// RoleConfig returns roles config path.
func RoleConfig(role string) string {
	return "config/" + role + ".conf"
}

// VPN returns vpn path.
func VPN(t string) string {
	if strings.EqualFold(t, Config.DVPN) {
		return Config.DappVPN
	}
	return Config.OpenVPN
}

// VPNConfig returns vpn config path.
func VPNConfig(t, role string) string {
	if strings.EqualFold(t, Config.DVPN) {
		return Config.DappVPNConfig
	}
	return RoleConfig(role)
}
