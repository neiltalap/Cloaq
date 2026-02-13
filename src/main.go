package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cloaq <command>")
		return
	}

	switch os.Args[1] {
	case "circuits":
		UdatlCircuits()
	case "status":
		UdalStatus()
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
