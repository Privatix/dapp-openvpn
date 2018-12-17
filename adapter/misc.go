package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	chanPerm = 0644
)

func encode(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

func commonNameOrEmpty() string {
	return os.Getenv("common_name")
}

func commonName() string {
	cn := commonNameOrEmpty()
	if len(cn) == 0 {
		logger.Fatal("empty common_name")
	}
	return cn
}

func storeChannel(cn, ch string) {
	name := filepath.Join(conf.ChannelDir, encode(cn))

	logger := logger.Add("method", "storeChannel", "file", name, "channel", ch)

	err := ioutil.WriteFile(name, []byte(ch), chanPerm)
	if err != nil {
		logger.Fatal("failed to store channel: " + err.Error())
	}
}

func loadChannel() string {
	name := filepath.Join(conf.ChannelDir, encode(commonName()))

	logger := logger.Add("method", "loadChannel", "file", name)

	data, err := ioutil.ReadFile(name)
	if err != nil {
		logger.Fatal("failed to load channel: " + err.Error())
	}

	return string(data)
}

func getCreds() (string, string) {
	logger := logger.Add("method", "getCreds")

	user := os.Getenv("username")
	pass := os.Getenv("password")

	if len(user) != 0 && len(pass) != 0 {
		return user, pass
	}

	if flag.NArg() < 1 {
		logger.Fatal("no filename passed to read credentials")
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		logger.Fatal(
			"failed to open file with credentials: " + err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	user = scanner.Text()
	scanner.Scan()
	pass = scanner.Text()

	if err := scanner.Err(); err != nil {
		logger.Fatal(
			"failed to read file with credentials: " + err.Error())
	}

	return user, pass
}

func storeActiveChannel(ch string) {
	name := filepath.Join(conf.ChannelDir, "active")

	logger := logger.Add("method", "storeActiveChannel",
		"file", name, "channel", ch)

	err := ioutil.WriteFile(name, []byte(ch), chanPerm)
	if err != nil {
		logger.Fatal("failed to store active channel: " + err.Error())
	}
}

func loadActiveChannel() string {
	name := filepath.Join(conf.ChannelDir, "active")

	logger := logger.Add("method", "loadActiveChannel", "file", name)

	data, err := ioutil.ReadFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		logger.Fatal("failed to load active channel: " + err.Error())
	}

	return string(data)
}

func removeActiveChannel() {
	name := filepath.Join(conf.ChannelDir, "active")

	logger := logger.Add("method", "removeActiveChannel", "file", name)

	if err := os.Remove(name); err != nil && !os.IsNotExist(err) {
		logger.Fatal("failed to remove active channel: " + err.Error())
	}
}
