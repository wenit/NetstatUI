//go:build !windows

package kill

func Process(pid uint32) Result {
	return Result{PID: pid, OK: false, Reason: "not supported on this platform"}
}
