//go:build windows
// +build windows

package network

func InitTunnel() (Tunnel, error) {
	if err := InitWinTun(); err != nil {
		return nil, err
	}

	// WIll need to rewrite a real wintun tunnel
	return nil, nil
}
