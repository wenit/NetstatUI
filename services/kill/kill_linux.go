//go:build linux

package kill

import (
	"fmt"

	gproc "github.com/shirou/gopsutil/v3/process"
)

// Process terminates the process with the given PID.
//
// On Linux, gopsutil's Kill() invokes syscall.Kill(pid, SIGKILL) — the same
// hard-kill semantics as the Windows TerminateProcess-based path and the
// macOS implementation. Use Kill() (not Terminate()) to ensure SIGKILL.
func Process(pid uint32) Result {
	p, err := gproc.NewProcess(int32(pid))
	if err != nil {
		return Result{PID: pid, OK: false, Reason: fmt.Sprintf("open process failed (pid=%d): %v", pid, err)}
	}
	if err := p.Kill(); err != nil {
		return Result{PID: pid, OK: false, Reason: fmt.Sprintf("terminate failed: %v", err)}
	}
	return Result{PID: pid, OK: true, Reason: "terminated"}
}
