# Builder

Builder is a tool to package service plug-in to compatible with dapp-installer form.

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

where

`keystore` - is private key of service plug-in maintainer in ethereum keystore format

`auth` - passphrase for keystore file decryption

`agent` - Agent OR Client service plug-in user role

Resulting archive and descriptor files will be placed in `build` directory.

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
