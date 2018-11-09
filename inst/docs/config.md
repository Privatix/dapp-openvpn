`config/installer.config.json` definition:

```
{
    Path:           workdir path, by default "."
    Role:           role: server - client, by default "server"
    Proto:          proto: udp - tcp, by default "udp"
    Host:           Server parameters
        IP:         address, by default "0.0.0.0"
        Port:       port, by default 443
    Managment:      managment interface	
        IP:         address, by default "127.0.0.1"
        Port:       port by default 7505
    Server:         VPN parameters
        IP:         address, by default "10.217.3.0",
        Mask:       subnet mask, by default "255.255.255.0"
    Validity        validity date to certificates and keys
        Year:       year, by default 10
        Month:      month, by default 0
        Day:        day, by default 0
}
```
