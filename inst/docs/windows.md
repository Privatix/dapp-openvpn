#### Install 
1. Choose `openvpn` version: server (`openvpn-win-server.zip`) or client (`openvpn-win-client.zip`)
2. Extract selected zip archive.
3. Run command shell (`cmd`) with elevate role `Run as administrator`.
4. Go to `bin` directory.
5. Execute command `installer.exe install -config ../config/installer.config.json`
6. Wait to end installation process. All steps print into console and write to `log`.
7. After successfull installation `config/.env` file with following key-values will be created: `DEVICE`, `INTERFACE`, `SERVICE`, `WORKDIR`
8. After successfull installation you should configure `openvpn` options in `config`.
9. Restart or manual run `installer.exe start`

#### Remove
1. Run command shell (`cmd`) with elevate role `Run as administrator`.
2. Go to `bin` directory.
3. Execute command `installer.exe remove`
4. Wait for removing process to end. All steps will be printed into console and written to `log`.

#### Help
1. Run command shell (`cmd`) with elevate role `Run as administrator`.
2. Go to `bin` directory.
3. Execute command `installer.exe help` to display help information.
