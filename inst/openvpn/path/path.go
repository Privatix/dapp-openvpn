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
	// Adapter file location
	Adapter string
	// AdapterConfig file location
	AdapterConfig string
	// DappCtrlConfig file location
	DappCtrlConfig string
	// DataDir all product related data storing location.
	DataDir string
	// VPN NAT powershell script location
	PowerShellVpnNat string
	// VPN Firewall powershell script location
	PowerShellVpnFirewall string
	// Schedule Task powershell script location
	PowerShellScheduleTask string
	// Re-enabled NAT powershell script location
	PowerShellReEnableNat string
}

// newConfig creates a default path configuration.
func newConfig() *config {
	return &config{
		OVPN:                   "openvpn",
		DVPN:                   "dappvpn",
		OpenVPN:                "bin/openvpn/openvpn",
		OpenSSL:                "bin/openvpn/openssl",
		OemVista:               "bin/openvpn/driver/OemVista.inf",
		TapInstall:             "bin/openvpn/tapinstall",
		DHParam:                "config/dh2048.pem",
		CACertificate:          "config/ca.crt",
		CAKey:                  "config/ca.key",
		ServerConfigTemplate:   "/ovpn/templates/server-config.tpl",
		Adapter:                "bin/dappvpn",
		AdapterConfig:          "config/adapter.config.json",
		DappCtrlConfig:         "../../dappctrl/dappctrl.config.json",
		DataDir:                "data",
		PowerShellVpnNat:       "bin/set-nat.ps1",
		PowerShellVpnFirewall:  "bin/set-vpnfirewall.ps1",
		PowerShellScheduleTask: "bin/new-startuptask.ps1",
		PowerShellReEnableNat:  "bin/reenable-nat.ps1",
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
		return Config.Adapter
	}
	return Config.OpenVPN
}

// VPNConfig returns vpn config path.
func VPNConfig(t, role string) string {
	if strings.EqualFold(t, Config.DVPN) {
		return Config.AdapterConfig
	}
	return RoleConfig(role)
}
