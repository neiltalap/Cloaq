package main

import (
	"fmt"
	"os"

	"cloaq/network"
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
	fmt.Println("Starting Cloaq...")

	err := network.InitWinTun()
	if err != nil {
		fmt.Println("WinTun error:", err)
		return
	}

	fmt.Println("Cloaq running (WinTun initialized)")
	select {}
}

func helpCommand() {
	fmt.Println("help text")
}

func settingsCommand() {
	fmt.Println("settings text")
}
