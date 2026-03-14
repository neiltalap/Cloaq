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
	"log"
	"os"

	"golang.org/x/sys/unix"
)

// Reads packets from Tunnel, FOR NOW it's just logs basic info (IPv4/IPv6)
func ReadLoop(f *os.File) error {
	if f == nil {
		return fmt.Errorf("file handle is nil")
	}

	fd := int(f.Fd())

	// Set file descriptor to blocking mode to avoid "not pollable" error in Go runtime
	err := unix.SetNonblock(fd, false)
	if err != nil {
		return fmt.Errorf("failed to set blocking mode: %w", err)
	}

	log.Println("ReadLoop: initialized successfully, waiting for traffic...")

	// Buffer size for standard MTU (1500) + some overhead
	buf := make([]byte, 2048)

	for {
		// Use direct syscall.Read to bypass Go's network poller
		n, err := unix.Read(fd, buf)

		if err != nil {
			// If the syscall was interrupted by a signal, just retry
			if err == unix.EINTR {
				continue
			}
			return fmt.Errorf("read error: %w", err)
		}

		if n > 0 {
			// Basic packet metadata logging
			// buf[0] >> 4 extracts the IP version (4 or 6)
			ipVersion := buf[0] >> 4
			fmt.Printf(">>> Received packet: %d bytes | IP Version: v%d\n", n, ipVersion)

			// TODO: Pass the packet to your processing logic
			// ProcessPacket(buf[:n])
		}
	}
}

func safeRuntime(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[RECOVERY] %s panicked: %v", name, r)
		}
	}()
	fn()
}
