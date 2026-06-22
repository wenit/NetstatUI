//go:build !windows

package system

func GetSystemLocale() string {
	return ""
}
