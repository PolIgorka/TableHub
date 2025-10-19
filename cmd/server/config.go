package main

import (
	"os"

	"github.com/dbaas/internal/server"
	db "github.com/dbaas/pkg/database"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server server.ServerConfig `yaml:"http-server"`
	DB     db.DBConfig         `yaml:"database"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	config := Config{}
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if db_user := os.Getenv("DB_USER"); db_user != "" {
		config.DB.User = db_user
	}
	if db_password := os.Getenv("DB_PASSWORD"); db_password != "" {
		config.DB.Password = db_password
	}
	return &config, nil
}
