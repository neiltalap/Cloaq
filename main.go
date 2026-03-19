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
	cli "cloaq/src/cli"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("usage: cloaq <command>")
		return
	}

	commandName := os.Args[1]
	commandArguments := os.Args[2:]
	for _, command := range cli.Commands {
		if command.Name() != commandName {
			continue // skip commands that don't match user input
		}

		switch cmd := command.(type) {
		case *cli.Run:
			cmd.Execute(commandArguments)
		case *cli.Help:
			cmd.Execute(commandArguments)
		case *cli.Settings:
			cmd.Execute(commandArguments)
		default:
			log.Println("unknown command:", commandName)
		}
	}
}
