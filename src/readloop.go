package network

import (
	"cloaq/src/utils"
	"fmt"
	"io"
)

func ReadLoop(device io.Reader, packetChan chan<- utils.Packet) error {
	if device == nil {
		return fmt.Errorf("tun device is nil")
	}

	defer close(packetChan)

	buf := make([]byte, 2048)
	for {
		n, err := device.Read(buf)
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		if n >= 20 {
			packetData := make([]byte, n)
			copy(packetData, buf[:n])

			packetChan <- utils.Packet{
				Data:    packetData,
				Version: buf[0] >> 4,
			}
		}
	}
}
