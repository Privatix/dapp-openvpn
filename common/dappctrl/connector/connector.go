// Package connector implements standard methods for communicating with dappctrl.
package connector

import (
	"net/http"
)

// Session server API paths.
const (
	PathAuth          = "/session/auth"
	PathStart         = "/session/start"
	PathStop          = "/session/stop"
	PathUpdate        = "/session/update"
	PathProductConfig = "/product/config"
)

// Config is a connector configuration.
type Config struct {
	SessionServerAddr     string
	TLS                   *TLSConfig
	Username              string
	Password              string
	DialTimeout           int64
	ResponseHeaderTimeout int64
	RequestTimeout        int64
}

// TLSConfig is a TLS configuration.
type TLSConfig struct {
	CertFile string
	KeyFile  string
}

// DefaultConfig is a default connector config.
func DefaultConfig() *Config {
	return &Config{
		SessionServerAddr:     "localhost:80",
		TLS:                   nil,
		DialTimeout:           5,
		ResponseHeaderTimeout: 30,
		RequestTimeout:        40,
	}
}

// Connector defines the methods for interacting with a dappctrl.
type Connector interface {
	AuthSession(args interface{}) error
	StartSession(args interface{}) error
	StopSession(args interface{}) error
	UpdateSessionUsage(args interface{}) error
	SetupProductConfiguration(args interface{}) error
}

type cntr struct {
	config *Config
	client *http.Client
}

// NewConnector implements standard connector for communicating with dappctrl.
func NewConnector(config *Config) Connector {
	return &cntr{
		config: config,
		client: httpClient(config),
	}
}

// AuthSession sends a request for session authentication.
func (c *cntr) AuthSession(args interface{}) error {
	return post(c.client, url(c.config, PathAuth),
		c.config.Username, c.config.Password, args, nil)
}

// StartSession sends a request for session start.
func (c *cntr) StartSession(args interface{}) error {
	return post(c.client, url(c.config, PathStart),
		c.config.Username, c.config.Password, args, nil)
}

// StopSession sends a request for session stop.
func (c *cntr) StopSession(args interface{}) error {
	return post(c.client, url(c.config, PathStop),
		c.config.Username, c.config.Password, args, nil)
}

// UpdateSessionUsage sends a request to update
// a information on the use of session.
func (c *cntr) UpdateSessionUsage(args interface{}) error {
	return post(c.client, url(c.config, PathUpdate),
		c.config.Username, c.config.Password, args, nil)
}

// SetupProductConfiguration  sends a request to update product configuration.
func (c *cntr) SetupProductConfiguration(args interface{}) error {
	return post(c.client, url(c.config, PathProductConfig),
		c.config.Username, c.config.Password, args, nil)
}
