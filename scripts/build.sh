#!/usr/bin/env bash

PROJECT=github.com/privatix/dapp-openvpn
PROJECT_PATH=$GOPATH/src/${PROJECT}

GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_RELEASE=$(git tag -l --points-at HEAD)

ADAPTER_MAIN=/adapter
INSTALLER_MAIN=/installer
OPENVPN_INSTALLER_MAIN=/inst

ADAPTER_NAME=dappvpn
INSTALLER_NAME=dapp-openvpn-inst

OPENVPN_INSTALLER_NAME=openvpn-inst

cd "${PROJECT_PATH}" || exit

go get -d ${PROJECT}/...
go get -u gopkg.in/reform.v1/reform
go get -u github.com/rakyll/statik

curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure

go generate ./...

go build -o $GOPATH/bin/${ADAPTER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${ADAPTER_MAIN}

go build -o $GOPATH/bin/${INSTALLER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${INSTALLER_MAIN}

go build -o $GOPATH/bin/${OPENVPN_INSTALLER_NAME} -tags=notest \
    ${PROJECT}${OPENVPN_INSTALLER_MAIN}
