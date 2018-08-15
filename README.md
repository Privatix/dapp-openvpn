[![Go report](https://goreportcard.com/badge/github.com/Privatix/dapp-openvpn)](https://goreportcard.com/badge/github.com/Privatix/dapp-openvpn)
[![Maintainability](https://api.codeclimate.com/v1/badges/af4e29689d76d8ccf974/maintainability)](https://codeclimate.com/github/Privatix/dapp-openvpn/maintainability)
[![GoDoc](https://godoc.org/github.com/Privatix/dapp-openvpn?status.svg)](https://godoc.org/github.com/Privatix/dapp-openvpn)

# OpenVPN Service Adapter

OpenVPN service adapter is an executable which integrates OpenVPN as a service
with the Privatix controller.

## Getting started

These instructions will help you to build and configure the OpenVPN service
adapter.

### Prerequisites

- Install OpenVPN 2.4+.

### Installation

Build the adapter:

```bash
export GIT_COMMIT=$(git rev-list -1 HEAD) && \
export GIT_RELEASE=$(git tag -l --points-at HEAD) && \
  go install -ldflags "-X main.Commit=$GIT_COMMIT \
   -X main.Version=$GIT_RELEASE"
```

#### Additional steps for agent

Insert a new product into a database of the corresponding agent. Then modify
the adapter configuration file:

```bash
CONF_FILE=$HOME/go/src/github.com/privatix/dapp-openvpn/dapp-openvpn.config.json
LOCAL_CONF_FILE=$HOME/dappvpn.config.json
PRODUCT_ID=<uuid> # ID of a newly inserted product.
PRODUCT_PASS=<password> # Password of a newly inserted product.

jq ".Server.Username=\"$PRODUCT_ID\" | .Server.Password=\"$PRODUCT_PASS\"" $CONF_FILE > $LOCAL_CONF_FILE
```

Add the following lines to the `OpenVPN`-server configuration file
(substituting file paths):

```
auth-user-pass-verify "/path/to/dapp-openvpn -config=/path/to/local/config" via-file
client-connect "/path/to/dapp-openvpn -config=/path/to/local/config"
client-disconnect "/path/to/dapp-openvpn -config=/path/to/local/config"
script-security 3
management localhost 7505
```

### Running the agent service

- Start the `OpenVPN`-server.
- Start the `dapp-openvpn` in the background with the configuration provided.

# Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/Privatix/dapp-openvpn/tags).

## Authors

* [ababo](https://github.com/ababo)
* [dzeckelev](https://github.com/dzeckelev)

See also the list of [contributors](https://github.com/Privatix/dapp-openvpn/contributors) 
who participated in this project.

# License

This project is licensed under the **GPL-3.0 License** - see the 
[COPYING](COPYING) file for details.
