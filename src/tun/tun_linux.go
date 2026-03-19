//go:build linux

// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package tun

import (
	"cloaq/src/tun/lintun"
	"os"
	"os/exec"
)

type linuxDevice struct {
	name string
	f    *os.File
}

func (t *linuxDevice) Name() string                { return t.name }
func (t *linuxDevice) Start() error                { return nil }
func (t *linuxDevice) Close() error                { return t.f.Close() }
func (t *linuxDevice) Write(p []byte) (int, error) { return t.f.Write(p) }
func (t *linuxDevice) File() *os.File              { return t.f }
func (t *linuxDevice) Fd() int                     { return int(t.f.Fd()) }

// InitDevice creates a L3 TUN on Linux
func InitDevice() (Device, error) {
	name := "cloaq0"
	f, err := lintun.CreateTUN(name)
	if err != nil {
		return nil, err
	}

	err = exec.Command("ip", "link", "set", name, "up").Run()
	if err != nil {
		return nil, err
	}

	// added ipv6 address
	err = exec.Command("ip", "-6", "addr", "add", "fc00::1/64", "dev", name).Run()
	if err != nil {
		return nil, err
	}

	return &linuxDevice{name: name, f: f}, nil

}

// Read reads packets from the TUN device
func (d *linuxDevice) Read(buf []byte) (int, error) {
	if d.f == nil {
		return 0, os.ErrClosed
	}

	// Use os.File.Read to leverage Go's non-blocking I/O and netpoller
	n, err := d.f.Read(buf)
	if err != nil {
		return n, err
	}

	return n, nil
}
