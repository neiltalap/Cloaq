//go:build linux
// +build linux

package network

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

type LinuxTUN struct {
	file *os.File
	name string
}

func (t *LinuxTUN) Start() error { return nil }

func (t *LinuxTUN) Read(p []byte) (int, error)  { return t.file.Read(p) }
func (t *LinuxTUN) Write(p []byte) (int, error) { return t.file.Write(p) }
func (t *LinuxTUN) Close() error                { return t.file.Close() }
func (t *LinuxTUN) Name() string                { return t.name }

// Minimal ifreq for TUNSETIFF: IFNAMSIZ name + short flags + padding.
// This layout is the standard way to do it from Go without relying on unix.Ifreq internals.
type ifreqTun struct {
	Name  [unix.IFNAMSIZ]byte
	Flags uint16
	Pad   [40]byte // plenty of padding for the ioctl
}

func InitLinuxTUN(requestedName string) (*LinuxTUN, error) {
	f, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("open /dev/net/tun: %w", err)
	}

	var req ifreqTun
	if requestedName != "" {
		if len(requestedName) >= unix.IFNAMSIZ {
			_ = f.Close()
			return nil, fmt.Errorf("interface name too long: %q", requestedName)
		}
		copy(req.Name[:], requestedName)
	}

	req.Flags = uint16(unix.IFF_TUN | unix.IFF_NO_PI)

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		f.Fd(),
		uintptr(unix.TUNSETIFF),
		uintptr(unsafe.Pointer(&req)),
	)
	if errno != 0 {
		_ = f.Close()
		return nil, fmt.Errorf("ioctl(TUNSETIFF): %v", errno)
	}

	actualName := unix.ByteSliceToString(req.Name[:])

	return &LinuxTUN{
		file: f,
		name: actualName,
	}, nil
}
