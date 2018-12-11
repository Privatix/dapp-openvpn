// +build !nopreparetest

package prepare

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/privatix/dappctrl/data"
	"github.com/privatix/dappctrl/util"
	"github.com/privatix/dappctrl/util/log"

	"github.com/privatix/dapp-openvpn/adapter/config"
	"github.com/privatix/dapp-openvpn/adapter/mon"
)

const (
	accessFileName = "access.ovpn"
	configFileName = "client.ovpn"
)

var (
	conf struct {
		VPNMonitor *mon.Config
	}

	logger log.Logger
)

func configDestination(dir string) string {
	return filepath.Join(dir, configFileName)
}

func accessDestination(dir string) string {
	return filepath.Join(dir, accessFileName)
}

func checkFile(t *testing.T, file string) {
	stat, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Fatal("file is empty")
	}
}

func TestClientConfig(t *testing.T) {
	rootDir, err := ioutil.TempDir("", util.NewUUID())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(rootDir)

	channel := util.NewUUID()

	adapterConfig := config.NewConfig()
	adapterConfig.Sess.Product = channel
	adapterConfig.Sess.Password = data.TestPassword
	adapterConfig.OpenVPN.ConfigRoot = rootDir
	adapterConfig.Monitor = conf.VPNMonitor

	ept := data.NewTestEndpoint(channel, util.NewUUID())
	ept.ServiceEndpointAddress = pointer.ToString("1.2.3.4")
	getEndpoint := func(clientKey string) (*data.Endpoint, error) {
		return ept, nil
	}

	err = ClientConfig(logger, channel, adapterConfig, getEndpoint)
	if err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(rootDir, channel)

	checkFile(t, configDestination(target))
	checkFile(t, accessDestination(target))
}

func TestMain(m *testing.M) {
	conf.VPNMonitor = mon.NewConfig()

	util.ReadTestConfig(&conf)

	logger = log.NewMultiLogger()

	os.Exit(m.Run())
}
