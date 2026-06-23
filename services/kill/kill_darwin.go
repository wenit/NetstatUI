//go:build darwin

package kill

import (
	"fmt"

	gproc "github.com/shirou/gopsutil/v3/process"
)

// Process terminates the process with the given PID.
//
// On macOS, gopsutil's Kill() invokes syscall.Kill(pid, SIGKILL), which is
// the same hard-kill behavior as the original TerminateProcess-based path on
// Windows. Use Kill() (not Terminate()) to ensure SIGKILL semantics.
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
