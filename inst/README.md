# Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Prerequisites

Install [Golang](https://golang.org/doc/install). Make sure that `$GOPATH/bin` is added to system path `$PATH`.

## Installation steps

Clone the `dapp-openvpn` repository using git:

```bash
git clone https://github.com/Privatix/dapp-openvpn.git
cd dapp-openvpn
git checkout master
```

Build `inst` package:

```bash
go get -d github.com/Privatix/dapp-openvpn
go get github.com/rakyll/statik

go generate ./...

cd inst
go build -o installer
```

# Usage

Place `installer` and `installer.config.json` in the appropriate folder, according to the package distribution folder structure.

Configurate `installer.config.json` (see [details on the configuration description](./docs/config.md)).

Simply run `installer <COMMAND>`

`installer` or `installer --help` will show usage and a list of commands:

```
Usage:
	installer [command] [flags]
Available Commands:
	install     Install product package
	remove      Remove product package
	run         Run service
	start       Start service
	stop        Stop service
Flags:
	--help      Display help information
	--version   Display the current version of this CLI
Use "installer [command] --help" for more information about a command.
 ```
More information about [installation](./docs/index.md).

# Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details of our code of conduct and of the process of submitting pull requests.

## Versioning

We use [semantic versioning](http://semver.org/), see available [versions](https://github.com/Privatix/dappctrl/tags).

## Authors

* [ubozov](https://github.com/ubozov)

See also the list of [contributors](https://github.com/Privatix/dapp-openvpn/contributors) who participated in this project.

# License

This project is licensed under the **GPL-3.0 License** - see the [COPYING](COPYING) file for details.
