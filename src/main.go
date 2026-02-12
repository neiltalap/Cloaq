package main

import (
	"log"
	"os"

	"cloaq/src/node"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: cloaq <command>")
		return
	}

	switch os.Args[1] {
	case "run":
		log.Println("Running Cloaq")

		if node.Bootstrap() {
			node.CreateListener() // listen to packets that will be passed via TUN
		}

	case "settings":
		settingsCommand()
	case "help":
		helpCommand()
	default:
		log.Println("Unknown command:", os.Args[1])
	}
}

func helpCommand() {
	log.Println("help text")
}

func settingsCommand() {
	log.Println("settings text")
}
