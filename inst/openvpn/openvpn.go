package openvpn

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
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
	Port      int
	Proto     string
	Host      *host
	Managment *host
	Server    *host
	Service   string
	Validity  *validity
	IsWindows bool
}

type validity struct {
	Year  int
	Month int
	Day   int
}

type host struct {
	IP       string
	Port     int
	Mask     string
	Protocol string
}

// NewOpenVPN creates a default OpenVPN configuration.
func NewOpenVPN() *OpenVPN {
	return &OpenVPN{
		Path:  ".",
		Tap:   &tapInterface{},
		Role:  "server",
		Proto: "udp",
		Host: &host{
			IP:       "0.0.0.0",
			Port:     443,
			Protocol: "tcp",
		},
		Managment: &host{
			IP:       "127.0.0.1",
			Port:     7505,
			Protocol: "tcp",
		},
		Server: &host{
			IP:   "10.217.3.0",
			Mask: "255.255.255.0",
		},
		Validity: &validity{
			Year: 10,
		},
		IsWindows: strings.EqualFold(runtime.GOOS, "windows"),
	}
}

// InstallTap installs a new tap interface.
func (o *OpenVPN) InstallTap() (err error) {
	if o.IsWindows {
		o.Tap, err = installTAP(o.Path, o.Role)
	}
	return err
}

// RemoveTap removes the tap interface.
func (o *OpenVPN) RemoveTap() (err error) {
	if o.IsWindows {
		err = o.Tap.remove(o.Path)
	}
	return err
}

// Configurate configurates openvpn config files.
func (o *OpenVPN) Configurate() error {
	if o.isClient() {
		return nil
	}

	if err := o.createCertificate(); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(o.Path, path.RoleConfig(o.Role)))
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := statik.ReadFile(path.ServerConfigTemplate)
	if err != nil {
		return err
	}

	templ, err := template.New("ovpnTemplate").Parse(string(data))
	if err != nil {
		return err
	}

	// Set dynamic port.
	o.Managment.Port = nextFreePort(*o.Managment)
	o.Host.Port = nextFreePort(*o.Host)

	return templ.Execute(file, &o)
}

// RemoveConfig removes openvpn configuration.
func (o *OpenVPN) RemoveConfig() error {
	if o.isClient() {
		return nil
	}

	pathsToRemove := []string{
		path.DHParam,
		path.CACertificate,
		path.CAKey,
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
	ossl := filepath.Join(o.Path, path.OpenSSL)
	dh := filepath.Join(p, "dh2048.pem")
	return exec.Command(ossl, "dhparam", "-out", dh, "2048").Run()
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
	o.Service = serviceName(o.Path)
	descr := fmt.Sprintf("dapp-openvpn %s %s", o.Service, o.Tap.Interface)

	if o.IsWindows {
		dependencies = []string{"tap0901", "dhcp"}
	}

	service, err := daemon.New(o.Service, descr, dependencies...)
	if err != nil {
		return "", err
	}

	return service.Install("run")
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

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	ovpn := filepath.Join(o.Path, path.OpenVPN)
	config := filepath.Join(o.Path, path.RoleConfig(o.Role))
	cmd := exec.Command(ovpn, "--config", config)

	if err := cmd.Run(); err != nil {
		return "failed to execute openvpn", err
	}

	// Waiting for interrupt by system signal.
	killSignal := <-interrupt
	return fmt.Sprintf("service exited. got signal: %v", killSignal), nil
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

	service, err := daemon.New(o.Service, "")
	if err != nil {
		return "", err
	}
	return service.Remove()
}
