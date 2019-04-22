#!/usr/bin/env bash

# Run tests
go test -v github.com/privatix/dapp-openvpn/... \
    -config="$1"
