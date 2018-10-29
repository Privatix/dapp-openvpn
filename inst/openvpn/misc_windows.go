package openvpn

import "fmt"

func serviceName(path string) string {
	return fmt.Sprintf("%s_%s", ovpnPrefix, hash(path))
}
