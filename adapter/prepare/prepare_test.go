// +build !nopreparetest

package prepare

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/privatix/dappctrl/data"
	"github.com/privatix/dappctrl/svc/connector"
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
		StderrLog  *log.WriterConfig
		VPNMonitor *mon.Config
	}

	logger log.Logger
	conn   *connector.Mock
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
	adapterConfig.Connector.Username = channel
	adapterConfig.Connector.Password = data.TestPassword
	adapterConfig.OpenVPN.ConfigRoot = rootDir
	adapterConfig.Monitor = conf.VPNMonitor
	adapterConfig.FileLog.WriterConfig = conf.StderrLog

	endpoint := "127.0.0.1"

	conn.Endpoint = data.NewTestEndpoint(channel, util.NewUUID())
	conn.Endpoint.ServiceEndpointAddress = &endpoint

	if err := ClientConfig(
		logger, channel, conn, adapterConfig); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(rootDir, channel)

	checkFile(t, configDestination(target))
	checkFile(t, accessDestination(target))
}

func TestMain(m *testing.M) {
	var err error

	conf.StderrLog = log.NewWriterConfig()
	conf.VPNMonitor = mon.NewConfig()

	args := &util.TestArgs{
		Conf: &conf,
	}
	util.ReadTestArgs(args)

	logger, err = log.NewStderrLogger(conf.StderrLog)
	if err != nil {
		panic(err)
	}

	conn = connector.NewMock()

	os.Exit(m.Run())
}
