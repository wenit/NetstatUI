//go:build windows

package netstat

func platformProvider() Provider { return NewWindowsProvider() }
