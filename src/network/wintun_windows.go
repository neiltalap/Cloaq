package network

import (
	"fmt"

	"golang.org/x/sys/windows"
)

var (
	wintunDLL = windows.NewLazyDLL("bin/amd64/wintun.dll")

	procWintunCreateAdapter = wintunDLL.NewProc("WintunCreateAdapter")
	procWintunDeleteAdapter = wintunDLL.NewProc("WintunDeleteAdapter")
)

func InitWinTun() error {
	if err := wintunDLL.Load(); err != nil {
		return fmt.Errorf("failed to load wintun.dll: %w", err)
	}

	fmt.Println("[WinTun] DLL loaded successfully")
	return nil
}
