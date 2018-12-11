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

cd "${DAPP_OPENVPN_DIR}"

echo
echo go get
echo

go get -d -v ${PROJECT}/...
go get -u -v gopkg.in/reform.v1/reform
go get -u -v github.com/rakyll/statik

echo
echo dep ensure
echo

if [ ! -f "${GOPATH}"/bin/dep ]; then
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
fi
rm -f Gopkg.lock
dep ensure -v

echo
echo go generate
echo

go generate -x ./...

echo
echo go build
echo

echo $GOPATH/bin/${ADAPTER_NAME}
go build -o $GOPATH/bin/${ADAPTER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${ADAPTER_MAIN} || exit 1

echo $GOPATH/bin/${INSTALLER_NAME}
go build -o $GOPATH/bin/${INSTALLER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${INSTALLER_MAIN} || exit 1

echo $GOPATH/bin/${OPENVPN_INSTALLER_NAME}
go build -o $GOPATH/bin/${OPENVPN_INSTALLER_NAME} -ldflags \
    "-X main.Commit=$GIT_COMMIT -X main.Version=$GIT_RELEASE" \
    ${PROJECT}${OPENVPN_INSTALLER_MAIN} || exit 1

echo
echo done
