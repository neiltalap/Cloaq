package main

import (
	"log"

	"cloaq/src/network"
)

func CreateIPv6PacketListener(tun network.Tunnel) {
	buf := make([]byte, 65535)
	for {
		n, err := tun.Read(buf)
		if err != nil {
			log.Println("tun.Read error:", err)
			continue
		}

		pkt := buf[:n]
		if len(pkt) < 40 {
			continue
		}
		if (pkt[0] >> 4) != 6 {
			continue
		}

		payload := pkt[40:]
		log.Printf("IPv6 packet: %d bytes, payload %d bytes\n", len(pkt), len(payload))
	}
}
