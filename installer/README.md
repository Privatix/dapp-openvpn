# Installer (legacy)

Given installer is used only in ubuntu installation process.

## Features 

Installer is able to:

* add `dapp-openvpn` products into the database
* add password-salt to dapp-openvpn configs

## Usage

```bash
    POSTGRES_USER=postgres
    POSTGRES_PORT=5432
    POSTGRES_PASSWORD=some_password
    
    BIN_DIR="$GOPATH/src/github.com/privatix/dapp-openvpn/bin"
    
    connection_string="dbname=dappctrl host=localhost sslmode=disable \
user=${POSTGRES_USER} \
port=${POSTGRES_PORT} \
${POSTGRES_PASSWORD:+ password=${POSTGRES_PASSWORD}}"

    $GOPATH/bin/bin/dapp-openvpn-inst \
     -rootdir=${BIN_DIR}/example \
     -connstr="$connection_string" -setauth
```

## Authors

* [dzeckelev](https://github.com/dzeckelev)
