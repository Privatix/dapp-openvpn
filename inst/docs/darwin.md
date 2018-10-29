#### Install 
1. Choose `openvpn` version: server (`openvpn-mac-server.zip`) or client (`openvpn-mac-client.zip`)
2. Extract selected zip archive.
3. Run command shell (`terminal`).
4. Go to `bin` directory.
5. Execute command `sudo installer install -config ../config/installer.config.json`
6. Wait to end installation process. All steps print into console and write to `log`.
7. After successfully install will created `config/.env` with follow key-values: `DEVICE`, `INTERFACE`, `SERVICE`, `WORKDIR`
8. After successfully install you should configurate `openvpn`-configs in `config`.
9. Restart or manual run `sudo installer run`

#### Remove
1. Run command shell (`terminal`).
2. Go to `bin` directory.
3. Execute command `sudo installer remove`
4. Wait to end removing process. All steps print into console and write to `log`.

#### Help
1. Run command shell (`terminal`).
2. Go to `bin` directory.
3. Execute command `sudo installer help` to display help information.
