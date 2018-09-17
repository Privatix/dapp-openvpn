package util

import (
	"flag"

	"github.com/privatix/dappctrl/util"
)

// ReadTestConfig reads test configuration.
func ReadTestConfig(config interface{}) {
	confFile := flag.String("config", "test.conf", "Test configuration")

	flag.Parse()

	if err := util.ReadJSONFile(*confFile,
		&config); err != nil {
		panic("failed to read configuration: " + err.Error())
	}
}
