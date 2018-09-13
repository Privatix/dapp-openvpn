#!/usr/bin/env bash

# Run tests
go test github.com/privatix/dapp-openvpn/... \
    -config="$GOPATH/src/github.com/privatix/dapp-openvpn/files/test/test.conf"
