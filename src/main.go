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
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: cloaq <command>")
		return
	}

	switch os.Args[1] {
	case "run":
		runCommand()
	case "settings":
		settingsCommand()
	case "help":
		helpCommand()
	default:
		log.Println("Unknown command:", os.Args[1])
	}
}

func runCommand() {
	if os.Geteuid() != 0 {
		log.Fatal("Run as root") // privileged kernel networking operations
	}
	log.Println("Running Cloaq")

	// Create TUN interface
	tunFD := NewTUN("tun0")

	router := &Router{}

	// Example static routes
	router.AddRoute("2001:db8:1::/64", "eth0")
	router.AddRoute("2001:db8:2::/64", "eth1")

	log.Println("IPv6 TUN gateway created")

	// CreateRouter(tunFD) to listen and forward packets to another nodes
	router.CreateIPv6PacketListener(tunFD)
}

func helpCommand() {
	log.Println("help text")
}

func settingsCommand() {
	log.Println("settings text")
}
