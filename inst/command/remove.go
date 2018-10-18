package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/privatix/dapp-openvpn/inst/openvpn"
	"github.com/privatix/dapp-openvpn/inst/pipeline"
)

func removeFlow() pipeline.Flow {
	return pipeline.Flow{
		newOperator("processed flags", processedRemoveFlags, nil),
		newOperator("validate", validateToRemove, nil),
		newOperator("stop service", stopService, nil),
		newOperator("remove tap", removeTap, nil),
		newOperator("remove service", removeService, nil),
		newOperator("remove env", removeEnv, nil),
		//newOperator("remove folder", removeFolder, nil),
	}
}

func processedRemoveFlags(ovpn *openvpn.OpenVPN) error {
	h := flag.Bool("help", false, "Display installer help")
	p := flag.String("workdir", ".", "Product install directory")

	flag.CommandLine.Parse(os.Args[2:])

	if *h {
		fmt.Println(installHelp)
		os.Exit(0)
	}

	ovpn.Path = *p
	return nil
}

func validateToRemove(o *openvpn.OpenVPN) error {
	path, err := filepath.Abs(o.Path)
	if err != nil {
		return err
	}
	o.Path = filepath.ToSlash(strings.ToLower(path))

	err = godotenv.Load(filepath.Join(o.Path, "config/.env"))
	if err != nil {
		return err
	}

	w := os.Getenv("WORKDIR")
	if !strings.EqualFold(o.Path, w) {
		return fmt.Errorf("env workdir %s is not equal to the path", w)
	}
	o.Tap.DeviceID = os.Getenv("DEVICE")
	o.Tap.Interface = os.Getenv("INTERFACE")
	o.Service = os.Getenv("SERVICE")

	return nil
}

func stopService(o *openvpn.OpenVPN) error {
	if err := o.StopService(); err != nil {
		return fmt.Errorf("failed to stop service: %v", err)
	}
	return nil
}

func removeTap(o *openvpn.OpenVPN) error {
	if err := o.RemoveTap(); err != nil {
		return fmt.Errorf("failed to remove tap interface: %v", err)
	}
	return nil
}

func removeService(o *openvpn.OpenVPN) error {
	if err := o.RemoveService(); err != nil {
		return fmt.Errorf("failed to remove service: %v", err)
	}
	return nil
}

func removeEnv(o *openvpn.OpenVPN) error {
	return os.Remove(filepath.Join(o.Path, "config/.env"))
}

func removeFolder(o *openvpn.OpenVPN) error {
	return os.RemoveAll(o.Path)
}
