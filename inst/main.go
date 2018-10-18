// +build windows

package main

import (
	"fmt"
	"os"

	"github.com/privatix/dapp-openvpn/inst/command"
	"github.com/privatix/dappctrl/util/log"
)

func main() {
	logConfig := &log.FileConfig{
		WriterConfig: log.NewWriterConfig(),
		Filename:     "../log/installer-%Y-%m-%d.log",
		FileMode:     0644,
	}

	logger, closer, err := log.NewFileLogger(logConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %s", err))
	}
	defer closer.Close()

	command.Execute(logger, os.Args[1:])
}
