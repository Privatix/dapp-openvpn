#!/usr/bin/env bash

if [ -z "${DAPP_OPENVPN_DIR}" ]
then
    MY_PATH="`dirname \"$0\"`" # relative bash file path
    DAPP_OPENVPN_DIR="`( cd \"$MY_PATH/..\" && pwd )`"  # absolutized and normalized dappctrl path
fi

GIT_COMMIT=$(git rev-list -1 HEAD | head -n 1)
GIT_RELEASE=$(git tag -l --points-at HEAD | head -n 1)

# if $GIT_RELEASE is zero:
GIT_RELEASE=${GIT_RELEASE:-$(git rev-parse --abbrev-ref HEAD | grep -o "[0-9]\{1,\}\.[0-9]\{1,\}\.[0-9]\{1,\}")}


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

go get -u -v github.com/rakyll/statik

echo
echo go generate
echo

go generate -x ./...

echo
echo go build
echo

if [[ ! -d "${GOPATH}/bin/" ]]; then
    mkdir "${GOPATH}/bin/" || exit 1
fi

go build -o $GOPATH/bin/${ADAPTER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${DAPP_OPENVPN_DIR}${ADAPTER_MAIN} || exit 1
echo $GOPATH/bin/${ADAPTER_NAME}

go build -o $GOPATH/bin/${INSTALLER_NAME} -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${DAPP_OPENVPN_DIR}${INSTALLER_MAIN} || exit 1
echo $GOPATH/bin/${INSTALLER_NAME}

go build -o $GOPATH/bin/${OPENVPN_INSTALLER_NAME} -ldflags \
    "-X main.Commit=$GIT_COMMIT -X main.Version=$GIT_RELEASE" \
    ${DAPP_OPENVPN_DIR}${OPENVPN_INSTALLER_MAIN} || exit 1
echo $GOPATH/bin/${OPENVPN_INSTALLER_NAME}

echo
echo done
