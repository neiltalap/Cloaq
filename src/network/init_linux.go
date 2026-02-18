//go:build linux
// +build linux

package network

func InitTunnel() (Tunnel, error) {
	tun, err := InitLinuxTUN("cloaq0")
	if err != nil {
		return nil, err
	}

	if err := ConfigureLinuxTUN(tun.Name(), "10.0.0.2/24"); err != nil {
		return nil, err
	}

	return tun, nil
}
