#!/bin/bash

wget -O - https://swupdate.openvpn.net/repos/repo-public.gpg|apt-key add -
echo "deb http://build.openvpn.net/debian/openvpn/release/2.4 stretch main" > /etc/apt/sources.list.d/openvpn-aptrepo.list

apt-get update && apt-get install -y openvpn

/product/73e17130-2a1d-4f7d-97a8-93a9aaa6f10d/bin/inst update --config /product/73e17130-2a1d-4f7d-97a8-93a9aaa6f10d/config/installer.client.config.json

