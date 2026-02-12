package main

import (
	"os"
)

const (
	// To get L3 interface
	InterfaceFlag_TUN = 0x0001
	// To get L2 interface
	InterfaceFlag_TAP = 0x0002
	// We don't want any header (4 bytes) added by the kernel
	// We want clean, raw Ethernet frames exactly as they
	// appear on the wire
	InterfaceFlag_NO_PI = 0x1000
)

type interfaceRequest struct {
	Name  [16]byte
	Flags uint16
}

func NewTAP(name string) (*os.File, error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)

	if err != nil {
		return nil, err
	}

	return file, nil
}
