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
	"cloaq/src/cli"
	"cloaq/src/config"
	"flag"
	"fmt"

	"log"
	"os"
)

func main() {
	config.Init()

	if len(os.Args) < 2 {
		log.Println("Usage: cloaq <command>")
		return
	}

	port := flag.Int("port", 8080, "port to listen on")
	flag.Int("peers", 1, "number of peers to connect to")

	err := flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		return
	}

	fmt.Printf("starting Cloaq on port %d...\n", *port)

	commandName := os.Args[1]
	args := os.Args[2:]

	switch commandName {
	case "run":

		cmd := &cli.Run{}
		if err := cmd.Execute(args); err != nil {
			log.Fatalf("error while executing run command %v", err)
		}

	case "settings":
		cmd := &cli.Settings{}
		if err := cmd.Execute(args); err != nil {
			log.Fatalf("error: %v", err)
		}

	case "help":
		cmd := &cli.Help{}
		err := cmd.Execute(args)
		if err != nil {
			return
		}

	case "monitor":
		cmd := &cli.Monitor{}
		err := cmd.Execute(args)
		if err != nil {
			return
		}
	case "version":
		cmd := &cli.Version{}
		err := cmd.Execute(args)
		if err != nil {
			return
		}

	default:
		log.Printf("unknown command: %s", commandName)
		log.Fatal("use help to see available commands")
	}

}
