package openvpn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/takama/daemon"

	"github.com/privatix/dapp-openvpn/inst/openvpn/path"
	"github.com/privatix/dapp-openvpn/statik"
)

// OpenVPN has a openvpn configuration.
type OpenVPN struct {
	Path      string
	Role      string
	Tap       *tapInterface
	Proto     string
	Host      *host
	Managment *host
	Server    *host
	Service   string
	Adapter   *DappVPN
	Validity  *validity
	IsWindows bool
	User      string
	Group     string
	Import    bool
	Install   bool
}

type validity struct {
	Year  int
	Month int
	Day   int
}

type host struct {
	IP   string
	Port int
	Mask string
}

// NewOpenVPN creates a default OpenVPN configuration.
func NewOpenVPN() *OpenVPN {
	return &OpenVPN{
		Path:  ".",
		Tap:   &tapInterface{},
		Role:  "server",
		Proto: "udp",
		Host: &host{
			IP:   "0.0.0.0",
			Port: 443,
		},
		Managment: &host{
			IP:   "127.0.0.1",
			Port: 7505,
		},
		Server: &host{
			IP:   "10.217.3.0",
			Mask: "255.255.255.0",
		},
		Validity: &validity{
			Year: 10,
		},
		IsWindows: strings.EqualFold(runtime.GOOS, "windows"),
		Adapter:   NewDappVPN(),
	}
}

// InstallTap installs a new tap interface.
func (o *OpenVPN) InstallTap() (err error) {
	if o.IsWindows {
		o.Tap, err = installTAP(o.Path, o.Role)
		if err != nil {
			return err
		}

		script := filepath.Join(o.Path, path.Config.PowerShellVpnNat)
		args := buildPowerShellArgs(script,
			"-TAPdeviceAddress", o.Tap.DeviceID,
			"-Enabled")
		err = runPowerShellCommand(args...)
	}
	return err
}

// RemoveTap removes the tap interface.
func (o *OpenVPN) RemoveTap() (err error) {
	if o.IsWindows {
		script := filepath.Join(o.Path, path.Config.PowerShellVpnNat)
		args := buildPowerShellArgs(script,
			"-TAPdeviceAddress", o.Tap.DeviceID)
		if err = runPowerShellCommand(args...); err != nil {
			return err
		}

		err = o.Tap.remove(o.Path)
	}
	return err
}

// Configurate configurates openvpn config files.
func (o *OpenVPN) Configurate() error {
	if o.isClient() {
		o.Managment.Port = nextFreePort(*o.Managment, "tcp")
		return nil
	}

	if err := o.createCertificate(); err != nil {
		return err
	}

	return o.createConfig()
}

func (o *OpenVPN) createConfig() error {
	file, err := os.Create(filepath.Join(o.Path, path.RoleConfig(o.Role)))
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := statik.ReadFile(path.Config.ServerConfigTemplate)
	if err != nil {
		return err
	}

	templ, err := template.New("ovpnTemplate").Parse(string(data))
	if err != nil {
		return err
	}

	// Set dynamic port.
	o.Managment.Port = nextFreePort(*o.Managment, "tcp")
	o.Host.Port = nextFreePort(*o.Host, o.Proto)

	if strings.EqualFold(o.Proto, "tcp") {
		o.Proto = fmt.Sprintf("%s-%s", o.Proto, o.Role)
	}

	if !o.IsWindows {
		o.User, o.Group, err = getUserGroup()
		if err != nil {
			return err
		}
	}

	return templ.Execute(file, &o)
}

// RemoveConfig removes openvpn configuration.
func (o *OpenVPN) RemoveConfig() error {
	if o.isClient() {
		return nil
	}

	pathsToRemove := []string{
		path.Config.DHParam,
		path.Config.CACertificate,
		path.Config.CAKey,
		path.RoleCertificate(o.Role),
		path.RoleKey(o.Role),
		path.RoleConfig(o.Role),
	}
	for _, path := range pathsToRemove {
		os.Remove(filepath.Join(o.Path, path))
	}

	return nil
}

func (o *OpenVPN) createCertificate() error {
	p := filepath.Join(o.Path, "config")
	t := time.Now().AddDate(o.Validity.Year,
		o.Validity.Month, o.Validity.Day)
	if err := buildServerCertificate(p, t); err != nil {
		return err
	}

	// Generate Diffie Hellman param.
	ossl := filepath.Join(o.Path, path.Config.OpenSSL)
	dh := filepath.Join(p, "dh2048.pem")
	err := exec.Command(ossl, "dhparam", "-out", dh, "2048").Run()
	if err != nil {
		cmd := exec.Command("openssl", "dhparam", "-out", dh, "2048")
		return cmd.Run()
	}
	return nil
}

func (o *OpenVPN) isClient() bool {
	return !strings.EqualFold(o.Role, "server")
}

// InstallService installs a openvpn service.
func (o *OpenVPN) InstallService() (string, error) {
	if o.isClient() {
		return "", nil
	}

	var dependencies []string
	o.Service = serviceName(path.Config.OVPN, o.Path)
	descr := fmt.Sprintf("Privatix %s OpenVPN %s", o.Role, hash(o.Path))

	if o.IsWindows {
		o.Service = fmt.Sprintf("Privatix OpenVPN %s", hash(o.Path))
		dependencies = []string{"tap0901", "dhcp"}
	}

	service, err := daemon.New(o.Service, descr, dependencies...)
	if err != nil {
		return "", err
	}

	if str, err := service.Install("run", "-workdir", o.Path); err != nil {
		return str, err
	}

	if !o.IsWindows {
		return "", nil
	}

	script := filepath.Join(o.Path, path.Config.PowerShellVpnFirewall)
	ovpn := filepath.Join(o.Path, path.Config.OpenVPN+".exe")
	args := buildPowerShellArgs(script, "-Create",
		"-ServiceName", strings.Join(strings.Fields(o.Service), "_"),
		"-ProgramPath", ovpn,
		"-Port", strconv.Itoa(o.Host.Port), "-Protocol", o.Proto[:3])

	return "", runPowerShellCommand(args...)
}

// StartService starts openvpn service.
func (o *OpenVPN) StartService() (string, error) {
	if o.isClient() {
		return "", nil
	}

	service, err := daemon.New(o.Service, "")
	if err != nil {
		return "", err
	}
	return service.Start()
}

// RunService executes openvpn service.
func (o *OpenVPN) RunService() (string, error) {
	if o.isClient() {
		return "", nil
	}

	service, err := daemon.New(o.Service, "")
	if err != nil {
		return "", err
	}

	return service.Run(&execute{Path: o.Path, Role: o.Role,
		Type: path.Config.OVPN})
}

// StopService stops openvpn service.
func (o *OpenVPN) StopService() (string, error) {
	if o.isClient() {
		return "", nil
	}

	service, err := daemon.New(o.Service, "")
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

// RemoveService removes the openvpn service.
func (o *OpenVPN) RemoveService() (string, error) {
	if o.isClient() {
		return "", nil
	}

	if o.IsWindows {
		script := filepath.Join(o.Path, path.Config.PowerShellVpnFirewall)
		args := buildPowerShellArgs(script, "-Remove",
			"-ServiceName", strings.Join(strings.Fields(o.Service), "_"))
		if err := runPowerShellCommand(args...); err != nil {
			return "", err
		}
	}

	service, err := daemon.New(o.Service, "")
	if err != nil {
		return "", err
	}
	return service.Remove()
}
