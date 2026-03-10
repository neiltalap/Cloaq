// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package commands

import (
	network "cloaq/src"
	"cloaq/src/routing"
	"cloaq/src/tun"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type Command interface {
	Name() string
	Description() string
	Execute(args []string) error
}

type CLI struct {
	Commands map[string]Command
}

func MakeCommandRegistry(cmds ...Command) *CLI {
	registry := make(map[string]Command)
	for _, cmd := range cmds {
		registry[cmd.Name()] = cmd
	}
	return &CLI{Commands: registry}
}

func (cli *CLI) Execute() {
	if len(os.Args) < 2 {
		log.Println("Usage: cloaq <command>")
		return
	}

	commandName := os.Args[1]
	command, exists := cli.Commands[commandName]
	if !exists {
		log.Println("Unknown command:", commandName)
	}

	if err := command.Execute(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
}

type RunCommand struct{}

func (r *RunCommand) Name() string {
	return "run"
}

func (r *RunCommand) Description() string {
	return "Start Cloaq gateway"
}

func (r *RunCommand) Execute(args []string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("run as root")
	}

	port := flag.Int("port", 8080, "Port to listen on")
	peers := flag.String("peers", "", "Comma-separated list of peer addresses")
	err := flag.CommandLine.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	fmt.Printf("Starting Cloaq on port %d...\n", *port)
	log.Println("goos:", runtime.GOOS, "goarch:", runtime.GOARCH)
	log.Println("starting tunnel on port", port, "with peers:", peers)

	// Initialize the identity for this node
	identity, err := network.GenerateIdentity()
	if err != nil {
		log.Fatal("identity creation failed: ", err)
	}
	// Logging the pubkey of the identity
	log.Println("current node's pubkey: ", string(identity.PublicKey.Bytes()))

	// Initialization of the VNIC on the node
	dev, err := tun.InitDevice()
	if err != nil {
		log.Fatal("tunnel init error:", err)
	}
	defer func() {
		err := dev.Close()
		if err != nil {
			log.Println("error closing device:", err)
		}
	}()

	log.Println("vnic has been initialized:", dev.Name())

	// Start VNIC processing
	if err := dev.Start(); err != nil {
		log.Fatal("vnic start error:", err)
	}

	log.Println("reading packets from the vnic...")

	startMonitor()
	log.Println("monitor started")

	// setting up readloop
	go func() {
		if err := network.ReadLoop(dev); err != nil {
			log.Println("readloop error:", err)
		}
	}()

	// Initialize the router
	router := &routing.Router{}

	// Example static routes
	_ = router.AddRoute("2001:db8:1::/64", "eth0")
	_ = router.AddRoute("2001:db8:2::/64", "eth1")

	log.Println("ipv6 tun gateway created")

	return nil
}

type HelpCommand struct {
	CLI *CLI
}

func (h *HelpCommand) Name() string {
	return "help"
}

func (h *HelpCommand) Description() string {
	return "Show help information"
}

func (h *HelpCommand) Execute(args []string) error {
	for _, cmd := range h.CLI.Commands {
		log.Printf("%s - %s\n", cmd.Name(), cmd.Description())
	}
	return nil
}

type SettingsCommand struct{}

func (s *SettingsCommand) Name() string {
	return "settings"
}

func (s *SettingsCommand) Description() string {
	return "Display configuration settings"
}

func (s *SettingsCommand) Execute(args []string) error {
	log.Println("----- [SETTINGS] -----")
	return nil
}

func startMonitor() {
	go func() {
		var m runtime.MemStats
		for {
			runtime.ReadMemStats(&m)

			log.Println("[monitor] alloc:", m.Alloc/1024/1024, "mb, sys:", m.Sys/1024/1024, "mb, goroutines:", runtime.NumGoroutine())

			time.Sleep(10 * time.Second)
		}
	}()
}
