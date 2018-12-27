// +build !nomsgtest

package msg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/privatix/dappctrl/util"
)

const (
	password               = "secret"
	serviceEndpointAddress = "example.com"

	caDataKey = "caData"
	remoteKey = "remote"
	portKey   = "port"

	managementPortKey = "management"
	tapInterfaceKey   = "dev-node"
	logAppendKey      = "log-append"
)

var (
	username = util.NewUUID()
)

func readStatikFile(t *testing.T, name string) []byte {
	data, err := readFileFromVirtualFS(name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func copyStrStrMap(params map[string]string) (dst map[string]string) {
	dst = make(map[string]string)

	for k, v := range params {
		dst[k] = v
	}
	return dst
}

func testAdditionalParams(t *testing.T, parameters map[string]string) map[string]string {
	ca := readStatikFile(t, sampleCa)

	result := copyStrStrMap(parameters)
	result[caDataKey] = string(ca)
	return result
}

func checkAccess(t *testing.T, file, username, password string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	expected := []byte(fmt.Sprintf("%s\n%s\n", username, password))

	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("expected %s, got %s", expected, data)
	}
}

func checkCA(t *testing.T, config string, ca []byte) {
	configData, err := ioutil.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}

	start := strings.Index(string(configData), `<ca>`)
	stop := strings.LastIndex(string(configData), `</ca>`)

	data := configData[start+5 : stop]

	if !reflect.DeepEqual(ca, data) {
		t.Fatalf("expected %s, got %s", ca, data)
	}
}

func checkConf(t *testing.T, config string,
	keys []string, options map[string]interface{}) {
	specialKeys := []string{remoteKey, tapInterfaceKey,
		managementPortKey, logAppendKey}
	keys = append(keys, specialKeys...)

	result, err := vpnParams(logger, config, keys)
	if err != nil {
		t.Fatal(err)
	}

	// Checks special argument `dev-node`.
	gotTapIf, ok := result[tapInterfaceKey]
	if !ok {
		t.Fatal(`special argument "dev-node" not exists`)
	}

	tapInterfaceVal := options[TapInterface]
	if runtime.GOOS == "windows" {
		tapInterfaceVal = fmt.Sprintf(`"%s"`, options[TapInterface])
	}

	if tapInterfaceVal != gotTapIf {
		t.Fatalf("special argument dev-node not exists,"+
			" wanted: %s, got: %s", tapInterfaceVal, gotTapIf)
	}

	// Checks special argument `log-append`.
	gotLogAppend, ok := result[logAppendKey]
	if !ok {
		t.Fatal(`special argument "log-append" not exists`)
	}

	var expLogAppend string

	logDir, ok := options[LogDir]
	if ok {
		if dir, ok := logDir.(string); ok {
			logName := fmt.Sprintf("openvpn-%s.log", username)
			openVpnLog := filepath.Join(dir, logName)

			if runtime.GOOS == "windows" {
				str := strings.Replace(openVpnLog,
					`\`, `\\`, -1)
				expLogAppend = fmt.Sprintf(
					`"%s"`, str)
			} else {
				expLogAppend = openVpnLog
			}
		}
	}

	if expLogAppend != gotLogAppend {
		t.Fatalf("special argument log-append not exists,"+
			" wanted: %s, got: %s", expLogAppend, gotLogAppend)
	}

	// Checks special argument `management`.
	manPortSrt, ok := result[managementPortKey]
	if !ok {
		t.Fatal(`special argument "management" not exists`)
	}

	manPortSrtWords := strings.Split(manPortSrt, " ")
	if len(manPortSrtWords) != 2 {
		t.Fatal(`special argument "management" not exists`)
	}

	expManPort := fmt.Sprintf("%d", options[VpnManagementPort])
	if expManPort != manPortSrtWords[1] {
		t.Fatalf("special argument management not exists,"+
			" wanted: %s, got: %s", expManPort, manPortSrtWords)
	}

	// Checks special argument "remote".
	val, ok := result[remoteKey]
	if !ok {
		t.Fatal(`special argument "remote" not exists`)
	}

	if !strings.HasSuffix(val, conf.TestVPNConfig[portKey]) {
		t.Fatal(`special argument "port" not exists`)
	}

	// Clears the map of special parameters.
	for _, v := range specialKeys {
		delete(result, v)
	}

	// Adds the just-tested parameter to the resulting map.
	result[portKey] = conf.TestVPNConfig[portKey]

	// Checks special parameter "proto".
	if result[paramProto] != defaultProto {
		t.Fatal(`special argument "proto" must be "tcp-client"`)
	}

	// Changes the just-tested parameter to the resulting map.
	// On the client "tcp" parameter is replaced by "tcp-client"
	result[paramProto] = conf.TestVPNConfig[paramProto]

	if !reflect.DeepEqual(conf.TestVPNConfig, result) {
		t.Fatal("result parameters not equals initial parameters")
	}
}

func TestMakeFiles(t *testing.T) {
	rootDir, err := ioutil.TempDir("", util.NewUUID())
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(rootDir)

	params := testAdditionalParams(t, conf.TestVPNConfig)

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	accessFile := filepath.Join(rootDir, defaultAccessFile)
	confFile := filepath.Join(rootDir, clientConfigFile)

	options := map[string]interface{}{
		VpnManagementPort: uint16(7777),
		TapInterface:      "12345",
		LogDir:            rootDir,
	}

	if err := MakeFiles(logger, rootDir, serviceEndpointAddress, username,
		password, data, options); err != nil {
		t.Fatal(err)
	}

	checkAccess(t, accessFile, username, password)
	checkCA(t, confFile, []byte(params[caDataKey]))
	checkConf(t, confFile, parameterKeys(conf.TestVPNConfig), options)
}
