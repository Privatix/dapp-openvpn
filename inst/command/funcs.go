package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/privatix/dappctrl/util"

	"github.com/privatix/dapp-openvpn/inst/env"
	"github.com/privatix/dapp-openvpn/inst/openvpn"
)

func processedRootFlags(printVersion func()) {
	v := flag.Bool("version", false, "Prints current inst version")

	flag.Parse()

	if *v {
		printVersion()
		os.Exit(0)
	}

	fmt.Printf(rootHelp)
}

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

	if s, err := o.Adapter.InstallService(o.Role, o.Path); err != nil {
		return fmt.Errorf("failed to install service: %v %v", s, err)
	}

	if !o.IsWindows {
		return nil
	}

	name := strings.Replace(o.Adapter.Service, " ", "_", -1)
	cmd := exec.Command("sc", "failure", name, "reset=", "0",
		"actions=", "restart/1000/restart/2000/restart/5000")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set service restart rule: %v", err)
	}

	return nil
}

func startService(o *openvpn.OpenVPN) error {
	if s, err := o.StartService(); err != nil {
		return fmt.Errorf("failed to start service: %v %v", s, err)
	}

	if s, err := o.Adapter.StartService(); err != nil {
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

func runAdapter(o *openvpn.OpenVPN) error {
	if s, err := o.Adapter.RunService(o.Role, o.Path); err != nil {
		return fmt.Errorf("failed to run service: %v %v", s, err)
	}
	return nil
}

func stopService(o *openvpn.OpenVPN) error {
	if s, err := o.Adapter.StopService(); err != nil {
		return fmt.Errorf("failed to stop adapter service: %v %v",
			s, err)
	}

	if s, err := o.StopService(); err != nil {
		return fmt.Errorf("failed to stop service: %v %v", s, err)
	}
	return nil
}

func removeService(o *openvpn.OpenVPN) error {
	if s, err := o.Adapter.RemoveService(); err != nil {
		return fmt.Errorf("failed to remove adapter service: %v %v",
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
	role := flag.String("role", "", "Product role")
	p := flag.String("workdir", "", "Product install directory")

	flag.CommandLine.Parse(os.Args[2:])

	if *h {
		fmt.Println(installHelp)
		os.Exit(0)
	}

	if len(*config) > 0 {
		if err := util.ReadJSONFile(*config, &ovpn); err != nil {
			return err
		}
	}

	if len(*p) > 0 {
		ovpn.Path = *p
	}
	if len(*role) > 0 {
		ovpn.Role = *role
	}
	return nil
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
	o.Import = v.ProductImport

	if v.ProductInstall || strings.EqualFold(o.Path, v.Workdir) {
		err = errors.New("openvpn was installed at this workdir")
	}
	return err
}

func createConfig(o *openvpn.OpenVPN) error {
	if err := o.Configurate(); err != nil {
		return fmt.Errorf("failed to configure openvpn: %v", err)
	}

	if err := o.Adapter.Configurate(o); err != nil {
		return fmt.Errorf("failed to configure adapter: %v", err)
	}
	return nil
}

func removeConfig(o *openvpn.OpenVPN) error {
	return o.RemoveConfig()
}

func recordPortsToOpen(o *openvpn.OpenVPN) error {
	portsF, err := os.Create(filepath.Join(o.Path, "config/ports.txt"))
	if err != nil {
		return err
	}
	defer portsF.Close()

	_, err = fmt.Fprintf(portsF, "%d(%s)", o.Host.Port, o.Proto)
	return err
}

func removePortsFile(o *openvpn.OpenVPN) error {
	return os.Remove(filepath.Join(o.Path, "config/ports.txt"))
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

	o.Tap.DeviceID = v.Device
	o.Tap.GUID = v.GUID
	o.Tap.Interface = v.Interface
	o.Service = v.Service
	o.Role = v.Role
	o.Adapter.Service = v.Adapter
	o.Import = v.ProductImport
	o.Install = v.ProductInstall
	o.ForwardingState = v.ForwardingState

	return nil
}

func createEnv(o *openvpn.OpenVPN) error {
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("/usr/sbin/sysctl", "-n",
			"net.inet.ip.forwarding").Output()
		if err != nil {
			return err
		}
		o.ForwardingState = strings.Replace(string(out), "\n", "", -1)
	}

	v := env.NewConfig()

	v.Workdir = o.Path
	v.Device = o.Tap.DeviceID
	v.GUID = o.Tap.GUID
	v.Interface = o.Tap.Interface
	v.Service = o.Service
	v.Role = o.Role
	v.Adapter = o.Adapter.Service
	v.ProductImport = o.Import
	v.ProductInstall = o.Install
	v.ForwardingState = o.ForwardingState

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

func finalize(o *openvpn.OpenVPN) error {
	if !strings.EqualFold(o.Role, "server") {
		return nil
	}

	if !o.IsWindows {
		return o.CreateForwardingDaemon()
	}

	if err := o.CheckServiceStatus("running"); err != nil {
		return err
	}

	if err := stopService(o); err != nil {
		return fmt.Errorf("failed to stop services: %v", err)
	}

	if err := o.CheckServiceStatus("stopped"); err != nil {
		return err
	}

	if err := startServices(o); err != nil {
		return fmt.Errorf("failed to start services: %v", err)
	}
	return nil
}

func changeOwner(o *openvpn.OpenVPN) error {
	if runtime.GOOS != "darwin" {
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

	file := filepath.Join(o.Path, envFile)
	if err := os.Chown(file, uid, gid); err != nil {
		return fmt.Errorf("failed to change owner: %v", err)
	}

	return nil
}

func update(o *openvpn.OpenVPN) error {
	if runtime.GOOS == "linux" {
		// doesn't need stop/start services,
		// because container is used on linux
		return o.Update()
	}

	if err := stopService(o); err != nil {
		return err
	}

	if err := o.Update(); err != nil {
		if err := startService(o); err != nil {
			return err
		}
		return fmt.Errorf("failed to update product: %v", err)
	}

	return nil
}
