//go:build windows

package kill

import (
	"fmt"

	gproc "github.com/shirou/gopsutil/v3/process"
)

// Process terminates the process with the given PID.
//
// On Windows, gopsutil's Kill() invokes TerminateProcess, which matches the
// hard-kill semantics of the previous direct TerminateProcess implementation.
func Process(pid uint32) Result {
	p, err := gproc.NewProcess(int32(pid))
	if err != nil {
		return Result{PID: pid, OK: false, Reason: fmt.Sprintf("access denied (pid=%d): %v", pid, err)}
	}
	if err := p.Kill(); err != nil {
		return Result{PID: pid, OK: false, Reason: fmt.Sprintf("terminate failed: %v", err)}
	}
	return Result{PID: pid, OK: true, Reason: "terminated"}
}
