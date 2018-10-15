package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/privatix/dapp-openvpn/inst/openvpn"
	"github.com/privatix/dappctrl/util/log"
)

// Execute executes a CLI command.
func Execute(logger log.Logger, args []string) {
	if len(args) == 0 {
		args = append(args, "help")
	}

	var flow openvpn.Flow

	switch strings.ToLower(args[0]) {
	case "install":
		flow = installFlow()
	case "remove":
		flow = removeFlow()
	default:
		fmt.Println(rootHelp)
		return
	}

	ovpn := openvpn.NewOpenVPN()
	if err := flow.Run(ovpn, logger); err != nil {
		logger.Error(fmt.Sprintf("failed to execute command: %v", err))
		os.Exit(2)
	}

	logger.Info("command was successfully executed")
}
