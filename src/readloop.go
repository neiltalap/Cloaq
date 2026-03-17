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

type Packet struct {
	Data    []byte
	Version uint8 // 4 for IPv4, 6 for IPv6
}

// Reads packets from Tunnel, FOR NOW it's just logs basic info (IPv4/IPv6)
func ReadLoop(device io.ReadCloser, packetChan chan<- Packet) error {
	if device == nil {
		return fmt.Errorf("TUN device is not initialized (nil)")
	}

	defer close(packetChan) // Close the channel if the loop exits
	log.Println("ReadLoop: Started. Ready to capture outbound traffic.")

	// Standard MTU is 1500. 2048 provides enough headroom for headers.
	buf := make([]byte, 2048)

	for {
		n, err := device.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("readLoop: Device closed, terminating loop.")
				return nil
			}
			return fmt.Errorf("read error on TUN device: %w", err)
		}

		if n > 0 {
			// 1. Basic Validation: IP header is at least 20 bytes for IPv4
			if n < 20 {
				continue
			}

			// 2. Extract Metadata
			// The first 4 bits of an IP packet indicate the version.
			version := buf[0] >> 4

			// 3. Create a copy of the packet data
			// This is crucial because 'buf' will be overwritten in the next iteration.
			packetData := make([]byte, n)
			copy(packetData, buf[:n])

			// 4. Send to Routing logic
			// We send it to a channel so the ReadLoop can immediately go back to reading.
			packetChan <- Packet{
				Data:    packetData,
				Version: version,
			}
		}
	}
}

func SafeRuntime(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[RECOVERY] %s panicked: %v", name, r)
		}
	}()
	fn()
}
