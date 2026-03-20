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
	"fmt"
	"log"
	"runtime/debug"
)

type Version struct{}

// Ensure Version implements the Command interface
var _ Command = (*Version)(nil)

func (v *Version) Name() string {
	return "version"
}

func (v *Version) Description() string {
	return "print the current build version and system architecture"
}

func (v *Version) Execute(args []string) error {
	log.Println("----- [version] -----")

	version := "v0.1.0-alpha"
	revision := "unknown"
	goVersion := "unknown"

	// 1. Retrieve build information from the binary itself
	if info, ok := debug.ReadBuildInfo(); ok {
		goVersion = info.GoVersion
		for _, setting := range info.Settings {
			// Extract the VCS revision (Git commit hash) if available
			if setting.Key == "vcs.revision" {
				revision = setting.Value
				break
			}
		}
	}

	// 2. Format and display the output
	fmt.Printf("cloaq version: %s\n", version)
	fmt.Printf("go runtime:    %s\n", goVersion)
	fmt.Printf("vcs revision:  %s\n", revision)

	return nil
}
