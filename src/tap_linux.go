// NOTICE

// Project Name: Cloaq
// Copyright © 2026 Neil Talap and/or its designated Affiliates.

// This software is licensed under the Dragonfly Public License (DPL) 1.0.

// All rights reserved. The names "Neil Talap" and any associated logos or branding
// are trademarks of the Licensor and may not be used without express written permission,
// except as provided in Section 7 of the License.

// For commercial licensing inquiries or permissions beyond the scope of this
// license, please create an issue in github.

package main

import (
	"log"
	"net"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

/*
	we're gonna have to disable
*/

const (
	// To get L3 interface
	InterfaceFlag_TUN = 0x0001
	// To get L2 interface
	InterfaceFlag_TAP = 0x0002
	// We don't want any header (4 bytes) added by the kernel
	// We want clean, raw Ethernet frames exactly as they
	// appear on the wire
	InterfaceFlag_NO_PI = 0x1000
)

type interfaceRequest struct {
	Name  [unix.IFNAMSIZ]byte //const unix.InterfaceNameSize untyped int = 0x10 => 16
	Flags uint16
}

type Route struct {
	Prefix *net.IPNet
	OutIf  string
}

var routes []Route

func CreateRouter(tunFileDescriptor *os.File) {
	// while(true) to start constantly reading incoming traffic stored in /dev/tun
	buf := make([]byte, 65535)
	for {
		n, err := tunFileDescriptor.Read(buf)
		if err != nil {
			continue
		}

		packet := buf[:n]

		if len(packet) < 40 {
			continue
		}

		// Decrement Hop Limit
		if packet[7] <= 1 {
			continue
		}
		packet[7]--

		dst := net.IP(packet[24:40])

		outIf := LookupRoute(dst)
		if outIf == "" {
			continue
		}

		/*
			FORWARD TRAFFIC TO ANOTHER NODE
		*/

		SendPacket(outIf, packet)
	}
}

func CreateIPv6PacketListener(tunFileDescriptor *os.File) {
	buf := make([]byte, 65535)
	for {
		n, err := tunFileDescriptor.Read(buf)
		if err != nil {
			continue
		}

		packet := buf[:n]

		payload := packet[40:]
		log.Printf("Payload (%d bytes): % x\n", len(payload), payload)
	}
}

func NewTUN(name string) *os.File {
	fileDescriptor, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)

	if err != nil {
		log.Fatalf("open /dev/net/tun failed: %v", err)
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
		fileDescriptor.Close()
		log.Fatalf("ioctl TUNSETIFF failed: %v", errno)
	}

	log.Println("TUN interface created: ", fileDescriptor)
	return fileDescriptor
}

func AddRoute(cidr, outIf string) {
	_, netw, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Fatal(err)
	}

	routes = append(routes, Route{
		Prefix: netw,
		OutIf:  outIf,
	})
}

func LookupRoute(dst net.IP) string {
	for _, r := range routes {
		if r.Prefix.Contains(dst) {
			return r.OutIf
		}
	}
	return ""
}

// send packet through a raw socket
func SendPacket(ifName string, packet []byte) {
	iface, err := net.InterfaceByName(ifName)
	if err != nil {
		return
	}

	fd, err := unix.Socket(
		unix.AF_PACKET,
		unix.SOCK_RAW,
		int(htons(0x86DD)),
	)
	if err != nil {
		return
	}
	defer unix.Close(fd)

	sll := &unix.SockaddrLinklayer{
		Ifindex:  iface.Index,
		Protocol: htons(0x86DD),
	}

	// Kernel will add L2 header automatically
	unix.Sendto(fd, packet, 0, sll)
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}
