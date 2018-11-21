package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/privatix/dappctrl/util"

	"github.com/privatix/dapp-openvpn/inst/env"
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

	if s, err := o.DappVPN.InstallService(o.Role, o.Path); err != nil {
		return fmt.Errorf("failed to install service: %v %v", s, err)
	}
	return nil
}

func startService(o *openvpn.OpenVPN) error {
	if s, err := o.StartService(); err != nil {
		return fmt.Errorf("failed to start service: %v %v", s, err)
	}

	if s, err := o.DappVPN.StartService(); err != nil {
		return fmt.Errorf("failed to start service: %v %v", s, err)
	}
	return nil
}

func runService(o *openvpn.OpenVPN) error {
	if s, err := o.RunService(); err != nil {
		return fmt.Errorf("failed to run service: %v %v", s, err)
	}
	return nil
}

func runDappVPN(o *openvpn.OpenVPN) error {
	if s, err := o.DappVPN.RunService(o.Role, o.Path); err != nil {
		return fmt.Errorf("failed to run service: %v %v", s, err)
	}
	return nil
}

func stopService(o *openvpn.OpenVPN) error {
	if s, err := o.DappVPN.StopService(); err != nil {
		return fmt.Errorf("failed to stop dappvpn service: %v %v",
			s, err)
	}

	if s, err := o.StopService(); err != nil {
		return fmt.Errorf("failed to stop service: %v %v", s, err)
	}
	return nil
}

func removeService(o *openvpn.OpenVPN) error {
	if s, err := o.DappVPN.RemoveService(); err != nil {
		return fmt.Errorf("failed to remove dappvpn service: %v %v",
			s, err)
	}

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
	if strings.EqualFold(o.Path, "..") {
		o.Path = filepath.Join(filepath.Dir(os.Args[0]), "..")
	}
	path, err := filepath.Abs(o.Path)
	if err != nil {
		return err
	}
	o.Path = filepath.ToSlash(path)
	return nil
}

func validateToInstall(o *openvpn.OpenVPN) error {
	err := validatePath(o)
	if err != nil {
		return err
	}

	v := env.NewConfig()
	// When installing the environment file may not be.
	// It is created on the installation finalize.
	_ = v.Read(filepath.Join(o.Path, envFile))

	if strings.EqualFold(o.Path, v.Workdir) {
		err = errors.New("openvpn was installed at this workdir")
	}
	return err
}

func createConfig(o *openvpn.OpenVPN) error {
	if err := o.Configurate(); err != nil {
		return fmt.Errorf("failed to configure openvpn: %v", err)
	}

	if err := o.DappVPN.Configurate(o); err != nil {
		return fmt.Errorf("failed to configure dappvpn: %v", err)
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

	v := env.NewConfig()
	if err := v.Read(filepath.Join(o.Path, envFile)); err != nil {
		return err
	}

	if !strings.EqualFold(o.Path, v.Workdir) {
		return fmt.Errorf("env workdir %s is not equal to the path",
			v.Workdir)
	}
	o.Tap.DeviceID = v.Device
	o.Tap.Interface = v.Interface
	o.Service = v.Service
	o.Role = v.Role
	o.DappVPN.Service = v.DappVPN
	o.Import = v.ProductImport
	o.Install = v.ProductInstall

	return nil
}

func createEnv(o *openvpn.OpenVPN) error {
	v := env.NewConfig()

	v.Workdir = o.Path
	v.Device = o.Tap.DeviceID
	v.Interface = o.Tap.Interface
	v.Service = o.Service
	v.Role = o.Role
	v.DappVPN = o.DappVPN.Service
	v.ProductImport = o.Import
	v.ProductInstall = o.Install

	if err := v.Write(filepath.Join(o.Path, envFile)); err != nil {
		return fmt.Errorf("failed to create env file: %v", err)
	}

	return nil
}

func removeEnv(o *openvpn.OpenVPN) error {
	if !o.Import {
		return os.Remove(filepath.Join(o.Path, envFile))
	}

	v := env.NewConfig()

	v.ProductImport = o.Import

	if err := v.Write(filepath.Join(o.Path, envFile)); err != nil {
		return fmt.Errorf("failed to write env file: %v", err)
	}

	return nil
}

func startServices(o *openvpn.OpenVPN) error {
	cmd := exec.Command(os.Args[0], "start")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start services: %v", err)
	}
	return nil
}

func changeOwner(o *openvpn.OpenVPN) error {
	if o.IsWindows {
		return nil
	}

	logname, err := exec.Command("logname").Output()
	if err != nil {
		return err
	}

	group, err := user.Lookup(strings.TrimSpace(string(logname)))
	if err != nil {
		return fmt.Errorf("failed to lookup user: %v", err)
	}
	uid, _ := strconv.Atoi(group.Uid)
	gid, _ := strconv.Atoi(group.Gid)

	if err := filepath.Walk(o.Path,
		func(name string, info os.FileInfo, err error) error {
			if err == nil {
				err = os.Chown(name, uid, gid)
			}
			return err
		}); err != nil {
		return fmt.Errorf("failed to change owner: %v", err)
	}
	return nil
}
