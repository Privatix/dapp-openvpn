package openvpn

import "fmt"

func serviceName(path string) string {
	return fmt.Sprintf("io.privatix.%s_%s", ovpnPrefix, hash(path))
}
