package env

import (
	"encoding/json"
	"os"

	"github.com/privatix/dappctrl/util"
)

// Config has a store environment variables.
type Config struct {
	Schema         string
	Workdir        string
	Device         string
	Interface      string
	Service        string
	Role           string
	DappVPN        string
	ProductImport  bool
	ProductInstall bool
}

// NewConfig creates a default Configs configuration.
func NewConfig() *Config {
	return &Config{Schema: "1.0"}
}

// Write saves the configs to json file.
func (c *Config) Write(path string) error {
	write, err := os.Create(path)
	if err != nil {
		return err
	}
	defer write.Close()
	return json.NewEncoder(write).Encode(c)
}

// Read reads the configs from json file.
func (c *Config) Read(path string) error {
	return util.ReadJSONFile(path, &c)
}
