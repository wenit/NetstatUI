//go:build !windows && !darwin

package netstat

// platformProvider returns nil on unsupported platforms (currently Linux).
// main.go fatals on these platforms before reaching here.
func platformProvider() Provider { return nil }
