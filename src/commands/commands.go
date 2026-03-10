package commands

import (
	network "cloaq/src"
	"cloaq/src/routing"
	"cloaq/src/tun"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"
)

// imitation of config , in real project it will be yaml file
type Config struct {
	IdentityPath string
	Port         int
}

var currentConfig = Config{IdentityPath: "./id.key", Port: 8080}

func RunCommand(port int, peers string) {
	log.Println("starting cloaq...")
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

	// Prevent program from exiting
	select {}
}

func HelpCommand() {

	helpText := `
Cloaq — Universal Decentralized Anonymity Layer (UDAL)

Usage:
  cloaq [command] [arguments]

Available Commands:
  start          Initialize the TUN/TAP interface and join the network.
  stop           Gracefully shut down the node and close active tunnels.
  status         Display node health, connected peers, and network stats.
  settings       Manage configuration (ports, interface names, identity).
  gen-identity   Generate a new cryptographic identity (keypair).
  version        Print the current build version and architecture.

Flags:
  -h, --help     Show this help message.
  -v, --verbose  Enable debug logging for troubleshooting.

Examples:
  $ cloaq gen-identity --path ./id.key
  $ cloaq start --config ./config.yaml
  $ cloaq settings --port 9090

Use "cloaq [command] --help" for more information about a command.
`
	fmt.Print(helpText)

}

func SettingsCommand(args []string) {
	settingsFlags := flag.NewFlagSet("settings", flag.ExitOnError)

	// Опрелеляем флаги для подкоманды
	newPath := settingsFlags.String("path", currentConfig.IdentityPath, "Путь к файлу идентификации")
	newPort := settingsFlags.Int("port", currentConfig.Port, "Порт для gRPC/Transport")

	if len(args) == 0 {
		fmt.Printf("Текущие настройки:\n  Path: %s\n  Port: %d\n", currentConfig.IdentityPath, currentConfig.Port)
		return
	}

	err := settingsFlags.Parse(args)
	if err != nil {
		return
	}

	// Логика обновления
	currentConfig.IdentityPath = *newPath
	currentConfig.Port = *newPort

	fmt.Println("✅ Настройки успешно обновлены.")
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
