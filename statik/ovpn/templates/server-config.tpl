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
auth-user-pass-verify "{{.Path}}/bin/dappvpn.exe -config={{.Path}}/config/dappvpn.config.json" via-file
client-cert-not-required
username-as-common-name
client-connect "{{.Path}}/bin/dappvpn.exe -config={{.Path}}/config/dappvpn.config.json"
client-disconnect "{{.Path}}/bin/dappvpn.exe -config={{.Path}}/config/dappvpn.config.json"
script-security 3
tls-server
server {{.Server.IP}} {{.Server.Mask}}
push "route {{.Server.IP}} {{.Server.Mask}}"
push "redirect-gateway def1"
ifconfig-pool-persist {{.Path}}/config/ipp.txt
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
