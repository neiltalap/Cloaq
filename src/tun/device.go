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

import "os"

// Device reads/writes raw IPv4/IPv6 packets on the TUN interface
type Device interface {
	Name() string
	Start() error
	Close() error

	Read(p []byte) (int, error)
	Write(p []byte) (int, error)

	File() *os.File
	Fd() int //Optional, linux can provide an fd-backed os.File, wintun won't
}
