// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package monitor

import (
	"log"
	"runtime"
	"time"
)

// global counters for traffic tracking
var (
	BytesSent     uint64
	BytesReceived uint64
)

type Monitor struct {
	// add specific monitor config here if needed
}

func (m *Monitor) Name() string {
	return "monitor"
}

func (m *Monitor) Description() string {
	return "display real-time system and network metrics"
}

func (m *Monitor) Execute(args []string) error {
	log.Println("monitor subsystem started")

	// ticker for periodic stats reporting
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.logStats()
		}
	}
}

func (m *Monitor) logStats() {
	var mStats runtime.MemStats
	runtime.ReadMemStats(&mStats)

	// logging system metrics in lowercase
	log.Printf("[monitor] stats report:")
	log.Printf("  - traffic: sent: %d bytes | received: %d bytes", BytesSent, BytesReceived)
	log.Printf("  - memory: alloc: %v MiB | sys: %v MiB | goroutines: %d",
		bToMb(mStats.Alloc), bToMb(mStats.Sys), runtime.NumGoroutine())
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
