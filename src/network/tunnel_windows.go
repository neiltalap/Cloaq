//go:build windows

package network

import (
	"fmt"
	"os"

	"cloaq/src/network/wintun"
)

type windowsTunnel struct {
	name    string
	adapter *wintun.Adapter
	session wintun.Session
	started bool
}

func (t *windowsTunnel) Name() string   { return t.name }
func (t *windowsTunnel) File() *os.File { return nil } //Since no fd on wintun

func (t *windowsTunnel) Start() error {
	if t.started {
		return nil
	}

	//4mb ring buffer suggested
	sess, err := t.adapter.StartSession(0x400000)
	if err != nil {
		return fmt.Errorf("Start session: %w", err)
	}

	t.session = sess
	t.started = true
	return nil
}

func (t *windowsTunnel) Close() error {
	if t.started {
		t.session.End()
		t.started = false
	}

	if t.adapter != nil {
		_ = t.adapter.Close()
		t.adapter = nil
	}

	return nil
}

func (t *windowsTunnel) Read(p []byte) (int, error) {
	if !t.started {
		return 0, fmt.Errorf("tunnel not started")
	}

	packet, err := t.session.ReceivePacket()
	if err != nil {
		return 0, err
	}

	n := copy(p, packet)
	t.session.ReleaseReceivePacket(packet)
	return n, nil
}

func (t *windowsTunnel) Write(p []byte) (int, error) {
	if !t.started {
		return 0, fmt.Errorf("tunnel not started")
	}

	buf, err := t.session.AllocateSendPacket(len(p))
	if err != nil {
		return 0, err
	}

	copy(buf, p)
	t.session.SendPacket(buf)
	return len(p), nil
}

// Init Tunnel creates/opens a WinTun adapter and returns a Tunnel
func InitTunnel() (Tunnel, error) {
	const name = "cloaq0"
	const tunnelType = "Cloaq"

	adapter, err := wintun.OpenAdapter(name)
	if err != nil {
		adapter, err = wintun.CreateAdapter(name, tunnelType, nil)
		if err != nil {
			return nil, fmt.Errorf("CreateAdapter: %w", err)
		}
	}

	return &windowsTunnel{
		name:    name,
		adapter: adapter,
	}, nil
}
