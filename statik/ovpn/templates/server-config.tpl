local {{.Host.IP}}
port {{.Host.Port}}
proto {{.Proto}}
dev tun
{{if not .IsWindows}}#{{end}}dev-node "{{.Tap.Interface}}"
ca "{{.Path}}/config/ca.crt"
cert "{{.Path}}/config/server.crt"
key "{{.Path}}/config/server.key"
dh "{{.Path}}/config/dh2048.pem"
management {{.Managment.IP}} {{.Managment.Port}}
auth-user-pass-verify "{{.Path}}/bin/dappvpn{{if .IsWindows}}.exe{{end}} -config={{.Path}}/config/adapter.config.json" via-file
verify-client-cert none
username-as-common-name
client-connect "{{.Path}}/bin/dappvpn{{if .IsWindows}}.exe{{end}} -config={{.Path}}/config/adapter.config.json"
client-disconnect "{{.Path}}/bin/dappvpn{{if .IsWindows}}.exe{{end}} -config={{.Path}}/config/adapter.config.json"
script-security 3
tls-server
server {{.Server.IP}} {{.Server.Mask}}
push "route {{.Server.IP}} {{.Server.Mask}}"
push "redirect-gateway def1"
ifconfig-pool-persist "{{.Path}}/config/ipp.txt"
keepalive 10 120
comp-lzo
persist-key
persist-tun
{{if .IsWindows}}#{{end}}user {{.User}}
{{if .IsWindows}}#{{end}}group {{.Group}}
status "{{.Path}}/log/openvpn-status.log"
log "{{.Path}}/log/server.log"
log-append "{{.Path}}/log/server-append.log"
verb 3
