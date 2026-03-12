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
	"cloaq/src/routing"
	"cloaq/src/tun"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type Run struct{}

var _ Command = (*Run)(nil) // enforcement of an interface

func (r *Run) Name() string {
	return "run"
}

func (r *Run) Description() string {
	return "start cloaq gateway"
}

func (r *Run) Execute(args []string) error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("run as root")
	}

	fs := flag.NewFlagSet("run", flag.ExitOnError) // do not use global flags because flags from different commands can collide

	port := fs.Int("port", 8080, "port to listen on")
	peers := fs.String("peers", "", "comma-separated peers")

	err := fs.Parse(args)
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
