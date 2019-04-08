#### Install 
1. Choose `openvpn` version: server (`openvpn-mac-server.zip`) or client (`openvpn-mac-client.zip`)
2. Extract selected zip archive.
3. Run command shell (`terminal`).
4. Go to `bin` directory.
5. Execute command `sudo installer install -config ../config/installer.config.json`
6. Wait to end installation process. All steps print into console and write to `log`.
7. After successfull installation `config/.env` file with following key-values will be created: `DEVICE`, `INTERFACE`, `SERVICE`, `WORKDIR`
8. After successfull installation you should configure `openvpn` options in `config`.
9. Restart or manually run `sudo installer run`

#### Remove
1. Run command shell (`terminal`).
2. Go to `bin` directory.
3. Execute command `sudo installer remove`
4. Wait for removing process to end. All steps will be printed to console and written to `log`.

#### Help
1. Run command shell (`terminal`).
2. Go to `bin` directory.
3. Execute command `sudo installer help` to display help information.
