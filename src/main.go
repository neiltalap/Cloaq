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
	"cloaq/src/tun"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("use sudo ./cloaq run <ipGen_port>")
		return
	}

	switch os.Args[1] {
	case "run":
		runCommand()
	case "help":
		helpCommand()
	case "settings":
		settingsCommand()
	default:
		log.Println("unknown command", os.Args[1])
	}
}

func runCommand() {
	// start cloaq process
	log.Println("starting cloaq...")

	// Настройка перехвата сигналов для корректного выхода
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// initialize node identity
	id, err := GenerateIdentity()
	if err != nil {
		log.Fatalf("failed to generate identity: %v", err)
	}

	log.Printf("node public key: %s\n", id.String())

	// initialize tun device
	dev, err := tun.InitDevice()
	if err != nil {
		log.Fatal("failed to initialize tun:", err)
	}
	log.Println("interface ready:", dev.Name())

	// start the tun device
	if err := dev.Start(); err != nil {
		log.Fatal("failed to start tun:", err)
	}

	// setup udp transport
	tr, err := NewTransport(":9000")
	if err != nil {
		log.Fatal("transport error:", err)
	}

	// handle target peers from arguments
	var targetPeers []net.UDPAddr
	if len(os.Args) > 2 {
		addr, err := net.ResolveUDPAddr("udp", os.Args[2])
		if err == nil {
			targetPeers = append(targetPeers, *addr)
			log.Printf("added peer for broadcast: %s\n", addr.String())
		}
	}

	// listen for incoming network traffic
	incoming := make(chan []byte, 1024)
	go tr.Listen(incoming)

	// process incoming data and write to tun
	go func() {
		for data := range incoming {
			if len(data) < 1 {
				continue
			}

			// check packet type (0x01 = data)
			if data[0] == 0x01 {
				ipPacket := data[1:]

				err := tun.WritePacket(dev, ipPacket)
				if err != nil {
					log.Println("tun write error:", err)
				} else {
					log.Printf("[net -> tun] received and processed packet (%d bytes)", len(ipPacket))
				}
			}
		}
	}()

	// start main read loop for outgoing traffic
	go func() {
		log.Println("cloaq is running. waiting for traffic...")
		if err := ReadLoop(dev, tr, targetPeers); err != nil {
			log.Fatal("readloop failed:", err)
		}
	}()

	// waiting for shutting down

	select{}
	log.Println("shutting down...")
}
func helpCommand() {
	log.Println("help text")
}

func settingsCommand() {
	log.Println("settings text")
}
