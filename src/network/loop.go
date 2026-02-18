package network

import "fmt"

func ReadLoop(t Tunnel) error {
	buf := make([]byte, 65535)

	for {
		n, err := t.Read(buf)
		if err != nil {
			return fmt.Errorf("tun read: %w", err)
		}
		pkt := buf[:n]
		fmt.Println("[TUN]", DescribeIPPacket(pkt))
	}
}
