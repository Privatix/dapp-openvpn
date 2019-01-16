package openvpn

import (
	"fmt"
)

func serviceName(prefix, path string) string {
	return fmt.Sprintf("%s_%s", prefix, hash(path))
}

func setRegValue(guid, name string) error {
	return nil
}
