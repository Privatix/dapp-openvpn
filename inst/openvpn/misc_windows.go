package openvpn

import "fmt"

func serviceName(prefix, path string) string {
	return fmt.Sprintf("%s_%s", prefix, hash(path))
}

func networkInterface() (string, error) {
	return "", nil
}

func openvpnInterface() string {
	return ""
}
