// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package main

import (
	"cloaq/src/tun"
	"log"
	"net"
)

// Reads packets from Tunnel, FOR NOW it's just logs basic info (IPv4/IPv6)
func ReadLoop(dev tun.Device, tr *Transport, targetPeers []net.UDPAddr) error {
	buf := make([]byte, 65535)

	for {
		n, err := dev.Read(buf)
		if err != nil {
			return err
		}
		if n < 1 {
			continue
		}

		pkt := make([]byte, n)
		copy(pkt, buf[:n])

		finalPkt := make([]byte, len(pkt)+1)
		finalPkt[0] = 0x01
		copy(finalPkt[1:], pkt)

		for _, peerAddr := range targetPeers {
			addr := peerAddr
			go func(a net.UDPAddr, data []byte) {
				err := tr.SendTo(&a, data)
				if err != nil {
					log.Printf("Ошибка отправки на %s: %v", a.String(), err)
				}
			}(addr, finalPkt)
		}

		ver := pkt[0] >> 4
		log.Printf("[TUN -> NET] Отправлен IPv%d пакет (%d байт) на %d узлов", ver, n, len(targetPeers))
	}
}
