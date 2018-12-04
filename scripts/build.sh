#!/usr/bin/env bash

PROJECT=github.com/privatix/dapp-openvpn

if [ -z "${DAPP_OPENVPN_DIR}" ]
then
    DAPP_OPENVPN_DIR=$GOPATH/src/${PROJECT}
fi

GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_RELEASE=$(git tag -l --points-at HEAD)

ADAPTER_MAIN=/adapter
ADAPTER_NAME=dappvpn

INSTALLER_MAIN=/installer
INSTALLER_NAME=dapp-openvpn-inst

OPENVPN_INSTALLER_MAIN=/inst
OPENVPN_INSTALLER_NAME=openvpn-inst

cd "${DAPP_OPENVPN_DIR}" || exit 1

echo
echo go get
echo

go get -d -x ${PROJECT}/... || exit 1
go get -u -x gopkg.in/reform.v1/reform || exit 1
go get -u -x github.com/rakyll/statik || exit 1

echo
echo go dep
echo

if [ ! -f "${GOPATH}"/bin/dep ]; then
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
fi
rm -f Gopkg.lock
dep ensure || exit 1

echo
echo go generate
echo

go generate -x ./... || exit 1

echo
echo go build
echo


go build -o $GOPATH/bin/${ADAPTER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${ADAPTER_MAIN} || exit 1

go build -o $GOPATH/bin/${INSTALLER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${INSTALLER_MAIN} || exit 1

go build -o $GOPATH/bin/${OPENVPN_INSTALLER_NAME} -ldflags \
    "-X main.Commit=$GIT_COMMIT -X main.Version=$GIT_RELEASE" \
    ${PROJECT}${OPENVPN_INSTALLER_MAIN} || exit 1
