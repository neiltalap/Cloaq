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

	if err := tun.Start(); err != nil {
		fmt.Println("Tunnel start error:", err)
		return
	}

	fmt.Println("Reading packets from tunnel...")
	if err := network.ReadLoop(tun); err != nil {
		fmt.Println("ReadLoop error:", err)
		return
	}
}

func helpCommand() {
	log.Println("help text")
}

func settingsCommand() {
	log.Println("settings text")
}
