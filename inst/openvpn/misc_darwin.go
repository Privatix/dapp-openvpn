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
	return fmt.Sprintf("io.privatix.%s_%s", prefix, hash(path))
}

func setRegValue(guid, name string) error {
	return nil
}

func daemonPath(name string) string {
	return filepath.Join("/Library/LaunchDaemons", name+".plist")
}

func createNatRules(p, server string, port int) error {
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

var daemonTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Disabled</key>
    <false/>
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
    </dict>
    <key>Label</key>
    <string>{{.Name}}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.Script}}</string>
	<string>on</string>
	<string>{{.Server}}</string>
	<string>{{.Port}}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
`
