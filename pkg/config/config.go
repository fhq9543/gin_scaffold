package config

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type Config struct {
	*ini.File
}

var defaultIniPath = "./configs/config.ini"

var cfg *Config

func Init() (*Config, error) {
	prefix := "./"

	config, err := ini.Load(prefix + defaultIniPath)
	if err != nil {
		fmt.Errorf(err.Error())
		return nil, err
	}
	cfg = &Config{config}
	return &Config{config}, nil
}

func (c *Config) GetConfig(section string, key string) string {
	return c.Section(section).Key(key).String()
}

func GetConfig(section string, key string) string {
	return cfg.GetConfig(section, key)
}
