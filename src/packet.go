package network

import (
	"encoding/binary"
	"log"
	"runtime/debug"
)

type Packet struct {
	Data    []byte
	Version uint8
}

// Encapsulate adds a 4-byte header: [version][type][len_high][len_low]
func Encapsulate(raw []byte) []byte {
	size := len(raw)
	buf := make([]byte, 4+size)

	buf[0] = 0x01 // version
	buf[1] = 0x07 // type: data
	binary.BigEndian.PutUint16(buf[2:4], uint16(size))

	copy(buf[4:], raw)
	return buf
}

// SafeRuntime prevents goroutine panics from crashing the app
func SafeRuntime(name string, f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[panic] %s recovered: %v", name, r)
				log.Printf("[stack] %s trace: %s", name, string(debug.Stack()))
			}
		}()
		f()
	}()
}
