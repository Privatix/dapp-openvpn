#!/usr/bin/env bash

# Example run installer
../bin/installer \
 -rootdir="$GOPATH/src/github.com/privatix/dapp-openvpn/files/example" \
 -connstr="dbname=dappctrl host=localhost user=postgres \
  sslmode=disable port=5433" -setauth
