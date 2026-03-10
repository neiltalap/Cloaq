// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package main

import (
	"cloaq/src/commands"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: cloaq <command>")
		return
	}

	commandRegistry := commands.MakeCommandRegistry(&commands.RunCommand{}, &commands.SettingsCommand{})
	help := &commands.HelpCommand{CLI: commandRegistry}
	commandRegistry.Commands[help.Name()] = help // hate that we have Commands and commands. needs to be changed.

	commandRegistry.Execute()
}
