package openvpn

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"github.com/privatix/dapp-openvpn/inst/openvpn/path"
)

const serverTapNamePrefix = "Privatix VPN Server"
const clientTapNamePrefix = "Privatix VPN Client"

type tapInterface struct {
	DeviceID  string
	GUID      string
	Interface string
}

func installedTapInterfaces(tap string) ([]string, error) {
	output, err := exec.Command(tap, "status", "tap0901").CombinedOutput()

	if err != nil {
		return nil, err
	}

	return matchTAPInterface(string(output)), nil
}

func installTAP(p, role string) (*tapInterface, error) {
	driver := filepath.Join(p, path.Config.OemVista)
	tapExec := filepath.Join(p, path.Config.TapInstall)

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

	return newTAP(deviceID, role)
}

func newTAP(deviceID, role string) (*tapInterface, error) {
	id, err := strconv.ParseInt(deviceID[len(deviceID)-4:], 10, 64)
	if err != nil {
		return nil, errors.New("failed to parse device id")
	}

	tapInterfaceName := clientTapNamePrefix
	if strings.EqualFold(role, "server") {
		tapInterfaceName = serverTapNamePrefix
	}
	if id > 0 {
		tapInterfaceName = fmt.Sprintf("%s %v", tapInterfaceName, id)
	}

	fmt.Println("device", deviceID)
	fmt.Println("interface", tapInterfaceName)

	guid, err := tapInterfaceGUID(deviceID)
	if err != nil {
		return nil, err
	}

	tap := &tapInterface{
		DeviceID:  deviceID,
		GUID:      guid,
		Interface: tapInterfaceName,
	}

	return tap, renameTapInterface(deviceID, tapInterfaceName)
}

func tapInterfaceGUID(device string) (string, error) {
	output, err := exec.Command("wmic", "nic",
		"where", "PNPDeviceID='"+strings.Replace(device, `\`, `\\`, -1)+"'",
		"get", "GUID", "/value").CombinedOutput()

	if err != nil {
		return "", err
	}

	guid := strings.Replace(strings.TrimSpace(string(output)),
		"GUID=", "", -1)

	return guid, nil
}

func renameTapInterface(device, name string) error {
	output, err := exec.Command("wmic", "nic",
		"where", "PNPDeviceID='"+strings.Replace(device, `\`, `\\`, -1)+"'",
		"get", "NetConnectionID", "/value").CombinedOutput()

	if err != nil {
		return err
	}

	oldName := strings.Replace(strings.TrimSpace(string(decode(output))),
		"NetConnectionID=", "", -1)

	cmd := exec.Command("netsh", "interface", "set", "interface",
		"name="+oldName, "newname="+name)

	return cmd.Run()
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

func (tap *tapInterface) remove(p string) error {
	tapExec := filepath.Join(p, path.Config.TapInstall)
	if len(tap.DeviceID) == 0 {
		tap.DeviceID, _ = deviceID(tap.Interface)
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

func decode(b []byte) []byte {
	output, _ := exec.Command("cmd", "/C", "chcp").CombinedOutput()
	s := strings.Split(strings.TrimSpace(string(output)), ":")
	if len(s) < 2 {
		return b
	}
	codepage := strings.TrimSpace(s[1])
	var chMap *charmap.Charmap

	switch codepage {
	case "850":
		chMap = charmap.CodePage850
	case "852":
		chMap = charmap.CodePage852
	case "855":
		chMap = charmap.CodePage855
	case "858":
		chMap = charmap.CodePage858
	case "860":
		chMap = charmap.CodePage860
	case "862":
		chMap = charmap.CodePage862
	case "863":
		chMap = charmap.CodePage863
	case "865":
		chMap = charmap.CodePage865
	case "866":
		chMap = charmap.CodePage866
	default:
		chMap = charmap.CodePage437
	}

	out, _ := chMap.NewDecoder().Bytes(b)
	return out
}
