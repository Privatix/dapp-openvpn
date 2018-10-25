package openvpn

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

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
	}
}

// InstallTap installs a new tap interface.
func (o *OpenVPN) InstallTap() (err error) {
	o.Tap, err = installTAP(o.Path, o.Role)
	return err
}

// RemoveTap removes the tap interface.
func (o *OpenVPN) RemoveTap() (err error) {
	return o.Tap.remove(o.Path)
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

// RegisterService registries a openvpn service.
func (o *OpenVPN) RegisterService() error {
	o.Service = ovpnName(o.Path)
	ovpnsvc := filepath.Join(o.Path, path.ServiceWrapper)
	s := &service{
		ID:          o.Service,
		GUID:        ovpnsvc,
		Name:        o.Tap.Interface,
		Description: "dapp openvpn " + o.Service,
		Command:     filepath.Join(o.Path, path.OpenVPN),
		Args: []string{
			"--config",
			filepath.Join(o.Path, path.RoleConfig(o.Role)),
		},
		AutoStart: true,
	}

	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	fileName := filepath.Join(o.Path, path.ServiceWrapperConfig)
	if err := ioutil.WriteFile(fileName, bytes, 0644); err != nil {
		return err
	}

	cmd := exec.Command("sc", "create", o.Service,
		"binpath="+ovpnsvc+" -config "+fileName,
		"type=own", "start=auto", "depend=tap0901/dhcp")
	return cmd.Run()
}

// StopService stops openvpn service.
func (o *OpenVPN) StopService() error {
	ok, err := isServiceRun(o.Service)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	return exec.Command("sc", "stop", o.Service).Run()
}

// RemoveService removes the openvpn service.
func (o *OpenVPN) RemoveService() error {
	return exec.Command("sc", "delete", o.Service).Run()
}
