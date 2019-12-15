package config

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type Config struct {
	Mode   string
	Server *ServerConfig `ini:"server"`
}

type ServerConfig struct {
	Address string `ini:"address"`
}

func LoadConfig() (*Config, error) {

	var config Config
	if err := ini.MapTo(&config, "config.ini"); err != nil {
		return nil, fmt.Errorf("lad config error: %v", err)
	}

	return &config, nil
}
