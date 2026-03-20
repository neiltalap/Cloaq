//go:build linux

package tun

import (
	"cloaq/src/tun/lintun"
	"fmt"
	"log"
	"os"
	"runtime"

	"golang.org/x/sys/unix"
)

type LinuxDevice struct {
	f    *os.File
	name string
}

func (t *LinuxDevice) Name() string                { return t.name }
func (t *LinuxDevice) Start() error                { return nil }
func (t *LinuxDevice) Close() error                { return t.f.Close() }
func (t *LinuxDevice) Write(p []byte) (int, error) { return t.f.Write(p) }
func (t *LinuxDevice) File() *os.File              { return t.f }
func (t *LinuxDevice) Fd() int                     { return int(t.f.Fd()) }

func (t *LinuxDevice) Read(buf []byte) (int, error) {
	if t.f == nil {
		return 0, os.ErrClosed
	}
	n, err := unix.Read(int(t.f.Fd()), buf)
	if err != nil {
		if err == unix.EAGAIN || err == unix.EWOULDBLOCK {
			return 0, nil
		}
		return n, fmt.Errorf("direct read error: %w", err)
	}

	return n, nil
}

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
