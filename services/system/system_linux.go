//go:build linux

package system

import "os"

// GetSystemLocale returns the Linux user locale by reading $LANG.
//
// On most Linux distros LANG is set to something like `en_US.UTF-8` or
// `zh_CN.UTF-8`. Returns empty string when unset (caller falls back to
// navigator / en-US).
func GetSystemLocale() string {
	return os.Getenv("LANG")
}
