//go:build linux

package network

import (
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

type LinuxTunnel struct {
	name string
	f    *os.File
}

func (t *LinuxTunnel) Name() string                { return t.name }
func (t *LinuxTunnel) Start() error                { return nil }
func (t *LinuxTunnel) Close() error                { return t.f.Close() }
func (t *LinuxTunnel) Read(p []byte) (int, error)  { return t.f.Read(p) }
func (t *LinuxTunnel) Write(p []byte) (int, error) { return t.f.Write(p) }

func NewTUN(name string) (*os.File, error) {

	fd, err := unix.Open("/dev/net/tun", unix.O_RDWR|unix.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}

	var req struct {
		Name  [16]byte
		Flags uint16
	}
	copy(req.Name[:], name)

	req.Flags = unix.IFF_TUN | unix.IFF_NO_PI

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(0x400454ca), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		err := unix.Close(fd)
		if err != nil {
			return nil, err
		}
		return nil, errno
	}

	return os.NewFile(uintptr(fd), "/dev/net/tun"), nil
}

func InitTunnel() (*LinuxTunnel, error) {
	f, err := NewTUN("cloaq0")
	if err != nil {
		return nil, err
	}
	return &LinuxTunnel{
		name: "cloaq0",
		f:    f,
	}, nil
}
