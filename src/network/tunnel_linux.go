//go:build linux

package network

import "os"

type linuxTunnel struct {
	name string
	f    *os.File
}

func (t *linuxTunnel) Name() string                { return t.name }
func (t *linuxTunnel) Start() error                { return nil }
func (t *linuxTunnel) Close() error                { return t.f.Close() }
func (t *linuxTunnel) Read(p []byte) (int, error)  { return t.f.Read(p) }
func (t *linuxTunnel) Write(p []byte) (int, error) { return t.f.Write(p) }
func (t *linuxTunnel) File() *os.File              { return t.f }

// InitTunnel creates a L3 TUN on Linux
func InitTunnel() (Tunnel, error) {
	f, err := InitLinuxTUN("cloaq0") // your existing function
	if err != nil {
		return nil, err
	}
	return &linuxTunnel{name: "cloaq0", f: f}, nil
}
