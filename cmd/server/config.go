package main

import (
	"os"

	"github.com/tablehub/internal/server"
	db "github.com/tablehub/pkg/database"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server server.ServerConfig `yaml:"http-server"`
	DB     DataBase            `yaml:"database"`
}

type DataBase struct {
	RightsConfig db.DBConfig `yaml:"rights"`
	TablesConfig db.DBConfig `yaml:"tables"`
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

	if db_user := os.Getenv("RIGHTS_DB_USER"); db_user != "" {
		config.DB.RightsConfig.User = db_user
	}
	if db_password := os.Getenv("RIGHTS_DB_PASSWORD"); db_password != "" {
		config.DB.RightsConfig.Password = db_password
	}
	if db_user := os.Getenv("TABLES_DB_USER"); db_user != "" {
		config.DB.TablesConfig.User = db_user
	}
	if db_password := os.Getenv("TABLES_DB_PASSWORD"); db_password != "" {
		config.DB.TablesConfig.Password = db_password
	}
	return &config, nil
}
