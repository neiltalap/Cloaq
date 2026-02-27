package main

import (
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

func (t *Transport) SendTo(dst *net.UDPAddr, data []byte) error {
	_, err := t.conn.WriteToUDP(data, dst)
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
