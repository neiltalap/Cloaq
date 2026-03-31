// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package config

import (
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	// path to the file because privateKey is a secret
	PrivateKeyPath string `yaml:"privateKey"`
	Port           int    `yaml:"port"`
}

func Load(filePath string) (*Config, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tempConfig *Config
	// attempting to unmarshal the config in the yaml
	if err := yaml.Unmarshal(fileData, &tempConfig); err != nil {
		return nil, err
	}

	return tempConfig, nil
}

func Init(filePath string) (*Config, error) {
	tempConfig, err := Load(filePath)
	if err != nil {
		return nil, err
	}

	return tempConfig, nil
}

func (config *Config) Save(filePath string) (*Config, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	// we allow for read and write for the user, except execution
	// the world and the group are granted read permissions only
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return nil, err
	}

	// the logic here is we return what we saved
	return config, nil
}
