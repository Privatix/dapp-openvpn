package openvpn

import (
	"fmt"
	"net"
	"os/exec"
)

func serviceName(prefix, path string) string {
	return fmt.Sprintf("io.privatix.%s_%s", prefix, hash(path))
}

func networkInterface() (string, error) {
	routeCmd := exec.Command("/bin/sh", "-c",
		"route get default| grep interface| awk '{print $2}'")
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func openvpnInterface() string {
	prefix := "utun"
	index := 0
	for {
		name := fmt.Sprintf("%s%v", prefix, index)
		if _, err := net.InterfaceByName(name); err != nil {
			return name
		}
		index++
	}
}
