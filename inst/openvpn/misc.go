package openvpn

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
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
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func nextFreePort(h host, proto string) int {
	port := h.Port
	for i := port; i < 65535; i++ {
		ln, err := net.Listen(proto, h.IP+":"+strconv.Itoa(i))
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

func connectorAddr(config string) (string, error) {
	read, err := os.Open(config)
	if err != nil {
		return "", err
	}
	defer read.Close()

	jsonMap := make(map[string]interface{})

	json.NewDecoder(read).Decode(&jsonMap)

	srv, ok := jsonMap["SessionServer"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("SessionServer params not found")
	}
	addr, ok := srv["Addr"]
	if !ok {
		return "", fmt.Errorf("Addr params not found")
	}
	return addr.(string), nil
}
