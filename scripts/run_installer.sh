#!/usr/bin/env bash

# Example run installer
../bin/installer/installer \
 -rootdir="$GOPATH/src/github.com/privatix/dapp-openvpn/files/example" \
 -agent=true -connstr="dbname=dappctrl host=localhost user=postgres \
  sslmode=disable port=5432" -setauth
