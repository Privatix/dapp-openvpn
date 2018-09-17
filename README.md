[![Go report](http://goreportcard.com/badge/github.com/Privatix/dapp-openvpn)](https://goreportcard.com/report/github.com/Privatix/dapp-openvpn)
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
/scripts/install.sh
```

#### Additional steps for agent

On the agent side it necessary to perform the following steps:
[additional steps](https://github.com/Privatix/dapp-openvpn/wiki/Additional-steps-for-an-agent)

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
