#!/usr/bin/env bash

# Run tests
go test -v github.com/privatix/dapp-openvpn/... \
    -config="$GOPATH/src/github.com/privatix/dapp-openvpn/files/test/test.conf"
