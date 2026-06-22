//go:build windows

package kill

import (
	"fmt"

	"golang.org/x/sys/windows"
)

const processTerminate = 0x0001

func Process(pid uint32) Result {
	h, err := windows.OpenProcess(processTerminate, false, pid)
	if err != nil {
		return Result{PID: pid, OK: false, Reason: fmt.Sprintf("access denied (pid=%d): %v", pid, err)}
	}
	defer windows.CloseHandle(h)
	if err := windows.TerminateProcess(h, 1); err != nil {
		return Result{PID: pid, OK: false, Reason: fmt.Sprintf("terminate failed: %v", err)}
	}
	return Result{PID: pid, OK: true, Reason: "terminated"}
}
