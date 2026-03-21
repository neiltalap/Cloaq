// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package cli

import (
	network "cloaq/src"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"log"
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
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
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
