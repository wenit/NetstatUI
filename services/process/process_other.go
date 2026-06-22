//go:build !windows

package process

func newWindowsProvider() provider { return nil }
