//go:build linux

package netstat

func platformProvider() Provider { return NewLinuxProvider() }
