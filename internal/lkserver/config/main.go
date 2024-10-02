package config

import (
	"time"

	"github.com/BurntSushi/toml"
)

var ServerTimeZone *time.Location

type Config struct {
	BindAddr        string `toml:"bind_addr"`
	SessionsKey     string `toml:"session_key"`
	StaticFilesPath string `toml:"files"`
	SessionMaxAge   int    `toml:"session_max_age"`
	ServerTimeZone  string `toml:"timezone"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{
		BindAddr:        ":8833",
		SessionsKey:     "7e71a7e7-b86c-455c-9be6-dfd917144696",
		StaticFilesPath: "data",
		SessionMaxAge:   60 * 30,
		ServerTimeZone:  "Local",
	}
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return nil, err
	}

	ServerTimeZone, err = time.LoadLocation(config.ServerTimeZone)
	if err != nil {
		return nil, err
	}
	return config, nil
}
