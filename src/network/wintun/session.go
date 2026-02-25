//go:build windows

package wintun

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type Session struct {
	handle uintptr
}

const (
	PacketSizeMax   = 0xffff
	RingCapacityMin = 0x20000
	RingCapacityMax = 0x4000000
)

type Packet struct {
	Next *Packet
	Size uint32
	Data *[PacketSizeMax]byte
}

var (
	procWintunAllocateSendPacket   = modwintun.NewProc("WintunAllocateSendPacket")
	procWintunEndSession           = modwintun.NewProc("WintunEndSession")
	procWintunGetReadWaitEvent     = modwintun.NewProc("WintunGetReadWaitEvent")
	procWintunReceivePacket        = modwintun.NewProc("WintunReceivePacket")
	procWintunReleaseReceivePacket = modwintun.NewProc("WintunReleaseReceivePacket")
	procWintunSendPacket           = modwintun.NewProc("WintunSendPacket")
	procWintunStartSession         = modwintun.NewProc("WintunStartSession")
)

func (wintun *Adapter) StartSession(capacity uint32) (Session, error) {
	if err := ensureLoaded(); err != nil {
		return Session{}, err
	}
	r0, _, e1 := procWintunStartSession.Call(wintun.handle, uintptr(capacity))
	if r0 == 0 {
		return Session{}, e1
	}
	return Session{handle: r0}, nil
}

func (session Session) End() {
	if session.handle == 0 {
		return
	}
	_ = ensureLoaded()
	_, _, _ = procWintunEndSession.Call(session.handle)
}

func (session Session) ReadWaitEvent() windows.Handle {
	if session.handle == 0 {
		return 0
	}
	_ = ensureLoaded()
	r0, _, _ := procWintunGetReadWaitEvent.Call(session.handle)
	return windows.Handle(r0)
}

func (session Session) ReceivePacket() ([]byte, error) {
	if session.handle == 0 {
		return nil, windows.ERROR_INVALID_HANDLE
	}
	_ = ensureLoaded()

	var packetSize uint32
	r0, _, e1 := procWintunReceivePacket.Call(session.handle, uintptr(unsafe.Pointer(&packetSize)))
	if r0 == 0 {
		return nil, e1
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(r0)), packetSize), nil
}

func (session Session) ReleaseReceivePacket(packet []byte) {
	if session.handle == 0 || len(packet) == 0 {
		return
	}
	_ = ensureLoaded()
	_, _, _ = procWintunReleaseReceivePacket.Call(session.handle, uintptr(unsafe.Pointer(&packet[0])))
}

func (session Session) AllocateSendPacket(packetSize int) ([]byte, error) {
	if session.handle == 0 {
		return nil, windows.ERROR_INVALID_HANDLE
	}
	_ = ensureLoaded()

	r0, _, e1 := procWintunAllocateSendPacket.Call(session.handle, uintptr(packetSize))
	if r0 == 0 {
		return nil, e1
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(r0)), packetSize), nil
}

func (session Session) SendPacket(packet []byte) {
	if session.handle == 0 || len(packet) == 0 {
		return
	}
	_ = ensureLoaded()
	_, _, _ = procWintunSendPacket.Call(session.handle, uintptr(unsafe.Pointer(&packet[0])))
}
