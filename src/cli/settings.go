// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package cli

import (
	"cloaq/src/config"
	"flag"
	"fmt"
	"log"
	"os"
)

type Settings struct{}

var _ Command = (*Settings)(nil) // enforcement of an interface

func (s *Settings) Name() string {
	return "settings"
}

func (s *Settings) Description() string {
	return "display configuration settings"
}

func (s *Settings) Execute(args []string) error {
	log.Println("----- [settings] -----")

	// 1. Load the existing configuration or initialize with defaults
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("notice: Configuration file not found, proceeding with default values.")
	}

	// 2. Define flags for the settings command
	fs := flag.NewFlagSet("settings", flag.ContinueOnError)
	newPath := fs.String("path", cfg.IdentityPath, "Path to the cryptographic identity file")
	newPort := fs.Int("port", cfg.Port, "Network port for the node to listen on")

	// 3. If no arguments are provided, display current settings and exit
	if len(args) == 0 {
		fmt.Printf("current Node Configuration:\n")
		fmt.Printf("  identity Path: %s\n", cfg.IdentityPath)
		fmt.Printf("  server Port:   %d\n", cfg.Port)
		return nil
	}

	// 4. Parse command-line arguments
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	// 5. Update the configuration structure
	cfg.IdentityPath = *newPath
	cfg.Port = *newPort

	// 6. Security Enforcement: Ensure the identity path has restrictive permissions (0600)
	// This satisfies the security requirements for handling sensitive keys.
	if _, err := os.Stat(cfg.IdentityPath); err == nil {
		if err := os.Chmod(cfg.IdentityPath, 0600); err != nil {
			log.Printf("warning: Failed to enforce 0600 permissions on %s: %v", cfg.IdentityPath, err)
		} else {
			log.Printf("security: File permissions for %s secured (0600)", cfg.IdentityPath)
		}
	}

	// 7. Persist updated configuration to disk
	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to persist configuration to config.json: %w", err)
	}

	fmt.Println("configuration updated and saved successfully.")
	return nil
}
