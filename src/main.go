package main

import (
	"cloaq/src/network"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cloaq <command>")
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
		fmt.Println("Unknown command:", os.Args[1])
	}
}

func runCommand() {
	fmt.Println("Starting Cloaq Node...")

	id, err := GenerateIdentity()
	if err != nil {
		log.Fatal("Identity error:", err)
	}
	fmt.Printf("Node ID (Public Key): %s\n", id.String())

	tun, err := network.InitTunnel()
	if err != nil {
		log.Fatal("TUN error:", err)
	}
	defer func(tun *network.LinuxTunnel) {
		err := tun.Close()
		if err != nil {

		}
	}(tun)
	fmt.Printf("Interface %s is UP\n", tun.Name())

	go CreateIPv6PacketListener(tun)

	fmt.Println("Network is active. Press Ctrl+C to exit.")

	select {}
}

func helpCommand()     { fmt.Println("Commands: run, settings, help") }
func settingsCommand() { fmt.Println("Settings not available yet") }
