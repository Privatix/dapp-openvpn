# tunX | tapX | null TUN/TAP virtual network device
# ( X can be omitted for a dynamic device.)
dev tun
{{if .TapInterface}}dev-node {{.TapInterface}}{{end}}
# Use protocol tcp for communicating
# with remote host
proto {{if .Proto}}{{.Proto}}{{else}}tcp-client{{end}}

# Encrypt packets with AES-256-CBC algorithm
cipher {{if .Cipher}}{{.Cipher}}{{else}}AES-256-GCM{{end}}

# Enable TLS and assume client role
# during TLS handshake.
tls-client

# Remote host name or IP address
# with port number and protocol tcp
# for communicating
{{if .ServerAddress}}{{if .Port}}remote {{.ServerAddress}} {{.Port}}{{end}}{{end}}

# If hostname resolve fails for --remote,
# retry resolve for n seconds before failing.
# Set n to "infinite" to retry indefinitely.
resolv-retry 30

# Do not bind to local address and port.
# The IP stack will allocate a dynamic
# port for returning packets.
# Since the value of the dynamic port
# could not be known in advance by a peer,
# this option is only suitable for peers
# which will be initiating connections
# by using the --remote option.
nobind

# Don't close and reopen TUN/TAP device
# or run up/down scripts across SIGUSR1
# or --ping-restart restarts.
# SIGUSR1 is a restart signal similar
# to SIGHUP, but which offers
# finer-grained control over reset options.
persist-tun

# Don't re-read key files across
# SIGUSR1 or --ping-restart.
persist-key

# Trigger a SIGUSR1 restart after n seconds
# pass without reception of a ping
# or other packet from remote.
ping-restart {{if .PingRestart}}{{.PingRestart}}{{else}}25{{end}}

# Ping remote over the TCP/UDP control
# channel if no packets have been sent for
# at least n seconds
ping {{if .Ping}}{{.Ping}}{{else}}10{{end}}

# Authenticate with server using
# username/password in interactive mode
auth-user-pass {{.AccessFile}}

# Pull settings from server
pull
# But allow only restricted parameters and values
pull-filter accept "route {{.Server}}"
pull-filter accept "ifconfig {{.Server}}"
pull-filter accept "route-gateway {{.Server}}"
pull-filter accept "peer-id"
pull-filter accept "topology"
pull-filter accept "cipher AES-256-GCM"
pull-filter accept "ping"
pull-filter accept "compress"
pull-filter accept "dhcp-option DNS 8.8.8.8"
pull-filter accept "dhcp-option DNS 8.8.4.4"
pull-filter ignore ""

# route all traffic through VPN
redirect-gateway def1 block-local

# take n as the number of seconds
# to wait between connection retries
connect-retry {{if .ConnectRetry}}{{.ConnectRetry}}{{else}}5{{end}}

# Server CA certificate for TLS validation

<ca>
{{.Ca}}</ca>

# Enable compression on the VPN link.
# Don't enable this unless it is also
# enabled in the server config file.
compress {{if .Compress}}{{.Compress}}{{else}}lz4{{end}}

# Set log file verbosity.
verb 3
{{if .LogAppend}}log-append {{.LogAppend}}{{end}}

# Management interface settings
management 0.0.0.0 {{if .ManagementPort}}{{.ManagementPort}}{{else}}7506{{end}}
management-hold
management-signal

# Remap SIGUSR1 to SIGTERM to prevent holding in unconnected state
remap-usr1 SIGTERM

# Set up/down trigger scripts
script-security 2
{{if .UpScript}}up {{.UpScript}}{{end}}
{{if .DownScript}}down {{.DownScript}}{{end}}
