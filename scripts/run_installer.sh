#!/usr/bin/env bash

if [ -z "${POSTGRES_PORT}" ]
then
    POSTGRES_PORT=5433
fi

# Example run installer
echo dapp-openvpn-inst \
 -rootdir="$GOPATH/src/github.com/privatix/dapp-openvpn/files/example" \
 -connstr="dbname=dappctrl host=localhost user=postgres \
  sslmode=disable port=${POSTGRES_PORT}" -setauth
