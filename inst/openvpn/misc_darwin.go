package openvpn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/privatix/dapp-openvpn/inst/openvpn/path"
	"github.com/privatix/dapp-openvpn/statik"
)

func serviceName(prefix, path string) string {
	return fmt.Sprintf("io.privatix.%s_%s", prefix, hash(path))
}

func setRegValue(guid, name string) error {
	return nil
}

func daemonPath(name string) string {
	return filepath.Join("/Library/LaunchDaemons", name+".plist")
}

// createNatRules creates daemon on Mac, which configures NAT rules.
func createNatRules(p, server string, port int) error {
	name := serviceName("nat", p)
	file, err := os.Create(daemonPath(name))
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := statik.ReadFile("/ovpn/templates/mac-daemon.tpl")
	if err != nil {
		return err
	}

	templ, err := template.New("daemonTemplate").Parse(string(data))
	if err != nil {
		return err
	}

	type natRule struct {
		Name   string
		Script string
		Server string
		Port   int
	}

	script := filepath.Join(p, path.Config.NatScript)
	if err := os.Chmod(script, 0777); err != nil {
		return err
	}
	d := &natRule{
		Name:   name,
		Script: script,
		Server: server,
		Port:   port,
	}
	if err := templ.Execute(file, &d); err != nil {
		return err
	}

	return exec.Command("launchctl", "load", daemonPath(name)).Run()
}
