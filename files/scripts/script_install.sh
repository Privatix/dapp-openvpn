#!/usr/bin/env bash

PROJECT=github.com/privatix/dapp-openvpn
GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_RELEASE=$(git tag -l --points-at HEAD)
ADAPTER_MAIN=/cmd/adapter
INSTALLER_MAIN=/cmd/installer

go install -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${ADAPTER_MAIN}

go install -ldflags "-X main.Commit=$GIT_COMMIT \
    -X main.Version=$GIT_RELEASE" -tags=notest \
    ${PROJECT}${INSTALLER_MAIN}
