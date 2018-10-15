package openvpn

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
)

const ovpnPrefix = "dapp_ovpn"

const serverTemplate = `
dev tun
dev-node "{{.Tap.Interface}}"
proto {{.Proto}}
port {{.Port}}
tls-server
server {{.ServerIP}} 255.255.255.0
comp-lzo
dh {{.Path}}/ssl/dh2048.pem
ca {{.Path}}/ssl/ca.crt
cert {{.Path}}/ssl/server.crt
key {{.Path}}/ssl/server.key
tls-auth {{.Path}}/ssl/ta.key 0
tun-mtu 1500
tun-mtu-extra 32
mssfix 1450
keepalive 10 120
status {{.Path}}/log/openvpn-status.log
log {{.Path}}/log/openvpn.log
verb 3
`

const clientTemplate = `
dev tun
dev-node "{{.Tap.Interface}}"
proto {{.Proto}}
tls-client
remote 11.8.0.1 1194
comp-lzo
ca {{.Path}}/ca.crt
tls-auth {{.Path}}/ta.key 1
tun-mtu 1500
tun-mtu-extra 32
mssfix 1450
keepalive 10 120
status {{.Path}}/log/openvpn-status.log
log {{.Path}}/log/openvpn.log
verb 3
`

type service struct {
	ID          string
	GUID        string
	Name        string
	Description string
	Command     string
	Args        []string
	AutoStart   bool
}

func diff(a, b []string) string {
	for i, v := range b {
		if len(a) <= i || a[i] != v {
			return v
		}
	}
	return ""
}

func ovpnName(path string) string {
	return fmt.Sprintf("%s_%s", ovpnPrefix, hash(path))
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func isServiceRun(service string) (bool, error) {
	output, err := exec.Command("sc", "queryex", service).CombinedOutput()

	if err != nil {
		return false, err
	}

	return strings.Contains(string(output), "RUNNING"), nil
}

func isServiceStop(service string) (bool, error) {
	output, err := exec.Command("sc", "queryex", service).CombinedOutput()

	if err != nil {
		return false, err
	}

	return strings.Contains(string(output), "STOPPED"), nil
}
