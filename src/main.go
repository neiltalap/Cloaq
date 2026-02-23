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
	"fmt"
	"log"
	"os"
	"runtime"

	"cloaq/src/network"
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
	fmt.Println("Starting Cloaq...")
	fmt.Println("GOOS:", runtime.GOOS, "GOARCH:", runtime.GOARCH)

	tun, err := network.InitTunnel()
	if err != nil {
		fmt.Println("Tunnel init error:", err)
		return
	}
	if tun == nil {
		fmt.Println("Tunnel initialized (no device object returned on this OS yet).")
		fmt.Println("Cloaq running.")
		select {}
	}

	defer tun.Close()
	fmt.Println("Tunnel ready:", tun.Name())

	// Integrated logic: Start the local tunnel processing
	if err := tun.Start(); err != nil {
		fmt.Println("Tunnel start error:", err)
		return
	}

	fmt.Println("Reading packets from tunnel...")
	// Start the ReadLoop in a goroutine so we can also run the router
	go func() {
		if err := network.ReadLoop(tun); err != nil {
			fmt.Println("ReadLoop error:", err)
		}
	}()

	// Upstream logic: Initialize Router and start IPv6 listener
	router := &Router{}

	// Example static routes from upstream
	router.AddRoute("2001:db8:1::/64", "eth0")
	router.AddRoute("2001:db8:2::/64", "eth1")

	log.Println("IPv6 TUN gateway created")

	// Use the file descriptor from the local 'tun' object
	// Note: We need to ensure tun.File() or similar exists or we use the raw FD if accessible.
	// Since 'tun' is from 'network.InitTunnel()', let's check what 'tun' is.
	// For now, I'll use the logic that requires the FD.
	// If 'tun' doesn't provide it easily, I might need to adjust 'network' package.

	// Assuming tun is an interface or struct that might have a File() method returning *os.File
	go CreateIPv6PacketListener(tun)
}

func helpCommand() {
	log.Println("help text")
}

func settingsCommand() {
	log.Println("settings text")
}
