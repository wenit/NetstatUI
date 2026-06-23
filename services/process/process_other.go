//go:build !windows && !darwin

package process

func platformProvider() provider { return nil }
