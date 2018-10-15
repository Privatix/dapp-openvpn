package openvpn

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

const ovpnPrefix = "dapp_ovpn"

const serverTemplate = `
local {{.Host.IP}}
port {{.Host.Port}}
proto {{.Proto}}
dev tun
dev-node "{{.Tap.Interface}}"
ca {{.Path}}/config/ca.crt
cert {{.Path}}/config/server.crt
key {{.Path}}/config/server.key
dh {{.Path}}/config/dh2048.pem
management {{.Managment.IP}} {{.Managment.Port}}
auth-user-pass-verify "{{.Path}}/bin/dappvpn -config={{.Path}}/config/dappvpn.config.json" via-file
client-cert-not-required
username-as-common-name
client-connect "{{.Path}}/bin/dappvpn -config={{.Path}}/config/dappvpn.config.json"
client-disconnect "{{.Path}}/bin/dappvpn -config={{.Path}}/config/dappvpn.config.json"
script-security 3
tls-server
server {{.Server.IP}} {{.Server.Mask}}
push "route {{.Server.IP}} {{.Server.Mask}}"
push "redirect-gateway def1"
ifconfig-pool-persist ipp.txt
keepalive 10 120
comp-lzo
persist-key
persist-tun
user root
group root
status {{.Path}}/log/openvpn-status.log
log {{.Path}}/log/server.log
log-append {{.Path}}/log/server-append.log
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

func freePort(h host) int {
	port := h.Port
	for i := port; i < 65535; i++ {
		ln, err := net.Listen(h.Protocol, h.IP+":"+strconv.Itoa(i))
		if err != nil {
			continue
		}

		if err := ln.Close(); err != nil {
			continue
		}
		port = i
		break
	}

	return port
}
