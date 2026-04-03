//go:build linux

package lintun

import (
	"log"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	InterfaceFlag_TUN   = 0x0001
	InterfaceFlag_TAP   = 0x0002
	InterfaceFlag_NO_PI = 0x1000
)

type interfaceRequest struct {
	Name  [unix.IFNAMSIZ]byte
	Flags uint16
}

// linus standart ioctl
const TUNSETIFF = 0x400454ca
const IFF_TUN = 0x0001
const IFF_NO_PI = 0x1000

type ifreq struct {
	name  [16]byte
	flags uint16
	_     [22]byte
}

func CreateTUN(name string) (*os.File, error) {
	fileDescriptor, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	var req interfaceRequest
	copy(req.Name[:], name)
	req.Flags = InterfaceFlag_TUN | InterfaceFlag_NO_PI

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		fileDescriptor.Fd(),
		uintptr(unix.TUNSETIFF),
		uintptr(unsafe.Pointer(&req)),
	)
	if errno != 0 {
		err := fileDescriptor.Close()
		if err != nil {
			return nil, err
		}
		return nil, errno
	}

	log.Println("tun interface created: ", fileDescriptor)
	return fileDescriptor, nil
}
