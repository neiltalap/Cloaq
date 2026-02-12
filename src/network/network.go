package network

import (
	"log"

	// "github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

/*

We will be using TUN to build our network library because UDAL operates at **Layer 3 (IP level)** and needs a clean way to capture and inject IP packets in user space without dealing with Ethernet, MAC addresses, or ARP.

* Direct access to raw IP packets
* Works naturally with the OS routing stack
* Lets UDAL act like a virtual router
* Avoids Layer 2 complexity (which TAP would introduce)
------------------------------------------------------
TUN is the only option that gives structured IP control without unnecessary lower-layer overhead or routing conflicts:

* Raw sockets don’t integrate cleanly with system routing and require manual packet management, making full traffic interception unreliable and messy.
* AF_PACKET/libpcap operate at Layer 2 and are designed for sniffing/monitoring, not structured packet forwarding in an overlay network.
* TAP introduces Ethernet-level complexity (frames, broadcasts, MAC handling) that we don’t need for an IP-based privacy mesh.

We’re using water because it’s a lightweight Go library that gives direct access to TUN/TAP interfaces without forcing us into a heavy framework or OS-specific code.

*/

func SendPacket() {
	log.Println("Sending packet")
	iface, _ := water.New(water.Config{
		DeviceType: water.TUN,
	})

	buf := make([]byte, 1500)
	log.Println(iface.Read(buf))
}
