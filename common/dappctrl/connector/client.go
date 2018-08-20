package connector

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

// Response is a dappctrl server reply.
type Response struct {
	Error  *Error          `json:"error,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

// httpClient makes new http client.
func httpClient(config *Config) *http.Client {
	return &http.Client{
		Transport: transport(config),
		Timeout: time.Duration(
			config.RequestTimeout) * time.Millisecond,
	}
}

func transport(config *Config) *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(
				config.DialTimeout) * time.Millisecond,
			DualStack: true,
		}).DialContext,
		ResponseHeaderTimeout: time.Duration(
			config.ResponseHeaderTimeout) * time.Millisecond,
	}
}

func url(config *Config, path string) string {
	var proto = "http"
	if config.TLS != nil {
		proto += "s"
	}

	return proto + "://" + config.SessionServerAddr + path
}

func request(url, username, password string,
	args interface{}) (*http.Request, error) {
	reqData, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url,
		bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)

	return req, err
}

func send(httpClient *http.Client, req *http.Request,
	result interface{}) error {
	httpResp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	var ctrlResp Response
	if err = json.NewDecoder(httpResp.Body).Decode(&ctrlResp); err != nil {
		return err
	}

	if ctrlResp.Error != nil {
		return ctrlResp.Error
	}

	if ctrlResp.Result != nil && string(ctrlResp.Result) == "null" {
		return nil
	}

	return json.Unmarshal(ctrlResp.Result, result)
}

// post posts a request with given arguments and returns a response result.
func post(httpClient *http.Client, url, username, password string,
	args, result interface{}) error {
	req, err := request(url, username, password, args)
	if err != nil {
		return err
	}

	return send(httpClient, req, result)
}
