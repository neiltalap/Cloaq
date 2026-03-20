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
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/sys/unix"
)

type LinuxDevice struct {
	name string
	f    *os.File
}

func (t *LinuxDevice) Name() string                { return t.name }
func (t *LinuxDevice) Close() error                { return t.f.Close() }
func (t *LinuxDevice) Write(p []byte) (int, error) { return t.f.Write(p) }
func (t *LinuxDevice) File() *os.File              { return t.f }
func (t *LinuxDevice) Fd() int                     { return int(t.f.Fd()) }

func InitDevice(name string) (*LinuxDevice, error) {
	log.Printf("Initializing TUN device: %s", name)

	dev, err := lintun.CreateTUN(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create TUN interface: %w", err)
	}

	var f *os.File
	switch v := any(dev).(type) {
	case *os.File:
		f = v
	case interface{ File() *os.File }:
		f = v.File()
	default:
		return nil, fmt.Errorf("unsupported TUN device type: %T", dev)
	}

	rawConn, err := f.SyscallConn()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw connection: %w", err)
	}

	err = rawConn.Control(func(fd uintptr) {
		if runtime.GOOS == "linux" {
			errSet := unix.SetNonblock(int(fd), false)
			if errSet != nil {
				log.Printf("warning: failed to set blocking mode on fd %d: %v", fd, errSet)
			} else {
				log.Printf("successfully set blocking mode on fd %d", fd)
			}
		}
	})

	if err != nil {
		return nil, fmt.Errorf("error during raw control: %w", err)
	}

	return &LinuxDevice{
		f:    f,
		name: name,
	}, nil
}
func (t *LinuxDevice) Read(p []byte) (n int, err error) {
	fd := int(t.f.Fd())
	var readErr error
	n, err = unix.Read(fd, p)
	if err != nil {
		if err == unix.EAGAIN || err == unix.EWOULDBLOCK {
			return 0, nil
		}
		return 0, fmt.Errorf("unix read error: %v", err)
	}

	return n, readErr
}

func (t *LinuxDevice) Start() error {
	log.Printf("automatic network setup for %s...", t.name)

	cmdUp := exec.Command("ip", "link", "set", "dev", t.name, "up")
	if err := cmdUp.Run(); err != nil {
		return fmt.Errorf("failed to bring up interface: %w", err)
	}

	cmdAddr := exec.Command("ip", "-6", "addr", "add", "fc00::1/64", "dev", t.name)
	if err := cmdAddr.Run(); err != nil {

		log.Printf("note: IPv6 address might already be set: %v", err)
	}

	log.Printf("interface %s is READY: UP and IPv6 assigned (fc00::1)", t.name)
	return nil
}
