//go:build !windows && !darwin

package system

func GetSystemLocale() string {
	return ""
}
