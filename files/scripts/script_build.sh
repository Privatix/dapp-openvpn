#!/usr/bin/env bash

PROJECT=github.com/privatix/dapp-openvpn
PROJECT_PATH=$GOPATH/src/${PROJECT}
PROJECT_BIN=${PROJECT_PATH}/bin
GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_RELEASE=$(git tag -l --points-at HEAD)
ADAPTER_MAIN=/cmd/adapter
INSTALLER_MAIN=/cmd/installer

rm -drf "${PROJECT_BIN}"
mkdir -p "${PROJECT_BIN}/adapter" "${PROJECT_BIN}/installer"

go build -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    -o "${PROJECT_BIN}/adapter/adapter" \
    "${PROJECT}${ADAPTER_MAIN}"

go build -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    -o "${PROJECT_BIN}/installer/installer" \
    "${PROJECT}${INSTALLER_MAIN}"
