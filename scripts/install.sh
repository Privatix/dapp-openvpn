#!/usr/bin/env bash

PROJECT=github.com/privatix/dapp-openvpn
PROJECT_PATH=$GOPATH/src/${PROJECT}
GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_RELEASE=$(git tag -l --points-at HEAD)
ADAPTER_MAIN=/adapter
INSTALLER_MAIN=/installer
ADAPTER_NAME=dappvpn
INSTALLER_NAME=installer

cd "${PROJECT_PATH}" || exit
dep ensure
go generate ./...

go build -o $GOPATH/bin/${ADAPTER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${ADAPTER_MAIN}

go build -o $GOPATH/bin/${INSTALLER_NAME}  -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${INSTALLER_MAIN}