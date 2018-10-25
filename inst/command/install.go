package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/privatix/dappctrl/util"

	"github.com/privatix/dapp-openvpn/inst/openvpn"
	"github.com/privatix/dapp-openvpn/inst/pipeline"
)

func installFlow() pipeline.Flow {
	return pipeline.Flow{
		newOperator("processed flags", processedInstallFlags, nil),
		newOperator("validate", validateToInstall, nil),
		newOperator("install tap", installTap, removeTap),
		newOperator("configuration", configurate, removeConfig),
		newOperator("registration", registerService, removeService),
		newOperator("create env", createEnv, removeEnv),
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

func validatePath(o *openvpn.OpenVPN) error {
	path, err := filepath.Abs(o.Path)
	if err != nil {
		return err
	}
	o.Path = filepath.ToSlash(strings.ToLower(path))
	return nil
}

func validateToInstall(o *openvpn.OpenVPN) error {
	err := validatePath(o)
	if err != nil {
		return err
	}

	// When installing the environment file may not be.
	// It is created on the installation finalize.
	_ = godotenv.Load(filepath.Join(o.Path, envFile))

	if strings.EqualFold(o.Path, os.Getenv(envWorkDir)) {
		err = errors.New("openvpn was installed at this workdir")
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

func removeConfig(o *openvpn.OpenVPN) error {
	return o.RemoveConfig()
}

func registerService(o *openvpn.OpenVPN) error {
	if err := o.RegisterService(); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}
	return nil
}

func createEnv(o *openvpn.OpenVPN) error {
	env := make(map[string]string)

	env[envWorkDir] = o.Path
	env[envDevice] = o.Tap.DeviceID
	env[envInterface] = o.Tap.Interface
	env[envServcie] = o.Service

	err := godotenv.Write(env, filepath.Join(o.Path, envFile))
	if err != nil {
		return fmt.Errorf("failed to create env file: %v", err)
	}

	return nil
}
