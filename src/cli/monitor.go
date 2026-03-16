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

	go func() {
		var m runtime.MemStats
		for {
			runtime.ReadMemStats(&m)

			log.Println("[monitor] alloc:", m.Alloc/1024/1024, "mb, sys:", m.Sys/1024/1024, "mb, goroutines:", runtime.NumGoroutine())

			time.Sleep(10 * time.Second)
		}
	}()
	return nil
}
