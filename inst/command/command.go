package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/privatix/dapp-openvpn/inst/openvpn"
	"github.com/privatix/dapp-openvpn/inst/pipeline"
	"github.com/privatix/dappctrl/util/log"
)

// Execute executes a CLI command.
func Execute(logger log.Logger, args []string) {
	if len(args) == 0 {
		args = append(args, "help")
	}

	var flow pipeline.Flow

	switch strings.ToLower(args[0]) {
	case "install":
		logger.Info("start install process")
		flow = installFlow()
	case "remove":
		logger.Info("start remove process")
		flow = removeFlow()
	default:
		fmt.Println(rootHelp)
		return
	}

	ovpn := openvpn.NewOpenVPN()
	if err := flow.Run(ovpn, logger); err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		fmt.Println("failed to execute command")
		os.Exit(2)
	}

	fmt.Println("command was successfully executed")
	logger.Info("command was successfully executed")
}
