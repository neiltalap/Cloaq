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
	"log"
	"net"
)

type Transport struct {
	conn *net.UDPConn
}

func NewTransport(listenAddr string) (*Transport, error) {
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Transport{
		conn: conn,
	}, nil
}

func (t *Transport) Close() error {
	if t.conn != nil {
		log.Println("[-] Closing UDP transport socket...")
		return t.conn.Close()
	}
	return nil
}

func (t *Transport) SendTo(addr string, data []byte) error {
	// Convert string "ip:port" to the correct pointer type
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	// FIX: Use udpAddr (*net.UDPAddr), NOT a uint8
	_, err = t.conn.WriteToUDP(data, udpAddr)
	return err
}

func (t *Transport) Listen(incoming chan<- []byte) {
	buf := make([]byte, 65535)

	for {
		n, _, err := t.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		packet := make([]byte, n)
		copy(packet, buf[:n])

		incoming <- packet
	}
}
