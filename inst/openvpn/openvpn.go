package openvpn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// OpenVPN has a openvpn configuration.
type OpenVPN struct {
	Path string
	Role string
	Tap  *tapInterface

	Port      int
	Proto     string
	Interface string
	ServerIP  string
}

// NewOpenVPN creates a default OpenVPN configuration.
func NewOpenVPN() *OpenVPN {
	return &OpenVPN{
		Path:     ".",
		Tap:      &tapInterface{},
		Role:     "server",
		Port:     1194,
		Proto:    "udp",
		ServerIP: "10.8.0.0",
	}
}

// DeviceID returns a register openvpn device ID.
func (o *OpenVPN) DeviceID() (string, error) {
	svcName := ovpnName(o.Path)
	return deviceID(svcName)
}

// InstallTap installs a new tap interface.
func (o *OpenVPN) InstallTap() (err error) {
	o.Tap, err = installTAP(o.Path)
	return err
}

// RemoveTap removes the tap interface.
func (o *OpenVPN) RemoveTap() (err error) {
	return o.Tap.remove(o.Path)
}

// Configurate configurates openvpn config files.
func (o *OpenVPN) Configurate() error {
	if err := o.createSertificate(); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(o.Path, "config/openvpn.conf"))
	if err != nil {
		return err
	}
	defer file.Close()

	temp := serverTemplate
	if o.isClient() {
		temp = clientTemplate
	}

	templ, err := template.New("ovpnTemplate").Parse(temp)
	if err != nil {
		return err
	}

	return templ.Execute(file, &o)
}

func (o *OpenVPN) createSertificate() error {
	if o.isClient() {
		return nil
	}

	path := filepath.Join(o.Path, "ssl")
	if err := buildServerCertificate(path); err != nil {
		return err
	}

	//generate DH param
	ossl := filepath.Join(o.Path, "bin/openssl")
	dh := filepath.Join(path, "dh2048.pem")
	err := exec.Command(ossl, "dhparam", "-out", dh, "2048").Run()
	if err != nil {
		return err
	}

	// generate secret key
	ovpn := filepath.Join(o.Path, "bin/openvpn")
	ta := filepath.Join(path, "ta.key")
	cmd := exec.Command(ovpn, "--genkey", "--secret", ta)

	return cmd.Run()
}

func (o *OpenVPN) isClient() bool {
	return !strings.EqualFold(o.Role, "server")
}

// RegisterService registries a openvpn service.
func (o *OpenVPN) RegisterService() error {
	svcName := ovpnName(o.Path)
	ovpnsvc := filepath.Join(o.Path, "bin/winsvc.exe")
	s := &service{
		ID:          svcName,
		GUID:        ovpnsvc,
		Name:        svcName,
		Description: "dapp openvpn " + svcName,
		Command:     filepath.Join(o.Path, "bin/openvpn"),
		Args:        []string{"--config", filepath.Join(o.Path, "config/openvpn.conf")},
		AutoStart:   true,
	}

	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	fileName := filepath.Join(o.Path, "bin/winsvc.config.json")
	if err := ioutil.WriteFile(fileName, bytes, 0644); err != nil {
		return err
	}

	cmd := exec.Command("sc", "create", svcName, "binpath="+ovpnsvc,
		"type=own", "start=auto", "depend=tap0901/dhcp")
	return cmd.Run()
}

// StopService stops openvpn service.
func (o *OpenVPN) StopService() error {
	service := ovpnName(o.Path)
	ok, err := isServiceRun(service)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	return exec.Command("sc", "stop", service).Run()
}

// RemoveService removes the openvpn service.
func (o *OpenVPN) RemoveService() error {
	service := ovpnName(o.Path)
	fmt.Println(service)
	return exec.Command("sc", "delete", service).Run()
}
