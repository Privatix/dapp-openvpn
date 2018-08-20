#!/usr/bin/env bash
DAPP_OPENVPN=github.com/privatix/dapp-openvpn
DAPP_OPENVPN_DIR=$HOME/go/src/${DAPP_OPENVPN}

go generate ${DAPP_OPENVPN}/...

GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_RELEASE=$(git tag -l --points-at HEAD)

export GIT_COMMIT
export GIT_RELEASE

go install -ldflags "-X main.Commit=$GIT_COMMIT -X main.Version=$GIT_RELEASE"
