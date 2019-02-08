package openvpn

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const networkRegistryPath = `SYSTEM\CurrentControlSet\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}`

func serviceName(prefix, path string) string {
	return fmt.Sprintf("%s_%s", prefix, hash(path))
}

func setRegValue(guid, name string) error {
	regPath := fmt.Sprintf(`%s\%s\Connection`, networkRegistryPath, guid)
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, regPath,
		registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	return key.SetStringValue("Name", name)
}

func createNatRules(p, server string, port int) error {
	return nil
}
