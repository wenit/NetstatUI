//go:build darwin

package netstat

func platformProvider() Provider { return NewDarwinProvider() }
