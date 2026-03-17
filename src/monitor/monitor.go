package monitor

import (
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	BytesSent     uint64
	BytesReceived uint64
)

// Структура Monitor (с большой буквы!)
type Monitor struct{}

func (h *Monitor) Name() string {
	return "monitor"
}

func (h *Monitor) Description() string {
	return "display monitoring information"
}

func (h *Monitor) Execute(args []string) error {
	const (
		ByteToMiB = 1024 * 1024
		Interval  = 10 * time.Second
	)

	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			log.Printf("[monitor] RAM: %d MB | Sent: %d bytes | Goroutines: %d",
				m.Alloc/ByteToMiB,
				atomic.LoadUint64(&BytesSent),
				runtime.NumGoroutine(),
			)

			time.Sleep(Interval)
		}
	}()

	return nil
}
