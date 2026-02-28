package main

import (
	"cloaq/src/network"
	"log"
)

func CreateIPv6PacketListener(tun *network.LinuxTunnel) {
	buf := make([]byte, 65535)
	for {
		n, err := tun.Read(buf)
		if err != nil {
			log.Println("TUN Read error:", err)
			continue
		}

		pkt := buf[:n]

		if len(pkt) < 40 {
			continue
		}

		if (pkt[0] >> 4) != 6 {
			continue
		}

		log.Printf("Captured IPv6 packet: %d bytes\n", n)
	}
}
