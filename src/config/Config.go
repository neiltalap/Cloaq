package config

import (
	"encoding/json"
	"os"
)

const configFileName = "config.yaml"

type Config struct {
	IdentityPath string `yaml:"identity_path"`
	Port         int    `yaml:"port"`
	Interface    string `yaml:"interface"`
}

var currentConfig = Config{IdentityPath: "./id.key", Port: 8080}

func LoadConfig() (*Config, error) {
	c := &Config{IdentityPath: "./id.key", Port: 8080}
	data, err := os.ReadFile("config.json")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(data, c)
	return c, err
}

func SaveConfig(c *Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}
