package tun

import "os"

// Device reads/writes raw IPv4/IPv6 packets on the TUN interface
type Device interface {
	Name() string
	Start() error
	Close() error

	Read(p []byte) (int, error)
	Write(p []byte) (int, error)

	File() *os.File
	Fd() int //Optional, linux can provide an fd-backed os.File, wintun won't
}
