#!/usr/bin/env bash

PROJECT=github.com/privatix/dapp-openvpn
PROJECT_PATH=$GOPATH/src/${PROJECT}
PROJECT_BIN=${PROJECT_PATH}/bin
GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_RELEASE=$(git tag -l --points-at HEAD)
ADAPTER_MAIN=/adapter
INSTALLER_MAIN=/installer
ADAPTER_NAME=dappvpn
INSTALLER_NAME=installer

cd "${PROJECT_PATH}" || exit
dep ensure
go generate ./...

rm -drf "${PROJECT_BIN}"

go build -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    -o "${PROJECT_BIN}/${ADAPTER_NAME}" \
    "${PROJECT}${ADAPTER_MAIN}"

go build -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    -o "${PROJECT_BIN}/${INSTALLER_NAME}" \
    "${PROJECT}${INSTALLER_MAIN}"
