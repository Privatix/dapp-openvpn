package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/inst/command"
)

func createLogger() (log.Logger, io.Closer, error) {
	elog, err := log.NewStderrLogger(log.NewWriterConfig())
	if err != nil {
		return nil, nil, err
	}

	f := flag.NewFlagSet("", flag.ContinueOnError)
	p := f.String("workdir", "..", "Product install directory")

	if len(os.Args) > 2 && !strings.EqualFold(os.Args[1], "install") {
		f.Parse(os.Args[2:])
	}

	if strings.EqualFold(*p, "..") {
		*p = filepath.Join(filepath.Dir(os.Args[0]), *p)
	}

	path, _ := filepath.Abs(*p)
	path = filepath.ToSlash(path)

	fileName := filepath.Join(path, "log/installer-%Y-%m-%d.log")

	logConfig := &log.FileConfig{
		WriterConfig: log.NewWriterConfig(),
		Filename:     fileName,
		FileMode:     0644,
	}

	flog, closer, err := log.NewFileLogger(logConfig)
	if err != nil {
		return nil, nil, err
	}

	logger := log.NewMultiLogger(elog, flog)

	return logger, closer, nil
}

func main() {
	logger, closer, err := createLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %s", err))
	}
	defer closer.Close()

	command.Execute(logger, os.Args[1:])
}
