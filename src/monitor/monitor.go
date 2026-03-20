package cli

import (
	"log"
	"runtime"
	"time"
)

type Monitor struct{}

var _ Command = (*Monitor)(nil) // enforcement of an interface

func (h *Monitor) Name() string {
	return "monitor"
}

func (h *Monitor) Description() string {
	return "display monitoring information"
}

func (h *Monitor) Execute(args []string) error {

	const (
		ByteToMiB = 1024 * 1024
		Interval  = 30 * time.Second
	)
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			log.Printf("[monitor] alloc: %d MB, sys: %d MB, goroutines: %d",
				m.Alloc/ByteToMiB,
				m.Sys/ByteToMiB,
				runtime.NumGoroutine(),
			)

			time.Sleep(Interval)
		}
	}()

	return nil
}
