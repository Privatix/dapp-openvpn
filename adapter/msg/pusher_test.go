// +build !nomsgtest

package msg

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/privatix/dappctrl/util"

	"github.com/privatix/dapp-openvpn/statik"
)

const (
	caFileName   = "ca.crt"
	ovpnFileName = "server.ovpn"
	sampleCa     = "/ovpn/samples/ca.crt"
	sampleConf   = "/ovpn/samples/server.ovpn"
)

func readFile(t *testing.T, name string) []byte {
	file, err := statik.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func createTestConfig(t *testing.T, dir string) *Config {
	cfgData := readFile(t, sampleConf)
	caData := readFile(t, sampleCa)

	cfgPath := filepath.Join(dir, ovpnFileName)
	caPath := filepath.Join(dir, caFileName)

	if err := ioutil.WriteFile(cfgPath, cfgData, filePerm); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(caPath, caData, filePerm); err != nil {
		t.Fatal(err)
	}

	return &Config{
		ExportConfigKeys: conf.Pusher.ExportConfigKeys,
		ConfigPath:       cfgPath,
		CaCertPath:       caPath,
		TimeOut:          conf.Pusher.TimeOut,
	}
}

func TestPushConfig(t *testing.T) {
	rootDir, err := ioutil.TempDir("", util.NewUUID())
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(rootDir)

	setProductConfig := func(config map[string]string) error {
		return nil
	}

	pusher := NewPusher(
		createTestConfig(t, rootDir), logger, setProductConfig)
	if err := pusher.PushConfiguration(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestConfigPushedFile(t *testing.T) {
	rootDir, err := ioutil.TempDir("", util.NewUUID())
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(rootDir)

	if IsDone(rootDir) {
		t.Fatal("configuration not yet updated")
	}
	if err := Done(rootDir); err != nil {
		t.Fatal(err)
	}
	if !IsDone(rootDir) {
		t.Fatal("configuration already updated")
	}
}
