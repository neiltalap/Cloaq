package network

import "os"

type Tunnel interface {
	Start() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	Name() string
	File() *os.File
}
