// +build !nomsgtest

package msg

import (
	"os"
	"testing"

	"github.com/privatix/dappctrl/svc/connector"
	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/adapter/mon"
	"github.com/privatix/dapp-openvpn/adapter/util"
)

var (
	conf struct {
		FileLog       *log.FileConfig
		Pusher        *Config
		VPNMonitor    *mon.Config
		TestVPNConfig map[string]string
	}

	logger log.Logger
	conn   *connector.Mock
)

func newTestVPNConfig() map[string]string {
	return make(map[string]string)
}

func TestMain(m *testing.M) {
	var err error

	conf.FileLog = log.NewFileConfig()
	conf.Pusher = NewConfig()
	conf.VPNMonitor = mon.NewConfig()
	conf.TestVPNConfig = newTestVPNConfig()

	util.ReadTestConfig(&conf)

	logger, err = log.NewStderrLogger(conf.FileLog)
	if err != nil {
		panic(err)
	}

	conn = connector.NewMock()

	os.Exit(m.Run())
}
