package commands

import (
	"fmt"
	"runtime/debug"
)

func VersionCommand() {
	version := "v0.1.0-alpha"
	revision := "unknown"
	goVersion := "unknown"

	if info, ok := debug.ReadBuildInfo(); ok {
		goVersion = info.GoVersion
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				revision = setting.Value
				break
			}
		}
	}

	fmt.Printf("cloaq %s\n", version)
	fmt.Printf("go version: %s\n", goVersion)
	fmt.Printf("revision:   %s\n", revision)
}
