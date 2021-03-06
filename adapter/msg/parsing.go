package msg

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/privatix/dappctrl/util"
	"github.com/privatix/dappctrl/util/log"
)

func vpnParams(logger log.Logger, file string,
	keys []string) (params map[string]string, err error) {
	logger = logger.Add("method", "vpnParams", "file", file)

	params = make(map[string]string)

	f, err := os.Open(file)
	if err != nil {
		logger.Error(err.Error())
		return nil, ErrReadConfig
	}
	defer f.Close()

	keyMap := make(map[string]bool)

	for _, key := range keys {
		keyMap[strings.TrimSpace(key)] = true
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if key, value, add :=
			parseLine(keyMap, scanner.Text()); add {
			if key != "" {
				params[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error(err.Error())
		return nil, ErrReadConfig
	}

	return params, err
}

func certificateAuthority(logger log.Logger,
	file string) (ca []byte, err error) {
	logger = logger.Add("method", "certificateAuthority", "file", file)

	mainCertPEMBlock, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Error(err.Error())
		return nil, ErrReadCert
	}

	if !util.IsTLSCert(string(mainCertPEMBlock)) {
		logger.Error(err.Error())
		return nil, ErrFindCert
	}

	return mainCertPEMBlock, nil
}

func parseLine(keys map[string]bool, line string) (string, string, bool) {
	str := strings.TrimSpace(line)

	for key := range keys {
		sStr := strings.Split(str, " ")
		if sStr[0] != key {
			continue
		}

		index := strings.Index(str, "#")

		if index == -1 {
			words := strings.Split(str, " ")
			if len(words) == 1 {
				return key, "", true
			}
			value := strings.Join(words[1:], " ")
			return key, value, true
		}

		subStr := strings.TrimSpace(str[:index])

		words := strings.Split(subStr, " ")

		if len(words) == 1 {
			return key, "", true
		}

		value := strings.Join(words[1:], " ")

		return key, value, true
	}
	return "", "", false
}
