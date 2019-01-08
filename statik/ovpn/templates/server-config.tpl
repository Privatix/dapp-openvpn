local {{.Host.IP}}
port {{.Host.Port}}
proto {{.Proto}}
dev tun
{{if not .IsWindows}}#{{end}}dev-node {{.Tap.GUID}}
ca "config/ca.crt"
cert "config/server.crt"
key "config/server.key"
dh "config/dh2048.pem"
management {{.Managment.IP}} {{.Managment.Port}}
auth-user-pass-verify "bin/dappvpn{{if .IsWindows}}.exe{{end}} -config config/adapter.config.json" via-file
verify-client-cert none
username-as-common-name
client-connect "bin/dappvpn{{if .IsWindows}}.exe{{end}} -config config/adapter.config.json"
client-disconnect "bin/dappvpn{{if .IsWindows}}.exe{{end}} -config config/adapter.config.json"
script-security 3
tls-server
server {{.Server.IP}} {{.Server.Mask}}
push "route {{.Server.IP}} {{.Server.Mask}}"
push "dhcp-option DNS 8.8.8.8"
push "dhcp-option DNS 8.8.4.4"
push "redirect-gateway def1"
ifconfig-pool-persist "config/ipp.txt"
keepalive 10 120
comp-lzo
persist-key
persist-tun
{{if .IsWindows}}#{{end}}user {{.User}}
{{if .IsWindows}}#{{end}}group {{.Group}}
status "log/openvpn-status.log"
log "log/server.log"
log-append "log/server-append.log"
verb 3
