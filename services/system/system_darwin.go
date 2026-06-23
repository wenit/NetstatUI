//go:build darwin

package system

import "os/exec"

// GetSystemLocale returns the macOS user locale via osascript.
//
// Output is typically `en_US@calendar=gregorian` or `zh_CN@calendar=gregorian`.
// Empty string if the lookup fails (caller falls back to navigator / en-US).
func GetSystemLocale() string {
	out, err := exec.Command("osascript", "-e", "user locale of (system info)").Output()
	if err != nil {
		return ""
	}
	return trim(string(out))
}

func trim(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r' || s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
