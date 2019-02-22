package openvpn

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/privatix/dapp-openvpn/inst/openvpn/path"
)

func diff(a, b []string) string {
	for i, v := range b {
		if len(a) <= i || a[i] != v {
			return v
		}
	}
	return ""
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(strings.ToLower(s)))
	return hex.EncodeToString(h.Sum(nil))
}

func nextFreePort(h host, proto string) int {
	hostname := h.IP
	if strings.EqualFold(hostname, "0.0.0.0") {
		hostname = "localhost"
	}
	port := h.Port
	for i := port; i < 65535; i++ {
		ln, err := net.Listen(proto, hostname+":"+strconv.Itoa(i))
		if err != nil {
			continue
		}

		if err := ln.Close(); err != nil {
			continue
		}
		port = i
		break
	}

	return port
}

func getUserGroup() (string, string, error) {
	u, err := user.Current()
	if err != nil {
		return "", "", err
	}

	g, err := user.LookupGroupId(u.Gid)
	if err != nil {
		return u.Username, "", err
	}

	return u.Username, g.Name, nil
}

func setConfigurationValues(jsonMap map[string]interface{},
	settings map[string]interface{}) error {
	for key, value := range settings {
		path := strings.Split(key, ".")
		length := len(path) - 1
		m := jsonMap
		for i := 0; i < length; i++ {
			item, ok := m[path[i]]
			if !ok {
				m[path[i]] = make(map[string]interface{})
				item, ok = m[path[i]]
			}
			if ok && reflect.TypeOf(m) == reflect.TypeOf(item) {
				m, _ = item.(map[string]interface{})
				continue
			}
			return fmt.Errorf("failed to set config params: %s", key)
		}
		m[path[length]] = value
	}
	return nil
}

func sessAddr(config string) (string, error) {
	read, err := os.Open(config)
	if err != nil {
		return "", err
	}
	defer read.Close()

	jsonMap := make(map[string]interface{})

	json.NewDecoder(read).Decode(&jsonMap)

	srv, ok := jsonMap["Sess"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Sess params not found")
	}
	addr, ok := srv["Addr"]
	if !ok {
		return "", fmt.Errorf("Addr params not found")
	}
	return addr.(string), nil
}

func runPowerShellCommand(args ...string) error {
	cmd := exec.Command("powershell", args...)

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	if err := cmd.Run(); err != nil {
		outStr, errStr := outbuf.String(), errbuf.String()
		return fmt.Errorf("%v\nout:\n%s\nerr:\n%s", err, outStr, errStr)
	}
	return nil
}

func buildPowerShellArgs(file string, args ...string) []string {
	a := []string{"-ExecutionPolicy", "Bypass", "-File", file}
	return append(a, args...)
}

func disableNAT(p, device string) error {
	script := filepath.Join(p, path.Config.PowerShellVpnNat)
	args := buildPowerShellArgs(script, "-TAPdeviceAddress", device)
	return runPowerShellCommand(args...)
}

func enableNAT(p, device string) error {
	script := filepath.Join(p, path.Config.PowerShellVpnNat)
	args := buildPowerShellArgs(script,
		"-TAPdeviceAddress", device,
		"-Enabled")
	return runPowerShellCommand(args...)
}

func createScheduleTask(p, device string) error {
	script := filepath.Join(p, path.Config.PowerShellScheduleTask)
	reEnableScript := filepath.Join(p, path.Config.PowerShellReEnableNat)
	args := []string{"-ExecutionPolicy", "Bypass", "-NoProfile",
		"-File", script, "-scriptPath", reEnableScript,
		"-TAPdeviceAddress", device,
	}

	return runPowerShellCommand(args...)
}

func removeScheduleTask() error {
	args := []string{"-ExecutionPolicy", "Bypass", "-NoProfile", "-Command",
		"& {Unregister-ScheduledTask -TaskName 'Privatix re-enable ICS' -confirm:0}",
	}
	return runPowerShellCommand(args...)
}

func copyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())
		if fd.IsDir() {
			err = copyDir(srcfp, dstfp)
		} else {
			err = copyFile(srcfp, dstfp)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	var err error
	if _, err = os.Stat(dst); err == nil {
		return nil
	}

	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()
	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()
	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func merge(src, dst string) error {
	var err error
	var fds []os.FileInfo
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		if !strings.EqualFold(filepath.Ext(fd.Name()), ".json") {
			continue
		}

		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if err = mergeJSONFile(srcfp, dstfp); err != nil {
			return err
		}
	}
	return nil
}

func mergeJSONFile(dstFile, srcFile string) error {
	dstRead, err := os.Open(dstFile)
	if err != nil {
		return err
	}
	defer dstRead.Close()

	dstMap := make(map[string]interface{})
	json.NewDecoder(dstRead).Decode(&dstMap)

	srcRead, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer srcRead.Close()

	srcMap := make(map[string]interface{})
	json.NewDecoder(srcRead).Decode(&srcMap)

	mergeJSON(dstMap, srcMap)

	write, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer write.Close()

	return json.NewEncoder(write).Encode(dstMap)
}

func mergeJSON(dstMap, srcMap map[string]interface{}) {
	for key := range dstMap {
		mergeJSONKey(key, dstMap[key], srcMap[key], dstMap)
	}
}

func mergeJSONKey(key string, dst interface{}, src interface{},
	result map[string]interface{}) {
	if !reflect.DeepEqual(dst, src) {
		switch dst.(type) {
		case map[string]interface{}:
			if _, ok := src.(map[string]interface{}); ok {
				dstMap := dst.(map[string]interface{})
				srcMap := src.(map[string]interface{})
				for k := range dstMap {
					mergeJSONKey(k, dstMap[k], srcMap[k], dstMap)
				}
			}
		default:
			if src != nil {
				result[key] = src
			}
		}
	}
}
