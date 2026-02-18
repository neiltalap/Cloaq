package network

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

var (
	wintunDLL               *windows.LazyDLL
	procWintunCreateAdapter *windows.LazyProc
	procWintunDeleteAdapter *windows.LazyProc
)

func InitWinTun() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("os.Executable: %w", err)
	}
	exeDir := filepath.Dir(exe)

	dllPath := filepath.Join(exeDir, "bin", "amd64", "wintun.dll")

	cwd, _ := os.Getwd()
	fmt.Println("[WinTun] cwd:", cwd)
	fmt.Println("[WinTun] exe:", exe)
	fmt.Println("[WinTun] trying dll:", dllPath)

	wintunDLL = windows.NewLazyDLL(dllPath)
	procWintunCreateAdapter = wintunDLL.NewProc("WintunCreateAdapter")
	procWintunDeleteAdapter = wintunDLL.NewProc("WintunDeleteAdapter")

	if err := wintunDLL.Load(); err != nil {
		return fmt.Errorf("failed to load wintun.dll from %q: %w", dllPath, err)
	}

	fmt.Println("[WinTun] DLL loaded successfully")
	return nil
}
