package openvpn

import (
	"crypto/sha1"
	"encoding/hex"
	"net"
	"strconv"
)

const ovpnPrefix = "dapp_ovpn"

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
