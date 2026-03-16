package config

import (
	"log"
	"os"

	"go.yaml.in/yaml/v4"
)

const configFileName = "config.yaml"

type Config struct {
	IdentityPath string `yaml:"identity_path"`
	Port         int    `yaml:"port"`
	Interface    string `yaml:"interface"`
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Init() {
	var err error
	AppConfig, err = LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
}
