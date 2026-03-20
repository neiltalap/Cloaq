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
	"net"
	"strconv"
)

type Transport struct {
	conn *net.UDPConn
}

const HeaderSize = 16

type CloaqHeader struct {
	Version    uint8
	PacketType uint8
	SessionID  uint32
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

func (t *Transport) SendTo(targetAddr string, data []byte) error {
	addr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		return err
	}
	finalFrame := t.Encapsulate(data)

	_, err = t.conn.WriteToUDP(finalFrame, addr)
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

func SendUDP(addr uint8, data []byte) error {
	udpAddr, err := net.ResolveUDPAddr("udp", strconv.Itoa(int(addr)))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	_, err = conn.Write(data)
	return err
}

func (t *Transport) Encapsulate(data []byte) []byte {

	frame := make([]byte, len(data)+HeaderSize)

	frame[0] = 1
	frame[1] = 0

	//copying the real ip packets
	copy(frame[HeaderSize:], data)

	return frame
}
