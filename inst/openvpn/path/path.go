package path

import (
	"strings"
)

const (
	// OVPN is a openvpn
	OVPN = "openvpn"
	// DVPN is a dappvpn
	DVPN = "dappvpn"
	// OpenVPN file location
	OpenVPN = `bin/openvpn/openvpn`
	// OpenSSL file location
	OpenSSL = `bin/openvpn/openssl`
	// OemVista driver file location
	OemVista = `bin/openvpn/driver/OemVista.inf`
	// TapInstall file location
	TapInstall = `bin/openvpn/tapinstall`
	// DHParam file location
	DHParam = `config/dh2048.pem`
	// CACertificate file location
	CACertificate = `config/ca.crt`
	// CAKey file location
	CAKey = `config/ca.key`
	// ServerConfigTemplate file location
	ServerConfigTemplate = `/ovpn/templates/server-config.tpl`
	// DappVPN file location
	DappVPN = `bin/dappvpn`
	// DappVPNConfig file location
	DappVPNConfig = `config/dappvpn.config.json`
	// DappCtrlConfig file location
	DappCtrlConfig = `../../dappctrl/dappctrl.config.json`
)

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
	if strings.EqualFold(t, DVPN) {
		return DappVPN
	}
	return OpenVPN
}

// VPNConfig returns vpn config path.
func VPNConfig(t, role string) string {
	if strings.EqualFold(t, DVPN) {
		return DappVPNConfig
	}
	return RoleConfig(role)
}
