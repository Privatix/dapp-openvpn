package openvpn

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type tapInterface struct {
	DeviceID  string
	Interface string
}

func installedTapInterfaces(tap string) ([]string, error) {
	output, err := exec.Command(tap, "status", "tap0901").CombinedOutput()

	if err != nil {
		return nil, err
	}

	return matchTAPInterface(string(output)), nil
}

func installTAP(path string) (*tapInterface, error) {
	tapInterfaceName := ovpnName(path)

	driver := filepath.Join(path, "driver/OemVista.inf")
	tapExec := filepath.Join(path, "bin/tapinstall")

	before, err := installedTapInterfaces(tapExec)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(tapExec, "install", driver, "tap0901")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	after, err := installedTapInterfaces(tapExec)
	if err != nil {
		return nil, err
	}

	if len(after)-len(before) != 1 {
		return nil, errors.New("failed to install tap driver")
	}

	deviceID := diff(before, after)
	fmt.Println("device", deviceID)
	fmt.Println("interface", tapInterfaceName)

	return newTAP(deviceID, tapInterfaceName)
}

func newTAP(deviceID, tapInterfaceName string) (*tapInterface, error) {
	if err := renameTapInterface(deviceID, tapInterfaceName); err != nil {
		return nil, err
	}

	tap := &tapInterface{
		DeviceID:  deviceID,
		Interface: tapInterfaceName,
	}

	return tap, nil
}

func renameTapInterface(device, name string) error {
	output, err := exec.Command("wmic", "nic",
		"where", "PNPDeviceID='"+strings.Replace(device, `\`, `\\`, -1)+"'",
		"get", "NetConnectionID", "/value").CombinedOutput()

	if err != nil {
		return err
	}

	oldName := strings.Replace(strings.TrimSpace(string(output)),
		"NetConnectionID=", "", -1)

	cmd := exec.Command("netsh", "interface", "set", "interface",
		"name="+oldName, "newname="+name)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func deviceID(name string) (string, error) {
	output, err := exec.Command("wmic", "nic",
		"where", "NetConnectionID='"+name+"'",
		"get", "PNPDeviceID", "/value").CombinedOutput()

	if err != nil {
		return "", err
	}

	if !strings.Contains(string(output), "PNPDeviceID=") {
		return "", nil
	}

	device := strings.Replace(strings.TrimSpace(string(output)),
		"PNPDeviceID=", "", -1)

	return device, nil
}

func (tap *tapInterface) remove(path string) error {
	tapExec := filepath.Join(path, "bin/tapinstall")
	if len(tap.DeviceID) == 0 {
		tap.DeviceID, _ = deviceID(ovpnName(path))
		if len(tap.DeviceID) == 0 {
			return errors.New("undefined tap device id")
		}
	}
	return exec.Command(tapExec, "remove", "=net", "@"+tap.DeviceID).Run()
}

func matchTAPInterface(str string) []string {
	pattern := `(?m)ROOT\\NET\\\d{4}`

	var list []string
	re := regexp.MustCompile(pattern)
	for _, match := range re.FindAllStringSubmatch(str, -1) {
		list = append(list, match[0])
	}
	return list
}
