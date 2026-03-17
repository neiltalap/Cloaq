package cli

import (
	network "cloaq/src"
	"cloaq/src/monitor"
	"cloaq/src/routing"
	"cloaq/src/tun"
	"log"
	"runtime"
	"sync/atomic"
	"time"
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
	log.Println("starting cloaq...")
	log.Println("goos:", runtime.GOOS, "goarch:", runtime.GOARCH)
	log.Println("starting tunnel on port", s.port, "with peers:", s.peers)

	tr, _ := network.NewTransport(":9000")

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

	// setting up readloop
	time.Sleep(100 * time.Millisecond)

	packetChan := make(chan network.Packet, 100)
	go network.SafeRuntime("ReadLoop", func() {
		if err := network.ReadLoop(dev, packetChan); err != nil {
			log.Println("readloop error:", err)
		}
	})

	// Initialize the router
	router := &routing.Router{}

	// Example static routes
	_ = router.AddRoute("2001:db8:1::/64", "eth0")
	_ = router.AddRoute("2001:db8:2::/64", "eth1")

	log.Println("ipv6 tun gateway created")

	//deleting the select{} to a better way of handling the readloop
	//select{}
	log.Println("ipv6 tun gateway created")

	for pkt := range packetChan {

		if len(s.peers) > 0 {
			target := s.peers[0]

			err := tr.SendTo(target, pkt.Data)
			if err == nil {
				atomic.AddUint64(&monitor.BytesSent, uint64(len(pkt.Data)))
				log.Printf("[sent] %d bytes to %s", len(pkt.Data), target)
			}
			log.Printf("[sent] %d bytes to %s", len(pkt.Data), target)
		}
	}
	m := &monitor.Monitor{}

	go network.SafeRuntime("Monitor", func() {
		if err := m.Execute(nil); err != nil {
			log.Printf("monitor error: %v", err)
		}
	})

	return nil
}
