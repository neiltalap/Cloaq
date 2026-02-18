package network

import (
	"encoding/binary"
	"fmt"
	"net"
)

func DescribeIPPacket(b []byte) string {
	if len(b) < 1 {
		return "empty"
	}
	v := b[0] >> 4
	switch v {
	case 4:
		return describeIPv4(b)
	case 6:
		return describeIPv6(b)
	default:
		return fmt.Sprintf("unknown ip version=%d len=%d", v, len(b))
	}
}

func describeIPv4(b []byte) string {
	if len(b) < 20 {
		return fmt.Sprintf("ipv4 short len=%d", len(b))
	}
	ihl := int(b[0]&0x0F) * 4
	if ihl < 20 || len(b) < ihl {
		return fmt.Sprintf("ipv4 bad ihl=%d len=%d", ihl, len(b))
	}
	totalLen := int(binary.BigEndian.Uint16(b[2:4]))
	proto := b[9]
	src := net.IP(b[12:16])
	dst := net.IP(b[16:20])
	return fmt.Sprintf("ipv4 %s -> %s proto=%d iplen=%d caplen=%d", src, dst, proto, totalLen, len(b))
}

func describeIPv6(b []byte) string {
	if len(b) < 40 {
		return fmt.Sprintf("ipv6 short len=%d", len(b))
	}
	nextHdr := b[6]
	payloadLen := int(binary.BigEndian.Uint16(b[4:6]))
	src := net.IP(b[8:24])
	dst := net.IP(b[24:40])
	return fmt.Sprintf("ipv6 %s -> %s next=%d paylen=%d caplen=%d", src, dst, nextHdr, payloadLen, len(b))
}
