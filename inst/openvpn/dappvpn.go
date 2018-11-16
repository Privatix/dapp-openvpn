package openvpn

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/privatix/dapp-openvpn/inst/openvpn/path"
	"github.com/takama/daemon"
)

// DappVPN has a dappvpn configuration.
type DappVPN struct {
	Service string
}

// NewDappVPN creates a default dappVPN configuration.
func NewDappVPN() *DappVPN {
	return &DappVPN{}
}

// Configurate configurates dappvpn config files.
func (d *DappVPN) Configurate(o *OpenVPN) error {
	p := o.Path
	configFile := filepath.Join(p, path.DappVPNConfig)

	read, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer read.Close()

	jsonMap := make(map[string]interface{})
	json.NewDecoder(read).Decode(&jsonMap)

	maps := make(map[string]interface{})

	maps["FileLog.Filename"] = filepath.Join(p, "log/dappvpn-%Y-%m-%d.log")
	maps["OpenVPN.Name"] = filepath.Join(p, path.OpenVPN)
	maps["OpenVPN.ConfigRoot"] = filepath.Join(p, "config")
	if o.IsWindows {
		maps["OpenVPN.TapInterface"] = o.Tap.Interface
	}
	maps["Pusher.CaCertPath"] = filepath.Join(p, path.CACertificate)
	maps["Pusher.ConfigPath"] = filepath.Join(p, path.RoleConfig(o.Role))

	addr := fmt.Sprintf("%s:%v", o.Managment.IP, o.Managment.Port)
	maps["Monitor.Addr"] = addr
	if addr, err = connectorAddr(path.DappCtrlConfig); err != nil {
		return err
	}
	maps["Connector.Addr"] = addr

	if err := setConfigurationValues(jsonMap, maps); err != nil {
		return err
	}

	write, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer write.Close()

	return json.NewEncoder(write).Encode(jsonMap)
}

// InstallService installs a dappvpn service.
func (d *DappVPN) InstallService(role, dir string) (string, error) {
	d.Service = serviceName(path.DVPN, dir)
	descr := fmt.Sprintf("Privatix %s dappvpn %s", role, hash(dir))

	if strings.EqualFold(runtime.GOOS, "windows") {
		d.Service = fmt.Sprintf("Privatix DappVPN %s", hash(dir))
	}

	service, err := daemon.New(d.Service, descr)
	if err != nil {
		return "", err
	}

	return service.Install("run-dappvpn", "-workdir", dir)
}

// StartService starts dappvpn service.
func (d *DappVPN) StartService() (string, error) {
	service, err := daemon.New(d.Service, "")
	if err != nil {
		return "", err
	}
	return service.Start()
}

// RunService executes dappvpn service.
func (d *DappVPN) RunService(role, dir string) (string, error) {
	service, err := daemon.New(d.Service, "")
	if err != nil {
		return "", err
	}

	return service.Run(&execute{Path: dir, Role: role, Type: path.DVPN})
}

// StopService stops dappvpn service.
func (d *DappVPN) StopService() (string, error) {
	service, err := daemon.New(d.Service, "")
	if err != nil {
		return "", err
	}

	status, err := service.Status()
	if err != nil {
		return "", err
	}

	if !strings.Contains(strings.ToLower(status), "running") {
		return "", nil
	}

	return service.Stop()
}

// RemoveService removes the dappvpn service.
func (d *DappVPN) RemoveService() (string, error) {
	service, err := daemon.New(d.Service, "")
	if err != nil {
		return "", err
	}
	return service.Remove()
}
