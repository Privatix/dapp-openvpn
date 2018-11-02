#!/usr/bin/env bash

if [ -z "${POSTGRES_PORT}" ]
then
    POSTGRES_PORT=5432
fi

EXAMPLE_DIR="$GOPATH/src/github.com/privatix/dapp-openvpn/files/example"
BIN_DIR="$GOPATH/src/github.com/privatix/dapp-openvpn/bin"

cp -a ${EXAMPLE_DIR}/. ${BIN_DIR}/example

# Example run installer
dapp-openvpn-inst \
 -rootdir=${BIN_DIR}/example \
 -connstr="dbname=dappctrl host=localhost user=postgres \
  sslmode=disable port=${POSTGRES_PORT}" -setauth
