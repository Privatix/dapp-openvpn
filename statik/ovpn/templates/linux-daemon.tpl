[Unit]
Description={{.Name}}
After=syslog.target network-online.target 
Wants=network-online.target
After=syslog.target
After=postgresql.service

[Service]
Type=onshot
ExecStart={{.Script}} on {{.Server}}
ExecStop={{.Script}} off {{.Server}}
Restart=on-failure
RemainAfterExit=yes
User=root
Group=root
StandardOutput=syslog
StandardError=syslog

[Install]
WantedBy=multi-user.target
