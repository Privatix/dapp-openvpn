package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/privatix/dapp-openvpn/inst/openvpn"
)

func removeFlow() openvpn.Flow {
	return openvpn.Flow{
		openvpn.NewOperator(processedRemoveFlags, nil),
		openvpn.NewOperator(validateToRemove, nil),
		openvpn.NewOperator(stopService, nil),
		openvpn.NewOperator(removeTap, nil),
		openvpn.NewOperator(removeService, nil),
		openvpn.NewOperator(removeFolder, nil),
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

func removeFolder(o *openvpn.OpenVPN) error {
	return os.RemoveAll(o.Path)
}
