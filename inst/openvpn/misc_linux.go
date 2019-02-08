package openvpn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/privatix/dapp-openvpn/inst/openvpn/path"
)

func serviceName(prefix, path string) string {
	return fmt.Sprintf("%s_%s", prefix, hash(path))
}

func setRegValue(guid, name string) error {
	return nil
}

func daemonPath(name string) string {
	return filepath.Join("/etc/systemd/system/", name+".service")
}

func createNatRules(p, server string, forwardingState int) error {
	name := serviceName("nat", p)
	file, err := os.Create(daemonPath(name))
	if err != nil {
		return err
	}
	defer file.Close()

	templ, err := template.New("daemonTemplate").Parse(daemonTemplate)
	if err != nil {
		return err
	}

	type natRule struct {
		Name            string
		Script          string
		Server          string
		ForwardingState int
	}

	script := filepath.Join(p, path.Config.NatScript)
	if err := os.Chmod(script, 0777); err != nil {
		return err
	}
	d := &natRule{
		Name:            name,
		Script:          script,
		Server:          server,
		ForwardingState: forwardingState,
	}
	if err := templ.Execute(file, &d); err != nil {
		return err
	}

	return exec.Command("systemctl", "enable", daemonPath(name)).Run()
}

var daemonTemplate = `[Unit]
Description={{.Name}}
After=syslog.target network-online.target 
Wants=network-online.target
After=syslog.target
After=postgresql.service

[Service]
Type=onshot
ExecStart={{.Script}} on {{.Server}}
ExecStop={{.Script}} off {{.ForwardingState}} {{.Server}}
Restart=on-failure
RemainAfterExit=yes
User=root
Group=root
StandardOutput=syslog
StandardError=syslog

[Install]
WantedBy=multi-user.target
`
