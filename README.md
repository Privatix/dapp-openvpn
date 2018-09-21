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

### Building Executables

Switch to the dapp-openvpn repository root directory.

You can build all binaries using the go tool, placing the resulting binary in `$GOPATH/bin`:

```bash
./scripts/install.sh
```

You can build all binaries using the go tool, placing the resulting binary in `./bin`:

```bash
./scripts/build.sh
```

#### Additional steps for agent

On the agent side it necessary to perform the following steps:
[additional steps](https://github.com/Privatix/dapp-openvpn/wiki/Additional-steps-for-an-agent)

### Running the agent service

- Start the `OpenVPN`-server.
- Start the `dapp-openvpn` in the background with the configuration provided.

## Command Line Options

### Installer

```bash
Usage of installer:
  -agent
        Whether to install agent
  -connstr string
        PostgreSQL connection string (default "user=postgres dbname=dappctrl sslmode=disable")
  -rootdir string
        Full path to root directory of service adapter
  -setauth
        Generate authentication credentials for service adapter
```

### dapppvn (adapter)

```bash
Usage of dappvpn:
  -channel string
        Channel ID for client mode
  -config string
        Configuration file (default "dappvpn.config.json")
  -version
        Prints current dappctrl version
```

## Tests

Run tests for all packages.

```bash
./scripts/run_tests.sh
```

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
