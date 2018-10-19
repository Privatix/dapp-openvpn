package openvpn

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

const ovpnPrefix = "dapp_ovpn"

type service struct {
	ID          string
	GUID        string
	Name        string
	Description string
	Command     string
	Args        []string
	AutoStart   bool
}

func diff(a, b []string) string {
	for i, v := range b {
		if len(a) <= i || a[i] != v {
			return v
		}
	}
	return ""
}

func ovpnName(path string) string {
	return fmt.Sprintf("%s_%s", ovpnPrefix, hash(path))
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func isServiceRun(service string) (bool, error) {
	output, err := exec.Command("sc", "queryex", service).CombinedOutput()

	if err != nil {
		return false, err
	}

	return strings.Contains(string(output), "RUNNING"), nil
}

func isServiceStop(service string) (bool, error) {
	output, err := exec.Command("sc", "queryex", service).CombinedOutput()

	if err != nil {
		return false, err
	}

	return strings.Contains(string(output), "STOPPED"), nil
}

func nextFreePort(h host) int {
	port := h.Port
	for i := port; i < 65535; i++ {
		ln, err := net.Listen(h.Protocol, h.IP+":"+strconv.Itoa(i))
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
