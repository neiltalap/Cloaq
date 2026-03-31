//go:build windows

package wintun

import (
	"log"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

type loggerLevel int

const (
	logInfo loggerLevel = iota
	logWarn
	logErr
)

const AdapterNameMax = 128

type Adapter struct {
	handle uintptr
}

var (
	modwintun = windows.NewLazyDLL("wintun.dll")

	procWintunSetLogger               = modwintun.NewProc("WintunSetLogger")
	procWintunCreateAdapter           = modwintun.NewProc("WintunCreateAdapter")
	procWintunOpenAdapter             = modwintun.NewProc("WintunOpenAdapter")
	procWintunCloseAdapter            = modwintun.NewProc("WintunCloseAdapter")
	procWintunDeleteDriver            = modwintun.NewProc("WintunDeleteDriver")
	procWintunGetAdapterLUID          = modwintun.NewProc("WintunGetAdapterLUID")
	procWintunGetRunningDriverVersion = modwintun.NewProc("WintunGetRunningDriverVersion")
)

type TimestampedWriter interface {
	WriteWithTimestamp(p []byte, ts int64) (n int, err error)
}

func logMessage(level loggerLevel, timestamp uint64, msg *uint16) int {
	if tw, ok := log.Default().Writer().(TimestampedWriter); ok {
		tw.WriteWithTimestamp([]byte(log.Default().Prefix()+windows.UTF16PtrToString(msg)), (int64(timestamp)-116444736000000000)*100)
	} else {
		log.Println(windows.UTF16PtrToString(msg))
	}
	return 0
}

func setupLogger() {
	var callback uintptr
	if runtime.GOARCH == "386" {
		callback = windows.NewCallback(func(level loggerLevel, timestampLow, timestampHigh uint32, msg *uint16) int {
			return logMessage(level, uint64(timestampHigh)<<32|uint64(timestampLow), msg)
		})
	} else if runtime.GOARCH == "arm" {
		callback = windows.NewCallback(func(level loggerLevel, _, timestampLow, timestampHigh uint32, msg *uint16) int {
			return logMessage(level, uint64(timestampHigh)<<32|uint64(timestampLow), msg)
		})
	} else {
		callback = windows.NewCallback(logMessage) // amd64/arm64
	}

	_, _, _ = procWintunSetLogger.Call(callback)
}

func ensureLoaded() error {
	if err := modwintun.Load(); err != nil {
		return err
	}
	setupLogger()
	return nil
}

func closeAdapter(w *Adapter) {
	_, _, _ = procWintunCloseAdapter.Call(w.handle)
}

func CreateAdapter(name string, tunnelType string, requestedGUID *windows.GUID) (wintun *Adapter, err error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}

	name16, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	tunnelType16, err := windows.UTF16PtrFromString(tunnelType)
	if err != nil {
		return nil, err
	}

	r0, _, e1 := procWintunCreateAdapter.Call(
		uintptr(unsafe.Pointer(name16)),
		uintptr(unsafe.Pointer(tunnelType16)),
		uintptr(unsafe.Pointer(requestedGUID)),
	)
	if r0 == 0 {
		return nil, e1
	}

	wintun = &Adapter{handle: r0}
	runtime.SetFinalizer(wintun, closeAdapter)
	return wintun, nil
}

func OpenAdapter(name string) (wintun *Adapter, err error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}

	name16, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}

	r0, _, e1 := procWintunOpenAdapter.Call(uintptr(unsafe.Pointer(name16)))
	if r0 == 0 {
		return nil, e1
	}

	wintun = &Adapter{handle: r0}
	runtime.SetFinalizer(wintun, closeAdapter)
	return wintun, nil
}

func (wintun *Adapter) Close() error {
	if wintun == nil || wintun.handle == 0 {
		return nil
	}
	_ = ensureLoaded()

	runtime.SetFinalizer(wintun, nil)
	r1, _, e1 := procWintunCloseAdapter.Call(wintun.handle)
	if r1 == 0 {
		return e1
	}
	wintun.handle = 0
	return nil
}

func Uninstall() error {
	if err := ensureLoaded(); err != nil {
		return err
	}
	r1, _, e1 := procWintunDeleteDriver.Call()
	if r1 == 0 {
		return e1
	}
	return nil
}

func RunningVersion() (uint32, error) {
	if err := ensureLoaded(); err != nil {
		return 0, err
	}
	r0, _, e1 := procWintunGetRunningDriverVersion.Call()
	v := uint32(r0)
	if v == 0 {
		return 0, e1
	}
	return v, nil
}

func (wintun *Adapter) LUID() (luid uint64) {
	if wintun == nil || wintun.handle == 0 {
		return 0
	}
	_ = ensureLoaded()

	_, _, _ = procWintunGetAdapterLUID.Call(
		wintun.handle,
		uintptr(unsafe.Pointer(&luid)),
	)
	return luid
}
