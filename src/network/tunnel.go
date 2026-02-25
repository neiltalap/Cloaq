package network

import "os"

//It reads/writes raw IPv4/IPv6 packets
type Tunnel interface {
	Name() string
	Start() error
	Close() error

	Read(p []byte) (int, error)
	Write(p[]byte) (int, error)

	File() *os.File //Optional, :inux can provide an fd-backed os.File, wintun won't
}