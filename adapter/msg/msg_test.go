// +build !nomsgtest

package msg

import (
	"os"
	"testing"

	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/adapter/mon"
	"github.com/privatix/dapp-openvpn/adapter/util"
)

var (
	conf struct {
		Pusher        *Config
		Monitor       *mon.Config
		TestVPNConfig map[string]string
	}

	logger log.Logger
)

func newTestVPNConfig() map[string]string {
	return make(map[string]string)
}

func TestMain(m *testing.M) {
	conf.Pusher = NewConfig()
	conf.Monitor = mon.NewConfig()
	conf.TestVPNConfig = newTestVPNConfig()

	util.ReadTestConfig(&conf)

	logger = log.NewMultiLogger()

	os.Exit(m.Run())
}
