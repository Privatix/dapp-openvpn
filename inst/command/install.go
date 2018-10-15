package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/privatix/dapp-openvpn/inst/openvpn"
	"github.com/privatix/dappctrl/util"
)

func installFlow() openvpn.Flow {
	return openvpn.Flow{
		openvpn.NewOperator(processedInstallFlags, nil),
		openvpn.NewOperator(validateToInstall, nil),
		openvpn.NewOperator(installTap, nil),
		openvpn.NewOperator(configurate, nil),
		openvpn.NewOperator(registerService, nil),
	}
}

func processedInstallFlags(ovpn *openvpn.OpenVPN) error {
	h := flag.Bool("help", false, "Display installer help")
	configFile := flag.String("config", "", "Configuration file")

	flag.CommandLine.Parse(os.Args[2:])

	if *h || len(*configFile) == 0 {
		fmt.Println(installHelp)
		os.Exit(0)
	}

	return util.ReadJSONFile(*configFile, &ovpn)
}

func validateToInstall(o *openvpn.OpenVPN) error {
	path, err := filepath.Abs(o.Path)
	if err != nil {
		return err
	}
	o.Path = filepath.ToSlash(strings.ToLower(path))

	deviceID, _ := o.DeviceID()
	if len(deviceID) > 0 {
		err = errors.New("tap was installed at this workdir")
	}
	return err
}

func installTap(o *openvpn.OpenVPN) error {
	if err := o.InstallTap(); err != nil {
		return fmt.Errorf("failed to install tap interface: %v", err)
	}
	return nil
}

func configurate(o *openvpn.OpenVPN) error {
	if err := o.Configurate(); err != nil {
		return fmt.Errorf("failed to configure openvpn: %v", err)
	}
	return nil
}

func registerService(o *openvpn.OpenVPN) error {
	if err := o.RegisterService(); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}
	return nil
}