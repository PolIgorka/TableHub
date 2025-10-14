package main

import (
	"os"

	"gopkg.in/yaml.v3"
	"github.com/dbaas/internal/server"
)

type Config struct {
	Server server.ServerConfig `yaml:"http-server"`
}

func LoadConfig(path string) (*Config, error) {
	config := Config{}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
