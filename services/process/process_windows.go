//go:build windows

package process

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	processQueryLimitedInfo = 0x1000
	processVmRead           = 0x0010
)

var (
	modKernel32               = windows.NewLazySystemDLL("kernel32.dll")
	procQueryFullProcessImage = modKernel32.NewProc("QueryFullProcessImageNameW")
)

type windowsProvider struct{}

func platformProvider() provider { return &windowsProvider{} }

// SnapshotAll enumerates every process via a single
// CreateToolhelp32Snapshot syscall (~1ms regardless of process count),
// returning PID/PPID/Name for each. Path is fetched lazily by QueryPath
// only for PIDs that actually appear in netstat output.
//
// This is dramatically faster than gopsutil's process.Processes(), which
// issues Name()+Exe()+Ppid() syscalls per process (~1300 syscalls for
// 400+ processes → 3-4s on a busy machine).
func (w *windowsProvider) SnapshotAll() (map[uint32]Info, error) {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("toolhelp snapshot: %w", err)
	}
	defer windows.CloseHandle(snap)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(snap, &entry); err != nil {
		return nil, fmt.Errorf("process32 first: %w", err)
	}
	out := make(map[uint32]Info, 256)
	for {
		pid := entry.ProcessID
		out[pid] = Info{
			PID:  pid,
			PPID: entry.ParentProcessID,
			Name: windows.UTF16ToString(entry.ExeFile[:]),
		}
		if err := windows.Process32Next(snap, &entry); err != nil {
			break
		}
	}
	return out, nil
}

func (w *windowsProvider) QueryPath(pid uint32) (string, error) {
	h, err := windows.OpenProcess(processQueryLimitedInfo|processVmRead, false, pid)
	if err != nil {
		h, err = windows.OpenProcess(processQueryLimitedInfo, false, pid)
		if err != nil {
			return "", err
		}
	}
	defer windows.CloseHandle(h)

	var buf [windows.MAX_PATH * 2]uint16
	size := uint32(len(buf))
	r1, _, e := procQueryFullProcessImage.Call(
		uintptr(h),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)
	if r1 == 0 {
		return "", fmt.Errorf("query image: %v", e)
	}
	return windows.UTF16ToString(buf[:size]), nil
}
