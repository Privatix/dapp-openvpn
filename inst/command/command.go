package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/inst/openvpn"
	"github.com/privatix/dapp-openvpn/inst/pipeline"
)

// Execute executes a CLI command.
func Execute(logger log.Logger, args []string) {
	if len(args) == 0 {
		args = append(args, "help")
	}

	var flow pipeline.Flow

	switch strings.ToLower(args[0]) {
	case "install":
		logger.Info("install process")
		flow = installFlow()
	case "remove":
		logger.Info("remove process")
		flow = removeFlow()
	case "start":
		logger.Info("start process")
		flow = startFlow()
	case "stop":
		logger.Info("stop process")
		flow = stopFlow()
	case "run":
		logger.Info("run process")
		flow = runFlow()
	case "run-adapter":
		logger.Info("run adapter process")
		flow = runAdapterFlow()
	default:
		fmt.Println(rootHelp)
		return
	}

	ovpn := openvpn.NewOpenVPN()
	if err := flow.Run(ovpn, logger); err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		os.Exit(2)
	}

	logger.Info("command was successfully executed")
}

func installFlow() pipeline.Flow {
	return pipeline.Flow{
		newOperator("processed flags", processedInstallFlags, nil),
		newOperator("validate", validateToInstall, nil),
		newOperator("install tap", installTap, removeTap),
		newOperator("create config", createConfig, removeConfig),
		newOperator("create service", createService, removeService),
		newOperator("create env", createEnv, removeEnv),
		newOperator("start services", startServices, nil),
		newOperator("change owner", changeOwner, nil),
	}
}

func removeFlow() pipeline.Flow {
	return pipeline.Flow{
		newOperator("processed flags", processedCommonFlags, nil),
		newOperator("validate", checkInstallation, nil),
		newOperator("stop service", stopService, nil),
		newOperator("remove tap", removeTap, nil),
		newOperator("remove service", removeService, nil),
		newOperator("remove env", removeEnv, nil),
	}
}

func startFlow() pipeline.Flow {
	return pipeline.Flow{
		newOperator("processed flags", processedCommonFlags, nil),
		newOperator("validate", checkInstallation, nil),
		newOperator("start service", startService, nil),
	}
}

func stopFlow() pipeline.Flow {
	return pipeline.Flow{
		newOperator("processed flags", processedCommonFlags, nil),
		newOperator("validate", checkInstallation, nil),
		newOperator("stop service", stopService, nil),
	}
}

func runFlow() pipeline.Flow {
	return pipeline.Flow{
		newOperator("processed flags", processedCommonFlags, nil),
		newOperator("validate", checkInstallation, nil),
		newOperator("run service", runService, nil),
	}
}

func runAdapterFlow() pipeline.Flow {
	return pipeline.Flow{
		newOperator("processed flags", processedCommonFlags, nil),
		newOperator("validate", checkInstallation, nil),
		newOperator("run adapter", runAdapter, nil),
	}
}
