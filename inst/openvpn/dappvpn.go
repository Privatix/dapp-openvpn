package openvpn

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/takama/daemon"

	"github.com/privatix/dapp-openvpn/inst/openvpn/path"
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
	configFile := filepath.Join(p, path.Config.AdapterConfig)

	read, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer read.Close()

	jsonMap := make(map[string]interface{})
	json.NewDecoder(read).Decode(&jsonMap)

	ovpnPath := filepath.Join(p, path.Config.OpenVPN)
	upScriptPath := filepath.Join(p, path.Config.UpScript)
	downScriptPath := filepath.Join(p, path.Config.DownScript)
	if runtime.GOOS == "linux" {
		ovpnPath = "/usr/sbin/openvpn"
		upScriptPath = filepath.Join(p, "bin/update-resolv-conf.sh")
		downScriptPath = upScriptPath
	}

	maps := make(map[string]interface{})

	maps["FileLog.Filename"] = filepath.Join(p, "log/dappvpn-%Y-%m-%d.log")
	maps["OpenVPN.Name"] = ovpnPath
	maps["OpenVPN.ConfigRoot"] = filepath.Join(p, path.Config.DataDir)
	if o.IsWindows {
		maps["OpenVPN.TapInterface"] = o.Tap.GUID
	} else if o.isClient() {
		maps["OpenVPN.UpScript"] = upScriptPath
		maps["OpenVPN.DownScript"] = downScriptPath
	}
	maps["Pusher.CaCertPath"] = filepath.Join(p, path.Config.CACertificate)
	maps["Pusher.ConfigPath"] = filepath.Join(p, path.RoleConfig(o.Role))

	addr := fmt.Sprintf("%s:%v", o.Managment.IP, o.Managment.Port)
	maps["Monitor.Addr"] = addr
	addr, err = sessAddr(filepath.Join(p, path.Config.DappCtrlConfig))
	if err != nil {
		return err
	}
	maps["Sess.Endpoint"] = fmt.Sprintf("ws://%s/ws", addr)
	maps["ChannelDir"] = filepath.Join(p, path.Config.DataDir)

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
	d.Service = serviceName(path.Config.DVPN, dir)
	descr := fmt.Sprintf("Privatix %s dappvpn %s", role, hash(dir))
	var dependencies []string

	if strings.EqualFold(runtime.GOOS, "windows") {
		d.Service = fmt.Sprintf("Privatix DappVPN %s", hash(dir))
		if strings.EqualFold(role, "server") {
			dependencies = []string{
				fmt.Sprintf("Privatix_OpenVPN_%s", hash(dir))}
		}
	}

	if runtime.GOOS == "linux" {
		dependencies = []string{"dappctrl.service"}
		if role == "server" {
			dependencies = append(dependencies,
				fmt.Sprintf("openvpn_%s.service", hash(dir)))
		}
	}

	service, err := daemon.New(d.Service, descr, dependencies...)
	if err != nil {
		return "", err
	}

	return service.Install("run-adapter", "-workdir", dir)
}

// StartService starts dappvpn service.
func (d *DappVPN) StartService() (string, error) {
	service, err := daemon.New(d.Service, "")
	if err != nil {
		return "", err
	}
	s, err := service.Start()
	if err != nil && err != daemon.ErrAlreadyRunning {
		return "", err
	}
	return s, nil
}

// RunService executes dappvpn service.
func (d *DappVPN) RunService(role, dir string) (string, error) {
	service, err := daemon.New(d.Service, "")
	if err != nil {
		return "", err
	}

	return service.Run(&execute{Path: dir, Role: role,
		Type: path.Config.DVPN})
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

	s, err := service.Stop()
	if err != nil && err != daemon.ErrAlreadyStopped {
		return "", err
	}

	return s, nil
}

// RemoveService removes the dappvpn service.
func (d *DappVPN) RemoveService() (string, error) {
	service, err := daemon.New(d.Service, "")
	if err != nil {
		return "", err
	}
	return service.Remove()
}
