package path

const (
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
