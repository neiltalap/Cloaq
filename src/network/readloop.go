package network

import "log"

// Reads packets from Tunnel, FOR NOW it's just logs basic info (IPv4/IPv6)
func ReadLoop(tun Tunnel) error {
	buf := make([]byte, 65535)

	for {
		n, err := tun.Read(buf)
		if err != nil {
			return err
		}
		pkt := buf[:n]
		if len(pkt) < 1 {
			continue
		}

		ver := pkt[0] >> 4
		switch ver {
		case 6:
			log.Printf("IPv6 packet: %d bytes\n", len(pkt))
		case 4:
			log.Printf("IPv4 packet: %d bytes\n", len(pkt))
		default:
			log.Printf("Unknown packet version %d: %d bytes\n", ver, len(pkt))
		}
	}
}
