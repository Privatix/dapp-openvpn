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

### Run installer
Installer is responsible for:
1. import templates: offering, access (to dappctrl database)
2. import product and link it to templates (to dappctrl database)
3. generate adapter config with proper authentication (same as in dappctrl database product table)
To install `dapp-openvpn` to `dappctrl`, please run the following script:

```bash
./scripts/run_installer.sh
```

#### Additional steps for agent

On the agent side it necessary to perform the following steps:
[additional steps](https://github.com/Privatix/dapp-openvpn/wiki/Additional-steps-for-an-agent)

### Running the agent service

- Start the `OpenVPN`-server.
- Start the `dapp-openvpn` in the background with the configuration provided by the installer.

## Build package

The package can be compiled for a specific operating system.
A `xgo`() is used to use cross-platform compilation.

To create package archive and descriptor:

#### Install dep dependency management tool

linux:
```bash
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

macos:
```bash
brew install dep
brew upgrade dep
```

windows:

Latest release can be downloaded from this link: https://github.com/golang/dep/releases

Run dep application:

Go to the root directory of the project and run command:
```bash
dep ensure
```

#### Install statik

```
go get github.com/rakyll/statik
```

Generate statik filesystem:
```
go generate ./...
```

#### Run builder
```
go run builder.go -agent=false \
                  -os=macos \
                  -keystore=/home/user/pk \
                  -auth="qwerty" \
                  -version=1_1 \
                  -min_core_version=0_123 \
                  -max_core_version=1_123
```

Archive and descriptor can be found in `build` directory.

## Command Line Options

### Installer

```bash
Usage of installer:
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
        Configuration file (default "dappvpn.config.json")
  -version
        Prints current dappctrl version
```

### builder
```bash
Usage of builder for dapp-openvpn:
  -agent
        Whether to install agent.
  -auth string
        Password to decrypt JSON private key.
  -keystore string
        Full path to JSON private key file.
  -max_core_version string
        Maximum version of Privatix core application for compatibility.
  -min_core_version string
        Minimal version of Privatix core application for compatibility. (default "undefined")
  -os string
        Target OS: linux, windows or macos (xgo usage). If is empty, a package will be created for a current operating system.
  -version string
        Product package distributive version. (default "undefined")

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
