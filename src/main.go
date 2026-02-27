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

	"cloaq/src/routing"
	"cloaq/src/tun"
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

	dev, err := tun.InitDevice()
	if err != nil {
		fmt.Println("Tunnel init error:", err)
		return
	}
	if dev == nil {
		fmt.Println("Tunnel initialized (no device object returned on this OS yet).")
		fmt.Println("Cloaq running.")
		select {}
	}

	defer dev.Close()
	fmt.Println("Tunnel ready:", dev.Name())

	// Integrated logic: Start the local tunnel processing
	if err := dev.Start(); err != nil {
		fmt.Println("Tunnel start error:", err)
		return
	}

	fmt.Println("Reading packets from tunnel...")
	// Start the ReadLoop in a goroutine so we can also run the router
	go func() {
		if err := ReadLoop(dev); err != nil {
			fmt.Println("ReadLoop error:", err)
		}
	}()

	// Initialize UDP transport
	tr, err := NewTransport(":9000")
	if err != nil {
		log.Fatal(err)
	}
	// Incoming packets
	incoming := make(chan []byte, 1024)

	// Start UDP Listener
	go tr.Listen(incoming)

	// Write packets to TUN
	go func() {
		for pkt := range incoming {
			if err := tun.WritePacket(dev, pkt); err != nil {
				log.Println("write to tun failed:", err)
			}
		}
	}()

	// Upstream logic: Initialize Router and start IPv6 listener
	router := routing.NewRouter()

	// Example static routes from upstream
	router.AddRoute("2001:db8:1::/64", "eth0")
	router.AddRoute("2001:db8:2::/64", "eth1")

	log.Println("IPv6 TUN gateway created")

	go routing.CreateIPv6PacketListener(dev)
}

func helpCommand() {
	log.Println("help text")
}

func settingsCommand() {
	log.Println("settings text")
}
