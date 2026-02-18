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
)

type Command interface {
	Name() string
	Description() string
	Execute(args []string) error
}

type CLI struct {
	commands map[string]Command
}

func MakeCommandRegistry(cmds ...Command) *CLI {
	registry := make(map[string]Command)
	for _, cmd := range cmds {
		registry[cmd.Name()] = cmd
	}
	return &CLI{commands: registry}
}

func (cli *CLI) Execute() {
	if len(os.Args) < 2 {
		log.Println("Usage: cloaq <command>")
		return
	}

	commandName := os.Args[1]
	command, exists := cli.commands[commandName]
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

	log.Println("Running Cloaq")

	tunFD, err := NewTUN("tun0")
	if err != nil {
		log.Println()
		return err
	}

	router := &Router{}

	router.AddRoute("2001:db8:1::/64", "eth0")
	router.AddRoute("2001:db8:2::/64", "eth1")

	log.Println("IPv6 TUN gateway created")

	router.CreateIPv6PacketListener(tunFD)

	return nil
}

type HelpCommand struct {
	cli *CLI
}

func (h *HelpCommand) Name() string {
	return "help"
}

func (h *HelpCommand) Description() string {
	return "Show help information"
}

func (h *HelpCommand) Execute(args []string) error {
	for _, cmd := range h.cli.commands {
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
