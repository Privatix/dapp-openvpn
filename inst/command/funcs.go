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
)

func installTap(o *openvpn.OpenVPN) error {
	if err := o.InstallTap(); err != nil {
		return fmt.Errorf("failed to install tap interface: %v", err)
	}
	return nil
}

func removeTap(o *openvpn.OpenVPN) error {
	if err := o.RemoveTap(); err != nil {
		return fmt.Errorf("failed to remove tap interface: %v", err)
	}
	return nil
}

func createService(o *openvpn.OpenVPN) error {
	if s, err := o.InstallService(); err != nil {
		return fmt.Errorf("failed to install service: %v %v", s, err)
	}
	return nil
}

func startService(o *openvpn.OpenVPN) error {
	if s, err := o.StartService(); err != nil {
		return fmt.Errorf("failed to stop service: %v %v", s, err)
	}
	return nil
}

func runService(o *openvpn.OpenVPN) error {
	s, err := o.RunService()
	if err != nil {
		return fmt.Errorf("failed to run service: %v %v", s, err)
	}
	if len(s) > 0 {
		return fmt.Errorf(s)
	}
	return nil
}

func stopService(o *openvpn.OpenVPN) error {
	if s, err := o.StopService(); err != nil {
		return fmt.Errorf("failed to stop service: %v %v", s, err)
	}
	return nil
}

func removeService(o *openvpn.OpenVPN) error {
	if s, err := o.RemoveService(); err != nil {
		return fmt.Errorf("failed to remove service: %v %v", s, err)
	}
	return nil
}

func processedInstallFlags(ovpn *openvpn.OpenVPN) error {
	h := flag.Bool("help", false, "Display installer help")
	config := flag.String("config", "", "Configuration file")

	flag.CommandLine.Parse(os.Args[2:])

	if *h || len(*config) == 0 {
		fmt.Println(installHelp)
		os.Exit(0)
	}

	return util.ReadJSONFile(*config, &ovpn)
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

func createConfig(o *openvpn.OpenVPN) error {
	if err := o.Configurate(); err != nil {
		return fmt.Errorf("failed to configure openvpn: %v", err)
	}
	return nil
}

func removeConfig(o *openvpn.OpenVPN) error {
	return o.RemoveConfig()
}

func processedCommonFlags(ovpn *openvpn.OpenVPN) error {
	h := flag.Bool("help", false, "Display installer help")
	p := flag.String("workdir", "..", "Product install directory")

	flag.CommandLine.Parse(os.Args[2:])

	if *h {
		fmt.Printf(templateHelp, os.Args[1])
		os.Exit(0)
	}

	ovpn.Path = *p
	return nil
}

func checkInstallation(o *openvpn.OpenVPN) error {
	err := validatePath(o)
	if err != nil {
		return err
	}

	err = godotenv.Load(filepath.Join(o.Path, envFile))
	if err != nil {
		return err
	}

	w := os.Getenv(envWorkDir)
	if !strings.EqualFold(o.Path, w) {
		return fmt.Errorf("env workdir %s is not equal to the path", w)
	}
	o.Tap.DeviceID = os.Getenv(envDevice)
	o.Tap.Interface = os.Getenv(envInterface)
	o.Service = os.Getenv(envServcie)

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

func removeEnv(o *openvpn.OpenVPN) error {
	return os.Remove(filepath.Join(o.Path, envFile))
}
