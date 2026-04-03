package cli

import (
	network "cloaq/src/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
)

type Run struct {
	port  int      `yaml:"port"`
	peers []string `yaml:"peers"`
}

var _ Command = (*Run)(nil) // enforcement of an interface

func (s *Run) Name() string {
	return "run"
}

func (s *Run) Description() string {
	return "display configuration run"
}
func (s *Run) Execute(args []string) error {
	fmt.Println(`_________ .__                       
\_   ___ \|  |   _________    ______
/    \  \/|  |  /  _ \__  \  / ____/
\     \___|  |_(  <_> ) __ \< <_|  |
 \______  /____/\____(____  /\__   |
        \/                \/    |__|`)

	if os.Geteuid() != 0 {
		log.Fatal("error: Run as root (sudo) to manage TUN device.")
	}

	node, err := NewCloaqNode(s.peers)
	if err != nil {
		log.Fatal("node init failed:", err)
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		log.Println("\n[!] signal received, cleaning up...")
		node.Shutdown()
		os.Exit(0)
	}()

	log.Printf("node %s started on fc00::1", node.ID[:12])

	packetChan := make(chan network.Packet, 100)
	node.Run(packetChan)

	for pkt := range packetChan {
		node.ProcessPacket(pkt)
	}

	return nil
}
