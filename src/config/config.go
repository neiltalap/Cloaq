package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

// SaveConfig writes the provided configuration to the config.yaml file.
func SaveConfig(cfg *Config) error {
	// 1. Serialize the Config struct into YAML format
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config to yaml: %w", err)
	}

	// 2. Persist the data to the file system.
	// Using 0644 (Read/Write for owner, Read for others).
	err = os.WriteFile(configFileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

type IdentityStore struct {
	Keys [][]byte `yaml:"keys"`
}

func getStorePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(exe)
	return filepath.Join(dir, "store.yaml"), nil
}
func LoadStore() (*IdentityStore, error) {
	path, err := getStorePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &IdentityStore{}, nil
		}
		return nil, err
	}

	var store IdentityStore
	if err := yaml.Unmarshal(data, &store); err != nil {
		return nil, err
	}

	return &store, nil
}

func SaveStore(store *IdentityStore) error {
	path, err := getStorePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(store)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
