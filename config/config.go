package config

import (
	"io"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"strings"
)

type Config struct {
	baseURL string
	errors  []string
}

type configJSON struct {
	BaseURL string `json:"baseURL"`
}

func LoadConfig(configFilePath string) *Config {
	contents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return &Config{errors: []string{
			fmt.Sprintf("failed to open configuration file '%s'", configFilePath),
			err.Error(),
		}}
	}
	reader := bytes.NewBuffer(contents)
	return newConfig(reader)
}

func newConfig(reader io.Reader) *Config {
	config := &Config{}
	configJSON := &configJSON{}

	err := json.NewDecoder(reader).Decode(&configJSON)
	if err != nil {
		config.errors = append(config.errors, fmt.Sprintf("error parsing JSON configuration: %s", err))
		return config
	}

	config.baseURL = configJSON.BaseURL
	return config
}

func (config *Config) HasErrors() bool {
	return len(config.errors) > 0
}

func (config *Config) Errors() string {
	if config.HasErrors() {
		return fmt.Sprintf("invalid configuration: %s", strings.Join(config.errors, ": "))
	}
	return ""
}

func (config *Config) BaseURL() string{
	return config.baseURL
}
