package network

import (
	"fmt"
	"os/exec"
)

func ConfigureLinuxTUN(ifName, cidr string) error {
	if out, err := exec.Command("ip", "link", "set", "dev", ifName, "up").CombinedOutput(); err != nil {
		return fmt.Errorf("ip link up failed: %w: %s", err, string(out))
	}

	if out, err := exec.Command("ip", "addr", "add", cidr, "dev", ifName).CombinedOutput(); err != nil {
		return fmt.Errorf("ip addr add failed: %w: %s", err, string(out))
	}

	return nil
}
