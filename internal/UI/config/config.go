package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	PostgreSQLURL string
	ClickHouseURL string
	Table         string
	Limit         int
	BatchSize     int
	Polling       PollingConfig
}

type PollingConfig struct {
	Enabled  bool
	Deltacol string
	Interval int
}

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = ".pgtoch.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("config file not found")
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, errors.New("cant parse config file")
	}

	return &config, nil
}
