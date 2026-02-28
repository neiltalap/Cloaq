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

package lintun

import (
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

func CreateTUN(name string) (*os.File, error) {

	fd, err := unix.Open("/dev/net/tun", unix.O_RDWR|unix.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}

	var req struct {
		Name  [unix.IFNAMSIZ]byte
		Flags uint16
	}
	copy(req.Name[:], name)

	req.Flags = unix.IFF_TUN | unix.IFF_NO_PI

	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		uintptr(unix.TUNSETIFF),
		uintptr(unsafe.Pointer(&req)),
	)

	if errno != 0 {
		err := unix.Close(fd)
		if err != nil {
			return nil, err
		}
		return nil, errno
	}

	// 4. Оборачиваем дескриптор в *os.File, чтобы Go мог с ним работать
	return os.NewFile(uintptr(fd), "/dev/net/tun"), nil
}
