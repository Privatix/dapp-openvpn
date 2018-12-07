package openvpn

import (
	"fmt"
)

func serviceName(prefix, path string) string {
	return fmt.Sprintf("io.privatix.%s_%s", prefix, hash(path))
}
