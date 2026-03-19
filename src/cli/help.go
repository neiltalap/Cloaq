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
)

type Help struct{}

var _ Command = (*Help)(nil) // enforcement of an interface

func (h *Help) Name() string {
	return "help"
}

func (h *Help) Description() string {
	return "display help information"
}

func (h *Help) Execute(args []string) error {
	log.Println("----- [help] -----")

	fmt.Println("cloaq — Universal Decentralized Anonymity Layer (UDAL)")
	fmt.Println("\nUsage:")
	fmt.Println("  cloaq [command] [arguments]")

	fmt.Println("\nAvailable Commands:")

	// 1. Dynamically list all registered commands
	for _, cmd := range Commands {
		// We use padding to align descriptions for better readability
		fmt.Printf("  %-15s %s\n", cmd.Name(), cmd.Description())
	}

	fmt.Println("\nFlags:")
	fmt.Println("  -h, --help      Show this help message")
	fmt.Println("  -v, --verbose   Enable debug logging")

	fmt.Println("\nExamples:")
	fmt.Println("  $ cloaq settings --port 9090")
	fmt.Println("  $ cloaq start")

	fmt.Println("\nUse \"cloaq [command] --help\" for more information about a specific command.")

	return nil
}
