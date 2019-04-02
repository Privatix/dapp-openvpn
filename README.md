[![Go report](http://goreportcard.com/badge/github.com/Privatix/dapp-openvpn)](https://goreportcard.com/report/github.com/Privatix/dapp-openvpn)
[![Maintainability](https://api.codeclimate.com/v1/badges/af4e29689d76d8ccf974/maintainability)](https://codeclimate.com/github/Privatix/dapp-openvpn/maintainability)
[![pullreminders](https://pullreminders.com/badge.svg)](https://pullreminders.com?ref=badge)

# OpenVPN Service Plug-in

OpenVPN [service plug-in](https://github.com/Privatix/privatix/blob/master/doc/service_plug-in.md) allows Agents and Clients to buy and sell their internet traffic in form of VPN service without 3rd party.

    This service plug-in is a PoC service plugin-in for Privatix core.

## Custom integration includes

-   Start and stop of session by Privatix Core
-   OpenVPN client sessions authentication by Privatix Core
-   OpenVPN traffic usage reporting to Privatix core (is a must for automatic payments)
-   Push OpenVPN server configuration
-   Traffic shaping based on offering parameters (Ubuntu only)

## Benefits from Privatix Core

-   Automatic billing
-   Automatic payment
-   Access control based on billing
-   Automatic credentials delivery
-   Automatic configuration delivery
-   Anytime increase of deposit
-   Privatix GUI for service control

## Service plug-in components:

-   Templates (offering and access)
-   Service adapter (with access to OpenVPN and Privatix core)
-   OpenVPN software (with management interface)

## Getting started

These instructions will help you to build and configure the OpenVPN service
adapter.

### Prerequisites

-   Install OpenVPN 2.4+.

### Building Executables

Switch to the dapp-openvpn repository root directory.

You can build all binaries using the go tool, placing the 
resulting binary in `$GOPATH/bin`:

```bash
./scripts/toml.sh ./Gopkg.toml.template > ./Gopkg.toml
./scripts/build.sh
```

### Run installer

Installer is responsible for:

1. import templates: offering, access (to Privatix core database)
2. import product and link it to templates (to Privatix core database)
3. generate adapter config with proper authentication (same as in Privatix core database database product table)

Install and register `OpenVPN service plug-in` in `Privatix core` using the following script:

```bash
./scripts/run_installer.sh
```

#### Additional steps for OpenVPN server configuration

On the agent side it necessary to perform the following steps:
[additional steps](https://github.com/Privatix/dapp-openvpn/wiki/Additional-steps-for-an-agent)

### Running the agent service

-   Start the `OpenVPN`-server.
-   Start the `adapter -config adapter.config.json` in the background with the configuration provided by the installer.

### Running the client service

-   Start the `adapter -config adapter.config.json` in the background with the configuration provided by the installer.

#### Install statik

```
go get github.com/rakyll/statik
```

Generate statik filesystem:

```
go generate ./...
```

## Command Line Options

### dapp-openvpn-inst

```bash
Usage of dapp-openvpn-inst:
  -connstr string
        PostgreSQL connection string (default "user=postgres dbname=dappctrl sslmode=disable")
  -rootdir string
        Full path to root directory of service adapter
  -setauth
        Generate authentication credentials for service adapter
```

### dapp-openvpn (adapter)

```bash
Usage of dapp-openvpn:
  -channel string
        Channel ID for client mode
  -config string
        Configuration file (default "adapter.config.json")
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

-   [ababo](https://github.com/ababo)
-   [dzeckelev](https://github.com/dzeckelev)

See also the list of [contributors](https://github.com/Privatix/dapp-openvpn/contributors)
who participated in this project.

# License

This project is licensed under the **GPL-3.0 License** - see the
[COPYING](COPYING) file for details.
