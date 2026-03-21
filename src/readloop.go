// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package network

import (
	"fmt"
	"io"
	"log"
)

// Reads packets from Tunnel, FOR NOW it's just logs basic info (IPv4/IPv6)
func ReadLoop(device io.Reader, packetChan chan<- Packet) error {
	if device == nil {
		return fmt.Errorf("TUN device is nil")
	}

	defer close(packetChan)
	log.Println("ReadLoop: Started.")

	buf := make([]byte, 2048)

	for {
		n, err := device.Read(buf)
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		if n >= 20 {
			packetData := make([]byte, n)
			copy(packetData, buf[:n])

			packetChan <- Packet{
				Data:    packetData,
				Version: buf[0] >> 4,
			}
		}
	}
}
